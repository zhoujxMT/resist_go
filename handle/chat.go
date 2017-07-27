package handle

import (
	"sync"

	"github.com/go-martini/martini"
)

type Chat struct {
	sync.Mutex
	Rooms map[string]*Room
}

var chat *Chat
var once sync.Once

// NewChat ...
func GetChat() *Chat {
	once.Do(func() {
		chat = &Chat{sync.Mutex{}, map[string]*Room{}}
	})
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

// HandleGameSocket ...
func HandleGameSocket(params martini.Params, recevier <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
}
