package handle

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"net/http"
	"resist_go/db"
	"resist_go/middleware"

	"github.com/go-martini/martini"
)

type Chat struct {
	sync.Mutex
	Rooms map[string]*Room
}

type CreateRoomInfo struct {
	RoomSize int `json:"roomSize"`
}
type ChatUserInfo struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
}

var chat *Chat
var once sync.Once

// NewChat ...
func GetChat() *Chat {
	if chat == nil {
		once.Do(func() {
			chat = &Chat{sync.Mutex{}, map[string]*Room{}}
		})
	}
	return chat
}

// (chat *Chat) GetRoomByName ...
func (chat *Chat) GetRoomByName(roomName string) *Room {
	chat.Lock()
	defer chat.Unlock()
	room, ok := chat.Rooms[roomName]
	if ok == true {
		return room
	} else {
		return nil
	}
}

// AddRoom ...
func (chat *Chat) AddRoom(roomName string, room *Room) {
	chat.Lock()
	defer chat.Unlock()
	chat.Rooms[roomName] = room
}

// RemoveRoom
func (chat *Chat) RemoveChat(roomName string) {
	delete(chat.Rooms, roomName)
}

func HandleCreateRoom(req *http.Request) (int, string) {
	var createRoomInfo CreateRoomInfo
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&createRoomInfo)
	if err != nil {
		log.Fatal(err)
	}
	// 生成uuid
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Fatal(err)
		return 0, ""
	}
	uuid := fmt.Sprintf("room_%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	room := CreteRoom(uuid, createRoomInfo.RoomSize)
	chat.AddRoom(uuid, room)
	roomInfo := fmt.Sprintf(`{"roomId":%s}`, uuid)
	return 200, roomInfo
}

// HandleGameSocket ...
func ResistSocket(req *http.Request, params martini.Params, recevier <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error, session *middleware.WxSessionManager) (int, string) {
	thridKey := req.URL.Query().Get("thirdKey")
	roomName := params["name"]
	info, ok := session.Get(thridKey, "userInfo")
	// 添加测试账号
	if ok == true {
		userInfo := info.(*db.User)
		cli := Client{Name: thridKey, UserInfo: userInfo, in: recevier, out: sender, done: done, err: err, diconnect: disconnect}
		room := chat.GetRoomByName(roomName)
		if room == nil {
			return 404, "{errorInfo:'can't find room'}"
		}
		room.AddClient(thridKey, &cli)
		addMsg := &Message{From: thridKey, EventName: "JOIN", Body: ""}
		room.BroadcastMessage(addMsg, &cli)
		for {
			select {
			case errMsg := <-cli.err:
				log.Fatal(errMsg)
				// 处理错误消息
			case msg := <-cli.in:
				// 消息处理器
				fmt.Println(msg)
				ResistGameHandle(room, msg, &cli)
			case <-cli.done:
				// 处理掉线
				room.RemoveClient(cli.Name)
				msg := &Message{From: thridKey, EventName: "DISCONTENT"}
				chatUserInfo := &ChatUserInfo{cli.UserInfo.NickName, cli.UserInfo.AvatarURL}
				body, _ := json.Marshal(chatUserInfo)
				msg.Body = string(body)
				room.BroadcastMessage(msg, &cli)
				if len(room.ClientNameList()) == 0 {
					chat.RemoveChat(roomName)
				}
				return 200, "ok"
			}
		}

	} else {
		return 403, "{'errorInfo':'no login'}"
	}
}

func ResistSocketTest(req *http.Request, params martini.Params, recevier <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error, session *middleware.WxSessionManager) (int, string) {
	thridKey := req.URL.Query().Get("thirdKey")
	isWechat := req.URL.Query().Get("isWeChat")
	roomName := params["name"]
	// 添加测试账号
	var cli Client
	var userInfo *db.User
	if isWechat == "true" {
		u, _ := session.Get(thridKey, "userInfo")
		userInfo = u.(*db.User)
	} else {
		userInfo = &db.User{OpenID: "TEST", NickName: "Test",
			AvatarURL: "http://img.qq745.com/uploads/allimg/160630/8-160630152935.jpg", Gender: "男"}
	}
	cli = Client{Name: thridKey, UserInfo: userInfo, in: recevier, out: sender, done: done, err: err, diconnect: disconnect}
	room := chat.GetRoomByName(roomName)
	if room == nil {
		return 404, "{errorInfo:'can't find room'}"
	}
	room.AddClient(thridKey, &cli)
	for {
		select {
		case <-cli.err:
			// 处理错误消息
		case msg := <-cli.in:
			// 消息处理器
			ResistGameHandle(room, msg, &cli)
		case <-cli.done:
			// 处理掉线
			room.RemoveClient(cli.Name)
			msg := &Message{From: thridKey, EventName: "DISCONTENT"}
			chatUserInfo := &ChatUserInfo{cli.UserInfo.NickName, cli.UserInfo.AvatarURL}
			body, _ := json.Marshal(chatUserInfo)
			msg.Body = string(body)
			room.Lock()
			room.BroadcastMessage(msg, &cli)
			room.Unlock()
			if len(room.ClientNameList()) == 0 {
				chat.RemoveChat(roomName)
			}
			return 200, "ok"
		}
	}
}
