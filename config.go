package main

import (
	"github.com/spf13/viper"
)

type Mattermost struct {
	Token string
}

type Config struct {
	Mattermost Mattermost
}

// loadConfigFile receive a path and return config object
func loadConfigFile() (*Config, error) {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("COWSAY")

	_ = viper.BindEnv("TOKEN")
	token := viper.GetString("TOKEN")

	if token == "" {
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		token = viper.GetStringMap("mattermost")["token"].(string)
	}

	return &Config{
		Mattermost: Mattermost{
			Token: token,
		},
	}, nil
}
