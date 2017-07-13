package db

import (
	"time"
)

type User struct {
	id         int       `xorm:"int(11) pk not null autoincr"`
	opnID      string    `xorm:"varhar(100) not null index"` // 小程序内的唯一id
	nickName   string    `xorm:"varchar(50) not null"`       // 用户昵称
	avatarURL  string    `xorm:"varchar(50) not null"`       // 用户头像
	gender     string    `xorm:"varchar(50) not null"`       // 用户性别
	cretedDate time.Time `xorm:"datetime created"`           // 创建时间
	updateDate time.Time `xorm:"datetime updated"`           // 上次更新用户信息的时间
}
