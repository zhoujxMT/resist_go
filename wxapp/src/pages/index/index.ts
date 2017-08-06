import { WeApp } from '../../common/common';

const app: WeApp = getApp() as WeApp;

interface RoomSizeObj{
    id:number;
    name:string;
}
interface IndexPageData {
    index:number;
    rangeRooms:object[];
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
    }

    public bindViewTap(): void {
        wx.navigateTo({
            url: '../logs/logs'
        });
    }
}

Page(new IndexPage());
