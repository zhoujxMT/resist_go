package handle

import (
	"log"
)

// 游戏主流程
func ResistGameHandle(room *Room, msg *Message, cli *Client) {
	switch msg.EventName {
	// 聊天
	case "CHAT":
		chatMsg := &Message{From: cli.Name, Body: msg.Body, EventName: "CHAT"}
		chatMsg.UserInfo.NickName = cli.UserInfo.NickName
		chatMsg.UserInfo.AvatarURL = cli.UserInfo.AvatarURL
		room.BroadcastMessage(chatMsg, cli)
	// 演讲
	case "SPEECH":
		chatMsg := &Message{From: cli.Name, Body: msg.Body, EventName: "SPEECH"}
		chatMsg.UserInfo.NickName = cli.UserInfo.NickName
		chatMsg.UserInfo.AvatarURL = cli.UserInfo.AvatarURL
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
			speechCli := room.GetClientByName(speecher)
			choiceNextSpeecher.UserInfo.AvatarURL = speechCli.UserInfo.AvatarURL
			choiceNextSpeecher.UserInfo.NickName = speechCli.UserInfo.NickName
			room.BroadcastAll(choiceNextSpeecher)
		}
	// 队长投票
	case "TEAMVOTE":
		if msg.Body == "True" {
			room.VoteAgreeVote(cli.Name)
		} else {
			room.VoteDisVote(cli.Name)
		}
		agrvotes, disvotes := room.GetVotes()
		if agrvotes+disvotes == room.RoomSize {
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
						teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", TeamSize: teamSize}
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
					teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", TeamSize: teamSize}
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
					teamMsg := &Message{From: "SYSTEM", EventName: "TEAM", TeamSize: teamSize}
					room.SendMessage(teamMsg, capName)
				}
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
	// 获取用户信息
	case "GET_USERLIST":
		userList := []RoomUserInfo{}
		room.Lock()
		for _, cli := range room.Clients {
			roomUserInfo := RoomUserInfo{cli.Name, cli.UserInfo.NickName, cli.UserInfo.AvatarURL}
			userList = append(userList, roomUserInfo)
		}
		room.Unlock()
		msg := &Message{From: "SYSTEM", EventName: "GET_USERLIST"}
		msg.UserList = userList
		room.SendMessage(msg, cli.Name)
	// 队长选择的队员
	case "TEAM":
		for teampoint := range msg.TeamList {
			tMsg := &Message{From: cli.Name, EventName: "TEAM", Body: msg.TeamList[teampoint]}
			room.SendMessage(tMsg, msg.TeamList[teampoint])
		}
		bTeamMsg := &Message{From: cli.Name, EventName: "TEAM", TeamList: msg.TeamList}
		room.BroadcastAll(bTeamMsg)
	}
}
