package db

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var orm *xorm.Engine

func SetEngine(dbPath string) *xorm.Engine {
	orm, err := xorm.NewEngine("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	return orm
}
