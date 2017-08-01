package db

import (
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var orm *xorm.Engine

func SetEngine(dbPath string) *xorm.Engine {
	var err error
	orm, err = xorm.NewEngine("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	orm.TZLocation = time.Local
	return orm
}
