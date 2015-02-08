package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HostName      string
	RootDirectory string
}

func New(file *os.File) *Config {
	if file == nil {
		panic("Invalid configuration file!")
	}
	config := new(Config)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(config)
	if err != nil {
		panic(err)
	}
	return config
}
