import { WeApp } from '../../common/common';
import {Promise} from '../../plugin/es6-promise.js'

const app: WeApp = getApp() as WeApp;

interface RoomSizeObj{
    roomSize:number
}
interface IndexPageData {
    index:number;
    rangeRooms:RoomSizeObj[];
    pickList:string[];
}

interface IndexPage extends IPage {
}

class IndexPage {
    public data:IndexPageData = {
        index:0,
        rangeRooms:[{
            roomSize:5
        },{
            roomSize:6
        },{
            roomSize:7
        }],
        pickList:["5人房","6人房","7人房"]
    }
    
    public bindPickRoomSize(e): void{
        console.log(e.detail.value)
        var roomSizeObj = this.data.rangeRooms[e.detail.value]
        wx.navigateTo({
            url:"../game/game"
        })
        // this._createRoom(roomSizeObj.roomSize).then((res)=>{
        //     wx.redirectTo({
        //         url:"../game/game"
        //     })
        // }).catch((res)=>{
        //     wx.showToast({
        //         title:"服务器提了一个问题",
        //         icon:"loading"
        //     })
        // })
    }
    public _createRoom(roomSize:number){
        return new Promise((resolve:(res:wx.RequestResult)=>void,reject:(res:wx.RequestResult)=>void)=>{
            wx.request({
                url:"http://127.0.0.1:3000/room",
                method:"POST",
                data:{
                    roomSize:roomSize
                },
                success:res=>{
                    if (res.statusCode == 200){
                        resolve(res)
                    }else{
                        reject(res)
                    }
                }
            })
        })

    }
}

Page(new IndexPage());
