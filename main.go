package main

import (
	"resist_go/conf"

	"resist_go/db"

	"github.com/go-martini/martini"
)

func main() {

	m := martini.Classic()
	config := conf.CreateConfig("./config/config.yaml")
	ConfigMartini(m, config)
	m.Get("/", func() string {
		return "hello,word"
	})
	m.Run()
}

func ConfigMartini(m *martini.ClassicMartini, config *conf.Config) *martini.ClassicMartini {
	orm := db.SetEngine(config.DataBase.DbPath)
	// 配置DATABASES
	m.Map(orm)
	// 配置微信开发者配置信息
	m.Map(config.Wechat)
	return m
}
