// Package config ...
package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Config is a config
type Config struct {
	BotToken          string  `mapstructure:"BOT_TOKEN"`
	PathToTargetImage string  `mapstructure:"PATH_TO_TARGET_IMAGE"`
	TargetChannels    []int64 `mapstructure:"TARGET_CHANNELS"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get(pathToConfig string) *Config {
	once.Do(func() {
		viper.SetConfigFile(pathToConfig)
		viper.SetConfigType("env")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("error reading env file", err)
		}

		// Viper unmarshals the loaded env varialbes into the struct
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal(err)
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Configuration:", string(configBytes))
	})

	return &config
}
