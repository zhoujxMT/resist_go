package handle

import (
	"encoding/json"
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
			if room.GameNum <= 5 {
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
	case "TEAMVOTE":
		// 投票
		if msg.Body == "True" {
			room.VoteAgreeVote(cli.Name)
		} else {
			room.VoteDisVote(cli.Name)
		}
		agrvotes, disvotes := room.GetVotes()
		// 如果票数够了
		if agrvotes+disvotes == room.RoomSize {
			// 统计投票
			agr, _ := room.CountVote("team")
			if room.GameNum <= 5 {
				if agr == true {
					voteMsg := &Message{From: "SYSTEM", EventName: "TEAMVOTE", Body: "AGREE"}
					room.BroadcastAll(voteMsg)
				} else {
					// 如果任务没有成功执行，则坏人赢一局
					room.AddBadGuysWins()
					voteMsg := &Message{From: "SYSTEM", EventName: "TEAMVOTE", Body: "DISMISS"}
					room.BroadcastAll(voteMsg)
					// ...
					if room.GameNum == 5 {
						// 如果已经进行第5局则开始清算了
						var gameMsg *Message
						if room.GoodManWins > room.BadGuysWins {
							gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "GOODMAN"}
						} else {
							gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "BADGUYS"}
						}
						room.BroadcastAll(gameMsg)

					} else {
						// 如果失败进行下一局，重新选择一个队长
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
						room.SendMessage(teamMsg, capName)
					}

				}
				// 清算后清空票仓
				room.ClearVoteSet()
			} else {
				log.Fatal("GAMENUM ERROR")
			}
		}
	// 任务投票阶段
	case "MISSIONVOTE":
		if msg.Body == "True" {
			room.VoteAgreeVote(cli.Name)
		} else {
			room.VoteDisVote(cli.Name)
		}
		missionNum := room.GetMissionConfig()
		agrvotes, disvotes := room.GetVotes()
		// 查看票输决定是否进行下一局
		if agrvotes+disvotes == missionNum {
			agr, _ := room.CountVote("mission")
			if room.GameNum <= 5 {
				if agr == true {
					voteMsg := &Message{From: "SYSTEM", EventName: "MISSIONVOTE", Body: "AGREE"}
					room.BroadcastAll(voteMsg)
					room.AddGoodManWins()
					capName, _ := room.TakeRandCaptains()
					capMsg := &Message{From: "SYSTEM", EventName: "CHOICE_CAPTAIN", Body: capName}
					room.BroadcastAll(capMsg)
					// 发送当前局数选择队员人数
					teamSize := room.GetMissionConfig()
					teamSizeMap := map[string]int{
						"teamSize": teamSize,
					}
					body, _ := json.Marshal(teamSizeMap)

					teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", Body: string(body)}
					room.SendMessage(teamMsg, capName)
				} else {
					voteMsg := &Message{From: "SYSTEM", EventName: "MISSIONVOTE", Body: "DISMISS"}
					room.BroadcastAll(voteMsg)
					room.AddBadGuysWins()
					capName, _ := room.TakeRandCaptains()
					capMsg := &Message{From: "SYSTEM", EventName: "CHOICE_CAPTAIN", Body: capName}
					room.BroadcastAll(capMsg)
					// 发送当前局数选择队员人数
					teamSize := room.GetMissionConfig()
					teamSizeMap := map[string]int{
						"teamSize": teamSize,
					}
					body, _ := json.Marshal(teamSizeMap)
					teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", Body: string(body)}
					room.SendMessage(teamMsg, capName)
				}
				// 如果局数等于5
				if room.GameNum == 5 {
					var gameMsg *Message
					if room.GoodManWins > room.BadGuysWins {
						gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "GOODMAN"}
					} else {
						gameMsg = &Message{From: "SYSTEM", EventName: "GAMEOVER", Body: "BADGUYS"}
					}
					room.BroadcastAll(gameMsg)

				} else {
					room.AddGameNum()
				}
			} else {
				log.Fatal("GameNum Error")
			}
		}
	// 队长选择的队员
	case "TEAM":
		var teamListMsg TeamListMsg
		json.Unmarshal([]byte(msg.Body), &teamListMsg)
		bTeamMsg := &Message{From: cli.Name, EventName: "CHOICE_TEAM", Body: msg.Body}
		room.BroadcastAll(bTeamMsg)
	}
}
