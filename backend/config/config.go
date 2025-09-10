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

// LoadConfig loads application configuration from a file and returns a populated Config.
// 
// LoadConfig looks for a YAML file named `config.yaml` inside a `config` subdirectory of
// the current working directory (i.e. `<workDir>/config/config.yaml`). It uses Viper to
// read and unmarshal the file into a Config value and returns a pointer to that value.
// If reading or unmarshalling the configuration fails, the function logs a fatal error and
// terminates the program.
func LoadConfig() *Config {
	workDir, _ := os.Getwd()
	configFile := filepath.Join(workDir, "config", "config.yaml")
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
