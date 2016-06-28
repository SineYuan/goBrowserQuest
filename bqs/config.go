package bqs

import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Host        string
	Port        int
	DebugLevel  string `json:"debug_level"`
	MapFilePath string `json:"map_filepath"`
}

func LoadConf(confPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var c = Config{}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
