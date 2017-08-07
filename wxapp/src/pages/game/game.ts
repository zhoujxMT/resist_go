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
}

class GamePage implements GamePage {
    userInfoCache: { [name: string]: any } = {}
    public data: GameData = {
        waitShow: false,
        initAnimationData: {},
        initMsg: "天黑请闭眼",
        initMsgShow: false,
        showTeamAnimationData: {},
        initViewAnimationData: {},
        waitData: {},
        roomUsers: []
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
            //开始结束等待的界面
            // this.waitFinashAnimition()
            //开始初始化信息的界面
            // this.initAnimition()
            default:
                break;
        }
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
        promiseAnimition(1500).then(()=>{
            this.setData({
                waitShow:true
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
            return promiseAnimition(3000, "你的身份是：间谍")
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
                waitShow:true
            })
        })

    }

}

function promiseAnimition(timeout: number, setData?: any) {
    return new Promise((resolve, reject) => {
        setTimeout(function() {
            resolve(setData)
        }, timeout);
    })
}

Page(new GamePage());
