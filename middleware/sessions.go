package middleware

import (
	"sync"
	"time"
)

type WxSessionManager struct {
	mLifeTime int64               // 存活时间
	mLock     sync.RWMutex        //互斥锁头
	mSessions map[string]*Session //session指针指向内容
}

// session 内容
type Session struct {
	lastTime   time.Time                   // 登录的时间已用来回收
	wxSessions string                      // 微信的key
	mValues    map[interface{}]interface{} // session设置的具体值
}

func NewSessionManager(lifeTime int64) *WxSessionManager {
	manager := &WxSessionManager{
		mLifeTime: lifeTime,
		mSessions: make(map[string]*Session)}
	// todo 定时回收
	go manager.GC()
	return manager

}

func (manager *WxSessionManager) Set(thirdPartyKey string, key interface{}, value interface{}) {
	//加锁
	manager.mLock.Lock()
	defer manager.mLock.Unlock()
	if session, ok := manager.mSessions[thirdPartyKey]; ok {
		session.mValues[key] = value
	}

}

func (manager *WxSessionManager) Get(thirdPartyKey string, key interface{}) (interface{}, bool) {
	manager.mLock.RLock()
	defer manager.mLock.RUnlock()
	if session, ok := manager.mSessions[thirdPartyKey]; ok {
		if val, ok := session.mValues[key]; ok {
			return val, ok
		}
	}
	return nil, false
}

func (manager *WxSessionManager) GC() {
	manager.mLock.Lock()
	defer manager.mLock.Unlock()
	for sessionId, session := range manager.mSessions {
		if session.lastTime.Unix()+manager.mLifeTime > time.Now().Unix() {
			delete(manager.mSessions, sessionId)
		}
	}
	time.AfterFunc(time.Duration(manager.mLifeTime)*time.Second, func() { manager.GC() })
}
