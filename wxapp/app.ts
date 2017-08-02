import { WeApp } from './common/common'

class WeAppClass implements WeApp {

    public get globalData() {
        return {
            userInfo: null
        }
    }

    public onLaunch(): void {
        //调用API从本地缓存中获取数据
        const logs = wx.getStorageSync('logs') || [];
        logs.unshift(Date.now());
        wx.setStorageSync('logs', logs);
    }

    public getUserInfo(cb: (info: wx.IData) => void): void {
        if (this.globalData.userInfo) {
            typeof cb == "function" && cb(this.globalData.userInfo);
        } else {
            //调用登录接口
            wx.login({
                success: (res) => {
                    wx.request({
                        url:"https://yuchenyang1994.natapp4.cc/login",
                        method: "POST",
                        data: {
                            code:res.code
                        },
                        success: (success_res, statusCode?,header?) =>{
                            if (statusCode == 200) {
                                console.log("调用成功")
                            }else{
                                console.log("调用失败")
                            }
                        },
                    });
                    wx.getUserInfo({
                        success: res => {
                            this.globalData.userInfo = res.userInfo;
                            typeof cb == "function" && cb(this.globalData.userInfo)
                        }
                    })
                }
            });
        }
    }
};

App(new WeAppClass());
