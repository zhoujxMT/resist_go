# 抵抗组织狼人杀小游戏后端
## build with Glide martini go-xorm
todo someting
## How To buld
`glide install`

`glide up`

`go build`

`./resist_go`


## WebSocket EventName
```
    {
      from:"xxxx" // 从哪来
      eventName:"xxxx" // 时间名称
      body: "xxxx" // 消息内容
      roleInfo:{
          role: "角色/GOODMAN/BADGUYS"
          Captian: "队长"
      }
      userInfo:{
          nickName: "客户端对应的用户名称"
          AvatarURL: "用户头像"
      }
      UserList: [{userInfo}] // 用户列表 
      TeamSize // 组队可选人数
      TeamList // 组队任务列表，类型为cli的Name
    }
```
### 事件名称:
- `CHAT`: 聊天内容，发送的聊天信息将广播所有客户端除了自己
- `SPEECH`: 演讲内容: 发送的演讲将广播除了自己房间的所有客户端
- `TEAMVOTE`: 队伍投票: 将投出当前任务是否执行的选票
- `MISSIONVOTE`: 任务投票:投出是否成功的投票
- `GET_USERLIST`: 获取用户信息列表，这个存在小程序的localstronge里，日后有用
- `TEAM`: 队长用户选择将要参加用户人员的用户信息
- `GAMEOVER`: 游戏结束标志，游戏清算，将公布获胜结果
- `CHOICE_CAPTAIN`: 系统所发，将公布下一个队长是谁

