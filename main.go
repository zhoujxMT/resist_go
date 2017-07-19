package main

import (
	"resist_go/conf"
	"resist_go/db"
	"resist_go/handle"
	"resist_go/middleware"

	"github.com/go-martini/martini"
)

func main() {

	m := martini.Classic()
	config := conf.CreateConfig("./config/config.yaml")
	ConfigMartini(m, config)
	RouterConfig(m)
	m.Run()
}

func ConfigMartini(m *martini.ClassicMartini, config *conf.Config) *martini.ClassicMartini {
	orm := db.SetEngine(config.DataBase.DbPath)
	orm.Sync(new(db.User))
	sessionManager := middleware.GetSessionManager(7200)
	// 配置DATABASES
	m.Map(orm)
	// 全局配置信息
	m.Map(config)
	// 全局Wxssion管理器
	m.Map(sessionManager)

	return m
}

func RouterConfig(m *martini.ClassicMartini) {
	m.Get("/", func() string {
		return "hello,word"
	})
	m.Post("/login", handle.LoginWechatUser)
	m.Post("/registerUser", handle.RegisterWechatUser)
}
