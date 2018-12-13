package common

import (
	"github.com/globalsign/mgo"
	"log"
	"time"
)

type AppBus struct {
	MongoDbSession *mgo.Session
}

var (
	appBus *AppBus
)

func GetAppBus() *AppBus {
	if appBus != nil {
		return appBus
	} else {
		dialInfo := mgo.DialInfo{
			Addrs:     []string{"127.0.0.1"},
			Direct:    false,
			Timeout:   time.Second * 1,
			Database:  "ip_pool",
			Source:    "admin",
			Username:  "root",
			Password:  "315215241",
			PoolLimit: 4096, // Session.SetPoolLimit
		}
		session, err := mgo.DialWithInfo(&dialInfo)
		if err != nil {
			log.Fatalf("open mongo db fail. error: %v", err.Error())
		}
		return &AppBus{
			MongoDbSession: session,
		}
	}
}
