package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Sql    DBConfig     `mapstructure:"database"`
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

const defaultConfigDir = "/config"

func LoadConfig() (*Config, error) {
	var config *Config
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = defaultConfigDir
	}

	fmt.Printf("config path is %s\n", configDir)

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
