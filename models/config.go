package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DownloadPath      string   `json:"downloadPath" `
	DownloadThread    int32    `json:"downloadThread"`
	DownloadBlockSize int64    `json:"downloadBlockSize"`
	Tokens            []string `json:"tokens"`
}

var (
	DefaultConf *Config
)

func LoadConfig(file string) *Config {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("failed to read config: " + err.Error())
	}
	var conf Config
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		log.Fatalln("failed to decode config: " + err.Error())
	}
	return &conf
}

func (conf *Config) SaveFile(file string) {
	bytes, err := json.Marshal(conf)
	if err != nil {
		log.Fatalln("failed to encode config: " + err.Error())
	}
	err = ioutil.WriteFile(file, bytes, 0666)
	if err != nil {
		log.Fatalln("failed to write file: " + err.Error())
	}
}
