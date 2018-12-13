package common

import (
	"github.com/spf13/viper"
	"log"
)

type MongoCfg struct {
	Addrs     []string
	Username  string
	Password  string
	Database  string
	PoolLimit int
}

type config struct {
	Mongo MongoCfg
}

var (
	cfg *config
)

func init() {
	GetConfig()
}

func GetConfig() *config {
	if cfg != nil {
		return cfg
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		cfg = &config{}
		err = viper.Unmarshal(cfg)
		if err != nil {
			log.Panicf("viper unmarshal fail. %v", err.Error())
		}
		return cfg
	}
}
