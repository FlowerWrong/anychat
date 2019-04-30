package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ENV string

const DEVELOPMENT = "development"
const TEST = "test"
const PRODUCTION = "production"

func Setup() error {
	ENV = os.Getenv("APP_ENV")
	if ENV == "" {
		ENV = DEVELOPMENT
	}
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetConfigName("settings" + "." + ENV)
	err = viper.MergeInConfig()
	if err != nil {
		return err
	}
	return nil
}
