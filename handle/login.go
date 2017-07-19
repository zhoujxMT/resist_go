package handle

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"

	"fmt"
	"log"
	"resist_go/conf"
	"resist_go/db"
	"resist_go/middleware"
	"time"
)

type WxCode struct {
	code string `json:"code"`
}

type WxSessionKey struct {
	openID     string `json:"openid"`
	sessionKey string `json:"session_key"`
}

type JsonUserInfo struct {
	openID    string `json:"openID"`
	nickName  string `json:"nickName"`
	avatarURL string `json:"avatarUrl"`
	gender    string `json:"gender"`
}

func LoginWechatUser(req *http.Request, config *conf.Config, session *middleware.WxSessionManager) (int, string) {
	decoder := json.NewDecoder(req.Body)
	var wxcode WxCode
	var user db.User
	err := decoder.Decode(&wxcode)
	if err != nil {
		log.Fatal(err)
	}
	// 从微信官方后台中获取用户的sessionkey
	wxsessionKey := getWxSessionCode(&wxcode, config)
	// 根据获取到的openID来读取数据库中我们的用户信息
	u, has := user.GetUserByOpenId(wxsessionKey.openID)
	if has == true {
		now := time.Now()
		userUpdateDate := u.UpdateDate
		subTime := now.Sub(userUpdateDate)
		days := subTime.Hours() / 24
		// 如果大于七天，则让小程序重新拉取用户信息进行更新
		if days > 7 {
			rspStr := fmt.Sprintf("{'error':'userinfo need to update'}")
			return 404, rspStr
		}
		// 将用户信息存放在session当中,并返回第三方sessionkey，防止官方session在网络中传输
		thirdKey := createThirdPatyKey(wxsessionKey, u, session)
		rspStr := fmt.Sprintf("{'thirdKey':'%s'}", thirdKey)
		return 200, rspStr
	} else {
		rspStr := fmt.Sprintf("{}")
		return 404, rspStr
	}
}

func getWxSessionCode(wxcode *WxCode, config *conf.Config) *WxSessionKey {
	wxSessionAddr := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.Wechat.APPID, config.Wechat.AppSecret, wxcode.code)
	rsp, err := http.Get(wxSessionAddr)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(rsp.Body)
	var wxsessionKey WxSessionKey
	err = decoder.Decode(&wxsessionKey)
	if err != nil {
		log.Fatal(err)
	}
	return &wxsessionKey
}

func createThirdPatyKey(wxsessionKey *WxSessionKey, u *db.User, session *middleware.WxSessionManager) string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Fatal(err)
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	session.Set(uuid, "userInfo", u)
	session.Set(uuid, "wxsessionKey", wxsessionKey)
	return uuid
}
