package manager

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

/**
	定义一个全局的session管理器
 */

type Manager struct {
	cookieName string
	lock sync.Mutex
	provider Provider         //驱动
	maxLifeTime int64         // 最大存在时间
}

func NewManager(provideName,cookieName string,maxLifeTime int64) (*Manager,error)  {
	provider, ok := provides[provideName]
	if !ok {
		return nil,fmt.Errorf("session: unkown provide %q (forgotten import?)",provideName)
	}
	return &Manager{
		cookieName:  cookieName,
		//lock:        sync.Mutex{},
		provider:    provider,
		maxLifeTime: maxLifeTime,
	},nil
}

var globalSession *session.Manger

func init()  {
	globalSession,_ = NewManager("memory","gosessionid",3600)
}

type Provider interface{
	SessionInit(sid string)(Session ,error)
	SessionRead(sid string)(Session ,error)
	SessionDestory(sid string) error
	SessionGC(maxLifeTime  int64)
}


type Session  interface{
	Set(key, value interface{}) error //设置session值
	Get(key interface{}) interface{}  //获取session值
	Delete(key interface{}) error     //删除session值
	SessionID()  string               //返回正确的sessionID
}


var provides = make(map[string]Provider)

func Register(name string,provider Provider)  {
	if provider == nil {
		panic("Session: Register provider is nil")
	}
	if _,dup := provides[name];dup{
		panic("Seesion:Register called twice for provider" + name)
	}
	provides[name] = provider
}

func (manager *Manager) sessionID() string  {
	b := make([]byte,32)
	if _,err := rand.Read(b);err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter,r *http.Request)  {
	manager.lock.Lock()
	defer manager.lock.Lock()
	cookie,err := r.Cookie(manager.cookieName)
	if err != nil  || cookie.Value == ""{
		sid := manager.sessionID()
		session,_ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:      manager.cookieName,
			Value:      url.QueryEscape(sid),
			Path:       "/",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     int(manager.maxLifeTime),
			Secure:     false,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		http.SetCookie(w,&cookie)
	}else{
		sid , _ := url.QueryUnescape(cookie.Value)
	}
}