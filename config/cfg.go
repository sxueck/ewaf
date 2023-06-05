package config

import (
	"github.com/spf13/viper"
	"log"
)

func InitParse(inf any) any {
	v := viper.New()
	v.SetConfigName("sites-enabled")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("/configs") // docker mode
	v.AddConfigPath("/etc/ewaf")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}

	if err := v.Unmarshal(inf); err != nil {
		log.Fatalf("unmarshal configuration to model error : %s", err)
	}

	return inf
}
