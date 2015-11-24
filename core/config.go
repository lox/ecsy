package core

import (
	"encoding/json"
	"os"
)

type Config struct {
	Users []string
}

func LoadConfig() (conf Config, err error) {
	file, err := os.Open("./conf.json")
	if err != nil {
		return conf, err
	}
	return conf, json.NewDecoder(file).Decode(&conf)
}
