import { WeApp } from './common/common'
import * as Promise from './plugin/bluebird.js'
function _loginServer(code: string) {
        return new Promise((resolve, reject) => {
            wx.request({
                url: "https://yuchenyang1994.natapp4.cc/login",
                method: "POST",
                data: {
                    code: code
                },
                success: res => {
                        resolve(res)
                    }
            })
        })
    }

 function _registerServer(userRes: wx.GetUserInfoResult, thirdKey: string) {
        return new Promise((resolve, reject) => {
            wx.request({
                url: "https://yuchenyang1994.natapp4.cc/registerUser",
                method: "POST",
                data: {
                    iv: userRes.iv,
                    encryptedData: userRes.encryptData,
                    thirdKey: thirdKey
                },
                success: res => {
                    resolve(res)
                }
            })
        })
    }

 function _promiseGetUserInfo() {
        return new Promise((resolve, reject) => {
            wx.getUserInfo({
                withCredentials: true,
                success: res => {
                    resolve(res)
                }
            })
        })
    }

class WeAppClass implements WeApp {

    public get globalData() {
        return {
            userInfo: null,
            thirdKey: null,
            userRes:null
        }
    }

    public onLaunch(): void {
        //调用API从本地缓存中获取数据
        const logs = wx.getStorageSync('logs') || [];
        logs.unshift(Date.now());
        wx.setStorageSync('logs', logs);
    }

    private _wxLogin() {
        return new Promise((resolve, reject) => {
            wx.login({
                success: (res) => {
                    if (res.code) {
                        resolve(res.code)
                    } else {
                        reject(res.errMsg)
                    }
                }
            })
        })
    }

    

    public loginResistServer(): void {
        this._wxLogin().then((code: string) => {
            return _loginServer(code)
        }).then((res:wx.RequestResult) => {
            console.log("?")
            if (res.statusCode == 200){
                this.globalData.thirdKey = res.data.thirdKey
            }else{
                return _promiseGetUserInfo()
            }
        }).then((res:wx.GetUserInfoResult) => {
            console.log(3)
            this.globalData.userInfo=res.userInfo
            return _registerServer(res,this.globalData.thirdKey)
        }).then((res:wx.RequestResult)=>{
            console.log(res)
        }, (res:wx.RequestResult)=>{
            console.log(res)
        })

    }

    public getUserInfo(cb: (info: wx.IData) => void): void {
        if (this.globalData.userInfo) {
            typeof cb == "function" && cb(this.globalData.userInfo);
        } else {
            this.loginResistServer()
        }
    }
}

App(new WeAppClass());
