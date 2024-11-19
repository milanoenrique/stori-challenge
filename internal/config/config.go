package config

import (
	"encoding/json"
	"os"
)

type Database struct {
	Host string
	Database string
	Password string
	Port string
	User string
}

type EmailSender struct {
	AccountId string
	Token string
}

type Config struct {
	Database Database
	EmailSender EmailSender
}

func LoadConfig() (*Config, error){
	confFile, err :=os.ReadFile("./conf.json")

	if err != nil {
		return nil, err
	}
	conf := new(Config)
	err = json.Unmarshal(confFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}