interface GamePage extends IPage {

}
import { WeApp } from '../../common/common'

const app: WeApp = getApp() as WeApp;

interface RoomUserInfo {
    nickName: string
    avatarUrl: string
}

interface GameData {
    waitShow: boolean
    initAnimationData: any
    initMsg: string
    initMsgShow: boolean
    showTeamAnimationData: any
    initViewAnimationData: any
    waitData: any
    roomUsers: any[] //房间用户列表
    teamList: any[] // 当你为坏人时候你的队友是谁
    isBadMan: boolean
}

class GamePage implements GamePage {
    userInfoCache: { [name: string]: any } = {}
    gameInfoCache: any = {}
    public data: GameData = {
        waitShow: false,
        initAnimationData: {},
        initMsg: "天黑请闭眼",
        initMsgShow: false,
        showTeamAnimationData: {},
        initViewAnimationData: {},
        waitData: {},
        roomUsers: [],
        teamList: [],
        isBadMan: false
    }

    public onLoad(): void {

    }


    public onShow(): void {
        var thirdKey: string = app.getUserThirdKey()
        this.wechatSockets(`ws://localhost:3000/game/testroom/test?thirdKey=${thirdKey}&isWeChat=true`)
    }


    private wechatSockets(wss: string): void {
        wx.connectSocket({
            url: wss,
        })
        wx.onSocketOpen((res) => {
            console.log("已打开")
        })
        wx.onSocketError((res) => {
            console.log("啊哦,服务器提出了一个问题")
            wx.redirectTo({
                url: "../index/index"
            })
        })
        wx.onSocketMessage((res) => {
            let jsonData = JSON.parse(res.data)
            console.log(jsonData)
            this._handleSocketMsg(jsonData)
        })
    }

    private _handleSocketMsg(msg) {
        switch (msg.eventName) {
            case "GET_USERS":
                let userInfoList = JSON.parse(msg.body)
                this.setData({
                    roomUsers: userInfoList
                })
                for (let u of userInfoList) {
                    this.userInfoCache[u.name] = u
                }
                console.log(this.userInfoCache)
                break;
            case "JOIN":
                // 加入房间
                let userInfo = JSON.parse(msg.body)
                console.log(userInfo)
                let users = this.data.roomUsers
                users.push(userInfo)
                this.setData({
                    roomUsers: users
                })
                this.userInfoCache[userInfo.name] = userInfo
                console.log(this.userInfoCache)
                break;
            case "READY":
                this.waitFinashAnimition()
                this.initAnimition()
                break;
            case "INIT":
                this._onInit(msg)

            default:
                break;
        }
    }

    private _onInit(msg: any): void {
        let initInfo = JSON.parse(msg.body)
        this.gameInfoCache.captain = initInfo.captain
        this.gameInfoCache.teamList = initInfo.teamList
        this.gameInfoCache.role = initInfo.role
        this.gameInfoCache.speecher = initInfo.speecher
        let badManteamList: any[] = []
        if (initInfo.role == "BADMAN") {
            for (var teamName of initInfo.teamList) {
                let teamUser = this.userInfoCache[teamName]
                console.log(teamUser)
                badManteamList.push({
                    nickName: teamUser.nickName,
                    avatarUrl: teamUser.avatarUrl
                })
            }
        }
        let roleIsBandman: boolean
        if (initInfo.role == "BADMAN") {
            roleIsBandman = true
        } else {
            roleIsBandman = false
        }
        this.setData({
            teamList: badManteamList,
            isBadMan: roleIsBandman
        })
    }
    private waitFinashAnimition(): void {
        var waitAnimation = wx.createAnimation({
            duration: 1500,
            timingFunction: "ease-in"
        })
        waitAnimation.rotate(180).opacity(0).step()
        this.setData({
            waitData: waitAnimation.export()
        })
        promiseAnimition(1500).then(() => {
            this.setData({
                waitShow: true
            })
        })

    }

    private initAnimition(): void {
        // 初始提示信息的动画
        var animation = wx.createAnimation({
            duration: 3000,
            timingFunction: 'ease-in',
        })
        // 初始Team的动画
        var initViewAnimation = wx.createAnimation({
            duration: 3000,
            timingFunction: 'ease-in'
        })
        // 队友Team动画
        var teamAnimation = wx.createAnimation({
            duration: 3000,
            timingFunction: 'ease-in'
        })
        // 初始化动画集合
        animation.opacity(1).step()
        this.setData({
            initAnimationData: animation.export()
        })
        promiseAnimition(3000).then(() => {
            animation.opacity(0).step()
            this.setData({
                initAnimationData: animation.export()
            })
            let role: string = ""
            if (this.gameInfoCache.role == "GOODMAN") {
                role = "抵抗组织"
            } else {
                role = "间谍"
            }
            return promiseAnimition(3000, `你的身份是：${role}`)
        }).then((data) => {
            this.setData({
                initMsg: data
            })
            animation.opacity(1).step()
            this.setData({
                initAnimationData: animation.export()
            })
            return promiseAnimition(3000)
        }).then((data) => {
            animation.opacity(0).step()
            this.setData({
                initAnimationData: animation.export()
            })
            teamAnimation.opacity(1).step()
            this.setData({
                showTeamAnimationData: teamAnimation.export()
            })
            return promiseAnimition(3000)
        }).then(() => {
            initViewAnimation.opacity(0).step()
            this.setData({
                initViewAnimationData: initViewAnimation.export()
            })
            return promiseAnimition(3000)
        }).then(() => {
            this.setData({
                initMsgShow: true,
                waitShow: true
            })
        })

    }

}

function promiseAnimition(timeout: number, setData?: any) {
    return new Promise((resolve, reject) => {
        setTimeout(function () {
            resolve(setData)
        }, timeout);
    })
}

Page(new GamePage());
