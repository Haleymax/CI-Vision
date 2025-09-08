package config

import (
	"civ/config/autoload"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	appConfig *Config
	once      sync.Once
)

type Config struct {
	MySQL  autoload.MySQLConfig  `mapstructure:"mysql"`
	System autoload.SystemConfig `mapstructure:"system"`
}

func LoadConfig() *Config {
	workDir, _ := os.Getwd()
	configFile := filepath.Join(workDir, "config.yaml")
	log.Println("Loading config from ", configFile)

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	log.Println("Using config file:", viper.AllSettings())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}
	return &config
}

func GetConfig() *Config {
	once.Do(func() {
		appConfig = LoadConfig()
	})
	return appConfig
}
