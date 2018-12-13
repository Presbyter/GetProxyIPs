package repository

import (
	"errors"
	"get_proxy_ips/common"
	"github.com/globalsign/mgo/bson"
	"log"
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
	// 验证是否已经存在
	s := common.GetAppBus().MongoDbSession.Copy()
	defer s.Close()

	c := s.DB("ip_pool").C("pool")
	count, err := c.Find(bson.M{"ip": entity.Ip, "port": entity.Port}).Count()
	if err != nil {
		log.Printf("mgo find pool fail. error: %v", err.Error())
		return
	}

	if count > 0 {
		// 已经存在
		return errors.New("the ip is exist.")
	}

	err = c.Insert(&entity)
	if err != nil {
		log.Printf("mgo insert pool fail. error: %v", err.Error())
		return
	}

	log.Println("mgo insert pool success.")
	return
}
