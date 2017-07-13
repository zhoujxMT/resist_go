package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataBase struct {
		DbPath string `yaml:"db_path"`
	} `yaml:"db"`

	Wechat struct {
		APPID     string `yaml:"appid"`
		AppSecret string `yaml:"appsecret"`
	} `yaml:wechat`
}

func CreateConfig(configPath string) *Config {
	config := Config{}
	config_string, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(config_string), &config)

	return &config
}
