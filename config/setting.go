package config

import (
	"bufio"
	"os"

	"github.com/spf13/viper"
)

// ENV is app env
var APP_ENV string

const (
	// DEVELOPMENT env
	DEVELOPMENT = "development"
	// TEST env
	TEST = "test"
	// PRODUCTION env
	PRODUCTION = "production"
)

// Setup ...
func Setup(file string) error {
	APP_ENV = os.Getenv("APP_ENV")
	if APP_ENV == "" {
		APP_ENV = DEVELOPMENT
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(bufio.NewReader(f))

	return nil
}
