package common

import (
	"github.com/globalsign/mgo"
	"log"
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
		session, err := mgo.Dial(GetConfig().MongoConnStr)
		if err != nil {
			log.Fatalf("open mongo db fail. error: %v", err.Error())
		}
		return &AppBus{
			MongoDbSession: session,
		}
	}
}
