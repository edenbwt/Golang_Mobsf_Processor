package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Domain   string   `json:"domain"`
	ApiKey   string   `json:"apiKey"`
	InFolder string   `json:"InFolder"`
	Dump     string   `json:"dump"`
	Etx_list []string `json:"ext"`
	Version  string   `json:"version"`
}

func LoadConfig(filename string) (*Config, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
