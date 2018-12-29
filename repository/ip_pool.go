package repository

import (
	"errors"
	"get_proxy_ips/common"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

type Pool struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Ip          string        `bson:"ip" json:"ip"`
	Port        int           `bson:"port" json:"port"`
	Type        string        `bson:"type" json:"type"`
	Location    string        `bson:"location" json:"location"`
	LastTryTime int64         `bson:"last_try_time" json:"last_try_time"`
	Status      int           `bson:"status" json:"status"` // 状态; 0:不可用; 1:可用;
}

func (p *Pool) Create(entity Pool) (err error) {
	s := common.GetAppBus().MongoDbSession.Copy()
	defer s.Close()

	c := s.DB("ip_pool").C("pool")
	// 验证是否已经存在
	count, err := c.Find(bson.M{"ip": entity.Ip, "port": entity.Port}).Count()
	if err != nil {
		log.Printf("error | mgo find pool fail. error: %v", err.Error())
		return
	}

	if count > 0 {
		// 已经存在
		return errors.New("the ip is exist.")
	}

	err = c.Insert(&entity)
	if err != nil {
		log.Printf("error | mgo insert pool fail. error: %v", err.Error())
		return
	}

	log.Println("mgo insert pool success.")
	return
}

func (p *Pool) GetByPage(pageIndex, pageSize int) (list []Pool, err error) {
	s := common.GetAppBus().MongoDbSession.Copy()
	defer s.Close()

	c := s.DB("ip_pool").C("pool")
	err = c.Find(nil).Sort("_id").Skip(pageSize * (pageIndex - 1)).Limit(pageSize).All(&list)
	if err != nil {
		log.Printf("error | mgo find pool fail. error: %v", err.Error())
		return
	}
	return
}

func (p *Pool) GetTotalCount() (count int, err error) {
	s := common.GetAppBus().MongoDbSession.Copy()
	defer s.Close()

	c := s.DB("ip_pool").C("pool")
	count, err = c.Count()
	if err != nil {
		log.Printf("error | get pool count fail. error: %v", err.Error())
	}
	return
}

func (p *Pool) Modify() (err error) {
	s := common.GetAppBus().MongoDbSession.Copy()
	defer s.Close()

	c := s.DB("ip_pool").C("pool")
	err = c.UpdateId(p.Id, bson.M{"$set": bson.M{"last_try_time": time.Now().Unix(), "status": p.Status}})
	if err != nil {
		log.Printf("error | mgo modify pool fail. error: %v", err.Error())
	}
	return
}
