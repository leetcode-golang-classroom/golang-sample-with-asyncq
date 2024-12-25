package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RedisURL string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}

var AppConfig *Config

func init() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	FailOnError(v.BindEnv("REDIS_URL"), "failed to bind REDIS_URL")
	FailOnError(v.BindEnv("PORT"), "failed to bind PORT")
	err := v.ReadInConfig()
	if err != nil {
		log.Println("Load from environment variable")
	}
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		FailOnError(err, "Failed to read enivronment")
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
