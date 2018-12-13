package handler

import (
	"get_proxy_ips/repository"
	"github.com/globalsign/mgo/bson"
	"log"
	"strconv"
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
