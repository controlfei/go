package memory

import (
	"container/list"
	"time"
	"github.com/astaxie/beego/session"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
	sid string        								//session id 唯一标识
	timeAccessed time.Time							//最后的访问时间
	value map[interface{}]interface{}				//session里面存储的值
}

func (st *SessionStore) Set(key,value interface{}) error  {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

