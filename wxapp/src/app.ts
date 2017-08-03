import { WeApp } from './common/common'
import {Promise}from './plugin/es6-promise.js'

function _wxLogin() {
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

function _loginServer(code: string) {
    return new Promise((resolve: (res) => void, reject: (res) => void) => {
        wx.request({
            url: "http://127.0.0.1:3000/login",
            method: "POST",
            data: {
                code: code
            },
            success: (res) => {
                if (res.statusCode == 200){
                    resolve(res)
                }else{
                    reject(res)
                }
            },
        })
    })
}

function _registerServer(userRes: wx.GetUserInfoResult, thirdKey: string) {
    return new Promise((resolve, reject) => {
        wx.request({
            url: "http://127.0.0.1:3000/registerUser",
            method: "POST",
            data: {
                iv: userRes.iv,
                encryptedData: userRes.encryptedData,
                thirdKey: thirdKey
            },
            success: res => {
                resolve(res)
            },
            fail: (res) => {
                reject(res)
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
            userRes: null
        }
    }

    public onLaunch(): void {
        //调用API从本地缓存中获取数据
        const logs = wx.getStorageSync('logs') || [];
        logs.unshift(Date.now());
        wx.setStorageSync('logs', logs);
        this.loginResistServer()
    }





    public loginResistServer(): void {
        var temp = {
            code: null,
            userRes:null
        }
        _wxLogin().then((code: string) => {
            temp.code = code
            return _promiseGetUserInfo()
        }).then((res: wx.GetUserInfoResult) => {
            this.globalData.userInfo = res.userInfo
            temp.userRes = res
            return _loginServer(temp.code)
        }).then((res:wx.RequestResult)=>{
            this.globalData.thirdKey = res.data.thirdKey
        },(res:wx.RequestResult) =>{
            console.log(temp.userRes)
            this.globalData.thirdKey = res.data.thirdKey
            console.log(res.data)
            return _registerServer(temp.userRes, res.data.thirdKey)
        })

    }

    public getUserInfo(cb: (info: wx.IData) => void): void {
        if (this.globalData.userInfo) {
            typeof cb == "function" && cb(this.globalData.userInfo);
        }
    }
}

App(new WeAppClass());
