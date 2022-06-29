package user

import (
	"skynet-service/app/logs"
	"sync"
	"time"
)

type UserSession struct {
	sync.Mutex
	allUsers 		map[string]*UserInfo
}

type UserInfo struct {
	loginTime   int64
	// 保存的一些信息
}

var GAllUsers = newUser ()

func init () {
	logs.Log.Debug("初始化用户模块...")
}

func newUser () *UserSession {
	return &UserSession {
		allUsers : map[string]*UserInfo{},
	}
}

// 新登录的用户
func (self *UserSession) NewUser (userName string) {
	self.Lock()
	defer self.Unlock()

	userInfo := UserInfo {time.Now().Unix()}

	self.allUsers[userName] = &userInfo
}

// 验证用户是否合法,先判断是否已登录，如果没有登录则查找此用户是否注册
func (self *UserSession) ValidUser (userName string) bool {
	return true
}