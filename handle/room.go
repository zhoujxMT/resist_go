package handle

import (
	"sync"
)

type Room struct {
	Lock        sync.Mutex         // 互斥锁，保证线程安全
	RoomSize    int                // 房间人数
	Name        string             // 创建房间时的名字，创建时为uuid，并分享时候将该uuid带上
	GameNum     int                // 房间当前局数
	GoodManWins int                // 抵抗组织成员获胜局数
	BadGuysWins int                // 间谍成员获胜局数
	DisVote     *VoteSet           // 反对票仓
	AgrVote     *VoteSet           // 赞成票仓
	Clients     map[string]*Client // 客户端管理池
}

func CreteRoom(roomName string, roomSize int) *Room {
	dismissVote := NewVote() // 反对票仓
	agreeVote := NewVote()   // 同意票仓
	room := Room{Lock: sync.Mutex{}, Name: roomName, DisVote: dismissVote, AgrVote: agreeVote, Clients: map[string]*Client{}, RoomSize: roomSize}
	return &room
}

// 添加房间客户端
func (room *Room) AddClient(clientName string, client *Client) bool {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.RoomSize < 5 {
		room.Clients[clientName] = client
		return true
	} else {
		return false
	}
}

// 清空房间票仓
func (room *Room) ClearVoteSet() {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	room.DisVote.Clear()
	room.AgrVote.Clear()
}

// 投同意票
func (room *Room) VoteAgreeVote(clientName string) {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	room.AgrVote.Add(clientName)
}

// 投反对票
func (room *Room) VoteDisVote(clientName string) {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	room.DisVote.Add(clientName)
}

// 投机投票数，参数为模式，mission(任务执行模式)|team(组队模式) 返回如果味true则同意多，如果为fale则反对多 统计结束后则晴空票仓,
func (room *Room) CountVote(modle string) (bool, bool) {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	agrvotes := room.AgrVote.Len()
	disvotes := room.DisVote.Len()
	if modle == "mission" {
		if disvotes >= 1 {
			return false, true
		} else {
			return true, true
		}
	} else if modle == "team" {
		if agrvotes > disvotes {
			return true, true
		} else {
			return false, true
		}
	} else {
		return false, false
	}
}

// 增加局数
func (room *Room) AddGameNum() bool {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.GameNum < 5 {
		room.GameNum++
		return true
	} else {
		return false
	}
}

// 增加好人获胜局数
func (room *Room) AddGoodManWins() bool {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.GoodManWins < 3 && room.GoodManWins+room.BadGuysWins <= 5 {
		room.GameNum++
		return true
	} else {
		return false
	}
}

// 增加坏人获胜局数
func (room *Room) AddBadGuysWins() bool {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	if room.BadGuysWins < 3 && room.GoodManWins+room.BadGuysWins <= 5 {
		room.BadGuysWins++
		return true
	} else {
		return false
	}
}

// 客户端名字列表
func (room *Room) ClientNameList() []string {
	room.Lock.Lock()
	defer room.Lock.Unlock()
	list := []string{}
	for clientName := range room.Clients {
		list = append(list, clientName)
	}
	return list
}
