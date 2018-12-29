package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"get_proxy_ips/repository"
	"github.com/globalsign/mgo/bson"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func GetIpFromSource() {
	c := colly.NewCollector(
		colly.AllowedDomains("kuaidaili.com", "www.kuaidaili.com"),
	)

	c.OnHTML("div#list table tbody tr", func(e *colly.HTMLElement) {
		m := make(map[string]string)
		e.ForEach("td", func(i int, el *colly.HTMLElement) {
			switch i {
			case 0:
				m["ip"] = el.Text
			case 1:
				m["port"] = el.Text
			case 3:
				m["type"] = el.Text
			case 4:
				m["location"] = el.Text
			}
			//log.Println(i, el.Text)
		})
		log.Println(m)

		port, _ := strconv.Atoi(m["port"])

		// todo save to db
		pool := repository.Pool{}
		err := pool.Create(repository.Pool{
			Id:       bson.NewObjectId(),
			Ip:       m["ip"],
			Port:     port,
			Type:     m["type"],
			Location: m["location"],
		})
		if err != nil {
			log.Printf("error | save to db fail. error: %v", err.Error())
		}
	})

	c.OnHTML("div#listnav ul li a[href]", func(e *colly.HTMLElement) {
		//log.Printf("the url: %v", e.Attr("href"))
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		time.Sleep(1500 * time.Millisecond)
		log.Printf("visiting: %v", r.URL)
	})

	c.Visit("https://www.kuaidaili.com/free/inha/1/")
	return
}

func CleanIps(ctx context.Context) {
	pool := repository.Pool{}

	ticker := time.NewTicker(10 * time.Second)

	ch := make(chan repository.Pool)
	limitCh := make(chan bool, 10)

	go func() {
		index := 0
		for item := range ch {
			limitCh <- true
			index++
			log.Println(index)
			go TestIpByBaidu(limitCh, item)
		}
	}()

	for {
		pageIndex, pageSize := 1, 500
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go func() {
				count, err := pool.GetTotalCount()
				if err != nil {
					return
				}

				for i := 0; i < int(math.Ceil(float64(count)/float64(pageSize))); i++ {
					list, err := pool.GetByPage(pageIndex+i, pageSize)
					if err != nil {
						continue
					}

					for _, v := range list {
						ch <- v
					}
				}
			}()
		}
		ticker = time.NewTicker(4 * time.Hour)
	}
}

func TestIpByBaidu(limit chan bool, p repository.Pool) {
	defer func() { <-limit }()
	var uri string
	switch strings.ToLower(p.Type) {
	case "http":
		uri = fmt.Sprintf("http://%v:%v", strings.Trim(p.Ip, " "), p.Port)
	default:
		uri = fmt.Sprintf("http://%v:%v", strings.Trim(p.Ip, " "), p.Port)
	}
	log.Println(uri)
	proxyUrl, err := url.Parse(uri)
	if err != nil {
		log.Printf("error | the uri is bad. error: %v", err.Error())
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Get("https://www.baidu.com")
	if err != nil {
		log.Printf("error | test by baidu fail. error: %v", err.Error())
		p.Status = 0
		p.Modify()
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		p.Status = 0
		p.Modify()
	} else {
		p.Status = 1
		p.Modify()
	}
}
