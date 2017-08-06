interface GamePage extends IPage {

}
interface GameData {
    initAnimationData: any
    initMsg: string
    showTeamAnimationData: any
}

class GamePage {
    showMsgAnimation: wx.Animation
    public data: GameData = {
        initAnimationData: {},
        initMsg: "天黑请闭眼",
        showTeamAnimationData: {}
    }
    public onShow(): void {
        var animation = wx.createAnimation({
            duration: 1500,
            timingFunction: 'ease-in',
        })
        // animation.opacity(1).scale(2,2).step()
        // this.setData({
        //     initMsg:"你的身份是：间谍",
        // })
        // this.setData({
        //     initAnimationData:animation.export()
        // }
        // )

    }

}
function setShowMsgAnimation(s:number){
    return new Promise((resolve)=>{
        setTimeout(resolve, 1500);
    })
}
Page(new GamePage());
