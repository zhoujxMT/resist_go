package db

import (
	"log"
	"time"
)

type User struct {
	id         int       `xorm:"int(11) pk not null autoincr"`
	openID     string    `xorm:"varhar(100) not null index"` // 小程序内的唯一id
	nickName   string    `xorm:"varchar(50) not null"`       // 用户昵称
	avatarURL  string    `xorm:"varchar(50) not null"`       // 用户头像
	gender     string    `xorm:"varchar(50) not null"`       // 用户性别
	cretedDate time.Time `xorm:"datetime created"`           // 创建时间
	updateDate time.Time `xorm:"datetime updated"`           // 上次更新用户信息的时间
}

func (this *User) GetUserByOpenId(openID string) (*User, bool) {
	// 根据微信的openId获得用户已经注册的用户信息
	user := &User{openID: openID}
	has, err := orm.Get(user)
	if err != nil {
		log.Panic(err)
	}
	return user, has
}

func (this *User) Insert() bool {
	_, err := orm.InsertOne(this)
	if err != nil {
		log.Panic(err)
	}
	return true
}

func (this *User) Delete() bool {
	_, err := orm.Delete(this)
	if err != nil {
		log.Panic(err)
	}
	return true
}

func (this *User) Update() bool {
	_, err := orm.Id(this.id).Update(this)
	if err != nil {
		log.Panic(err)
	}
	return true
}
