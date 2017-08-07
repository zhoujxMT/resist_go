interface GamePage extends IPage {

}
interface GameData {
    waitShow:boolean
    initAnimationData: any
    initMsg: string
    initMsgShow:boolean
    showTeamAnimationData: any
    initViewAnimationData:any
}

class GamePage {
    showMsgAnimation: wx.Animation
    public data: GameData = {
        waitShow:true,
        initAnimationData: {},
        initMsg: "天黑请闭眼",
        initMsgShow:false,
        showTeamAnimationData: {},
        initViewAnimationData:{}
    }
    public onShow(): void {
        this.initAnimition()

    }

    public initAnimition():void{
        // 初始提示信息的动画
        var animation = wx.createAnimation({
            duration: 3000,
            timingFunction: 'ease-in',
            })
        // 初始Team的动画
        var initViewAnimation = wx.createAnimation({
                duration:3000,
                timingFunction:'ease-in'
            })
        // 队友Team动画
        var teamAnimation = wx.createAnimation({
                duration:3000,
                timingFunction:'ease-in'
            })


        // 初始化动画集合
        animation.opacity(1).step()
        this.setData({
            initAnimationData:animation.export()
        })
        promiseAnimition(3000).then(()=>{
            animation.opacity(0).step()
            this.setData({
                initAnimationData:animation.export()
            })
            return promiseAnimition(3000,"你的身份是：间谍")
        }).then((data)=>{
            this.setData({
                initMsg:data
            })
            animation.opacity(1).step()
            this.setData({
                initAnimationData:animation.export()
            })
            return promiseAnimition(3000)
        }).then((data)=>{
            animation.opacity(0).step()
            this.setData({
                initAnimationData:animation.export()
            })
            teamAnimation.opacity(1).step()
            this.setData({
                showTeamAnimationData:teamAnimation.export()
            })
            return promiseAnimition(3000)
        }).then(()=>{
            initViewAnimation.opacity(0).step()
            this.setData({
                initViewAnimationData:initViewAnimation.export()
            })
            return promiseAnimition(3000)
        }).then(()=>{
            this.setData({
                initMsgShow:true
            })
        })
        
    }

}

function promiseAnimition(timeout:number,setData?:any){
    return new Promise((resolve, reject)=>{
        setTimeout(function() {
            resolve(setData)
        }, timeout);
    })
}

Page(new GamePage());
