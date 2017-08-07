package handle

import (
	"encoding/json"
	"fmt"
	"log"
)

type SpeechMsg struct {
	Message string `json:"message"`
}
type TeamListMsg struct {
	TeamList []string `json:"teamList"`
}

// 游戏主流程
func ResistGameHandle(room *Room, msg *Message, cli *Client) {
	switch msg.EventName {
	// 聊天
	case "CHAT":
		chatMsg := &Message{From: cli.Name, EventName: "CHAT"}
		body := &SpeechMsg{Message: msg.Body}
		b, _ := json.Marshal(body)
		chatMsg.Body = string(b)
		room.Lock()
		room.BroadcastMessage(chatMsg, cli)
		room.Unlock()
	// 演讲
	case "SPEECH":
		chatMsg := &Message{From: cli.Name, EventName: "SPEECH"}
		body := &SpeechMsg{Message: msg.Body}
		b, _ := json.Marshal(body)
		chatMsg.Body = string(b)
		room.BroadcastMessage(chatMsg, cli)
		speecher := room.TakeTurnsClientName()
		// 选择投票任务是否选择
		if speecher == "END" {
			if room.GameNum <= 4 {
				speechEndMsg := &Message{From: "SYSTEM", EventName: "CAPVOTE"}
				room.BroadcastAll(speechEndMsg)
			}

		} else {
			// 选取下一位演讲者
			choiceNextSpeecher := &Message{From: "SYSTEM", EventName: "SPEECHER"}
			speecherMsg := map[string]string{
				"speecher": speecher,
			}
			body, _ := json.Marshal(speecherMsg)
			choiceNextSpeecher.Body = string(body)
			room.BroadcastAll(choiceNextSpeecher)
		}
	// 队长投票
	// TODO:这里要重构
	case "CAPVOTE":
		// 投票
		// 当前局数
		fmt.Println("当前局数:", room.GameNum)
		if room.GameNum <= 4 {
			if msg.Body == "True" {
				room.VoteAgreeVote(cli.Name)
			} else {
				room.VoteDisVote(cli.Name)
			}
			agrvotes, disvotes := room.GetVotes()
			fmt.Println(agrvotes, disvotes)
			// 当投票够了时候
			if agrvotes+disvotes == room.RoomSize {
				agr, _ := room.CountVote("team")
				//如果同意通知所有人可以执行任务了
				if agr == true {
					voteMsg := &Message{From: "SYSTEM", EventName: "TEAMVOTE", Body: "AGREE"}
					room.BroadcastAll(voteMsg)
					// 清空票仓
				} else {
					// 如果任务没有成功执行，则坏人赢一局,进行下一局
					if room.GameNum < 4 {
						room.AddBadGuysWins()
						voteMsg := &Message{From: "SYSTEM", EventName: "TEAMVOTE", Body: "DISMISS"}
						room.BroadcastAll(voteMsg)
						capName, _ := room.TakeRandCaptains()
						capMsg := &Message{From: "SYSTEM", EventName: "CHOICE_CAPTAIN", Body: capName}
						room.BroadcastAll(capMsg)
						room.AddGameNum()
						teamSize := room.GetMissionConfig()
						teamSizeMap := map[string]int{
							"teamSize": teamSize,
						}
						body, _ := json.Marshal(teamSizeMap)
						teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", Body: string(body)}
						room.BroadcastAll(teamMsg)
						// 清空票仓
					}
				}
				room.ClearVoteSet()
				// 当坏人已经赢过3局，游戏直接失败
				fmt.Println("坏蛋分数", room.BadGuysWins, "好人分数", room.GoodManWins)
				if room.BadGuysWins >= 3 {
					var gameMsg *Message
					gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "BADMANWINS"}
					room.BroadcastAll(gameMsg)
				}
			}
			//清空票仓
		} else {
			log.Fatal("GAMENUM ERROR")
		}
	// 任务投票阶段
	case "MISSIONVOTE":
		if room.GameNum <= 4 {
			//先投票
			if msg.Body == "True" {
				room.VoteAgreeVote(cli.Name)
			} else {
				room.VoteDisVote(cli.Name)
			}
			// 都已经投了
			missionNum := room.GetMissionConfig()
			agrvotes, disvotes := room.GetVotes()
			fmt.Println(missionNum, agrvotes, disvotes)
			if agrvotes+disvotes == missionNum {
				agr, _ := room.CountVote("mission")
				missionsHandle(room, msg, cli, agr)
			}
			if room.BadGuysWins >= 3 {
				var gameMsg *Message
				gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "BADMANWINS"}
				room.BroadcastAll(gameMsg)
			}
			if room.GoodManWins >= 3 {
				var gameMsg *Message
				gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "GOODMANWINS"}
				room.BroadcastAll(gameMsg)
			}
		} else {
			log.Fatal("GameNum Error")
		}
	// 队长选择的队员
	case "TEAM":
		bTeamMsg := &Message{From: cli.Name, EventName: "CHOICE_TEAM", Body: msg.Body}
		room.BroadcastMessage(bTeamMsg, cli)
	case "GET_USERS":
		room.Lock()
		roomUserList := []RoomUserInfo{}
		for cliName, cli := range room.Clients {
			roomUserInfo := RoomUserInfo{Name: cliName,
				NickName:  cli.UserInfo.NickName,
				AvatarUrl: cli.UserInfo.AvatarURL}
			roomUserList = append(roomUserList, roomUserInfo)
		}
		jsonUserInfo, _ := json.Marshal(roomUserList)
		userListMsg := &Message{From: "SYSTEM", EventName: "GET_USERS", Body: string(jsonUserInfo)}
		room.Unlock()
		room.SendMessage(userListMsg, cli.Name)

	default:
		fmt.Println("哇哈哈")
	}
}

func missionsHandle(room *Room, msg *Message, cli *Client, isAgree bool) {
	if isAgree == true {
		// 任务成功
		voteMsg := &Message{From: "SYSTEM", EventName: "MISSIONVOTE", Body: "AGREE"}
		room.BroadcastAll(voteMsg)
		room.AddGoodManWins()

	} else {
		// 任务失败
		voteMsg := &Message{From: "SYSTEM", EventName: "MISSIONVOTE", Body: "DISSMISS"}
		room.BroadcastAll(voteMsg)
		room.AddBadGuysWins()
	}
	// 下一局
	if room.GameNum < 4 {

		capName, _ := room.TakeRandCaptains()
		capMsg := &Message{From: "SYSTEM", EventName: "CHOICE_CAPTAIN", Body: capName}
		room.BroadcastAll(capMsg)
		teamSize := room.GetMissionConfig()
		room.AddGameNum()
		teamSizeMap := map[string]int{
			"teamSize": teamSize,
		}
		body, _ := json.Marshal(teamSizeMap)
		teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", Body: string(body)}
		room.BroadcastAll(teamMsg)
		// 清空票仓
		room.ClearVoteSet()
	}
}
