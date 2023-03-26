package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

func GetValueFromEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

type Config struct {
	Port     int    `json:"port" yaml:"port"`
	DbUrl    string `json:"dbUrl" yaml:"db_url"`
	ProxyUrl string `json:"proxyUrl" yaml:"proxy_url"`
}

var config Config

func NewDefaultConfig() Config {
	portStr := GetValueFromEnv("FIBER_GPT_SERVER_PORT", "3000")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 3000
	}
	return Config{
		Port:  port,
		DbUrl: GetValueFromEnv("FIBER_GPT_DB_URL", "./db.sqlite"),
	}
}
func Setup(c Config) {
	config = c
}
func SetupFormYaml() {
	filePath := GetValueFromEnv("FIBER_GPT_CONFIG_PATH", "./config.yaml")
	byte, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		// create default file
		byte, err = yaml.Marshal(NewDefaultConfig())
		if err != nil {
			panic(err)
		}
		_, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(filePath, byte, 0644)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(byte, &config)
	if err != nil {
		panic(err)
	}
}
func GetConfig() Config {
	return config
}
