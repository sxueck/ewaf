package config

import (
	"github.com/spf13/viper"
	"log"
)

var Cfg = &cfg{}

func init() {
	v := viper.New()
	v.SetConfigName("sites-enabled")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("/configs") // docker mode
	v.AddConfigPath("/etc/ewaf")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	if err := v.Unmarshal(Cfg); err != nil {
		log.Fatalf("unmarshal configuration to model error : %s", err)
	}
}
