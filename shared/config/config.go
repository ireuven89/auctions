package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Sql    DBConfig     `mapstructure:"database"`
	Redis  DBConfig     `mapstructure:"redis"`
	Server ServerConfig `mapstructure:"server"`
	AWS    AWSConfig    `mapstructure:"aws"`
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

type AWSConfig struct {
	S3Buckets struct {
		Primary string `mapstructure:"primary"`
	} `mapstructure:"s3_buckets"`
	S3Region string
}

const defaultConfigDir = "/config"
const defaultPublicKeyPath = "/config/public.key"

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

func LoadRSAPublicKeyFromEnv() (*rsa.PublicKey, error) {
	publicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = defaultPublicKeyPath
	}
	publicKeyPemFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyPemFile)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %v", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}
	return rsaPub, nil
}
