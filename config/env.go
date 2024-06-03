package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost 		string `mapstructure:"POSTGRES_HOST"`
	DBUser 		string `mapstructure:"POSTGRES_USER"`
	DBPassword 	string `mapstructure:"POSTGRES_PASSWORD"`
	DBName 		string `mapstructure:"POSTGRES_DB"`
	DBPort 		string `mapstructure:"POSTGRES_PORT"`
	ServerPort 	string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (config Config, err error) {
	appEnv := os.Args[1]
    if appEnv == "" {
        appEnv = ""
    }

	viper.SetConfigName(appEnv)
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	return
}