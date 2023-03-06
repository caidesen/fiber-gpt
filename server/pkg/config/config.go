package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port  int    `json:"port" yaml:"port"`
	DbUrl string `json:"dbUrl" yaml:"db_url"`
}

var config Config

func NewDefaultConfig() Config {
	return Config{
		Port:  3000,
		DbUrl: "data.db",
	}
}
func Setup(c Config) {
	config = c
}
func SetupFormYaml(filePath string) {
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
