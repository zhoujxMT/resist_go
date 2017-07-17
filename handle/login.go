package handle

import (
	"encoding/json"
	"log"
	"net/http"

	"resist_go/db"

	"github.com/go-martini/martini"
)

type RepUser struct {
	openID   string
	nickName string
}
type JsonUserInfo struct {
	openID    string `json:"openID"`
	nickName  string `json:"nickName"`
	avatarURL string `json:"avatarUrl"`
	gender    string `json:"gender"`
}

func HandleCheckUser(req *http.Request, params martini.Params) (int, string) {
	var user db.User
	u, has := user.GetUserByOpenId(params["openID"])
	rspUser := &RepUser{openID: u.OpenID, nickName: u.NickName}
	body, err := json.Marshal(rspUser)
	if err != nil {
		log.Fatalln(err)
	}
	if has == true {
		return 200, string(body)
	} else {
		return 404, "{}"
	}
}

func RegisterWechatUser(req *http.Request) (int, string) {
	var jsonUser JsonUserInfo
	decode := json.NewDecoder(req.Body)
	decode.Decode(&jsonUser)
	var user db.User
	user = db.User{OpenID: jsonUser.openID, NickName: jsonUser.nickName,
		AvatarURL: jsonUser.avatarURL, Gender: jsonUser.gender}
	ok := user.Insert()
	if ok == true {
		rspUser := &RepUser{openID: user.OpenID, nickName: user.NickName}
		body, err := json.Marshal(rspUser)
		if err != nil {
			log.Fatalln(err)
		}
		return 200, string(body)
	} else {
		return 500, "{'info':'insert error'}"
	}
}
