import { WeApp } from '../../common/common';

const app: WeApp = getApp() as WeApp;

interface IndexPageData {
    motto: string;
    userInfo: wx.IData;
}

interface IndexPage extends IPage {
}

class IndexPage {

    public bindViewTap(): void {
        wx.navigateTo({
            url: '../logs/logs'
        });
    }
}

Page(new IndexPage());
