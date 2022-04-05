package logic

import (
	"chatroom/global"
	"expvar"
	"fmt"
)

// 饿汉式单例模式（全局广播器）
var Broadcasters = &Broadcaster{
	Users:                 make(map[string]*User),
	EnteringChannel:       make(chan *User),
	LeavingChannel:        make(chan *User),
	MessageChannel:        make(chan *Message, global.Conf.MessageQueue),
	CheckUserChannel:      make(chan string),
	CheckUserCanInChannel: make(chan bool),
	RequestUsersChannel:   make(chan struct{}),
	UsersChannel:          make(chan []*User),
}

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() interface{} {
	fmt.Println("===len=:", len(Broadcasters.MessageChannel))
	return len(Broadcasters.MessageChannel)
}

func (b *Broadcaster) Start() {
	for {
		select {
		case user := <-b.EnteringChannel:
			// 用户进入
			b.Users[user.NickName] = user

			OfflineProcessor.Send(user)
		case user := <-b.LeavingChannel:
			// 用户离开
			delete(b.Users, user.NickName)
			// 避免 goroutine 泄露
			user.CloseMessageChannel()
		case msg := <-b.MessageChannel:
			// 给所有在线用户发送消息
			for _, user := range b.Users {
				// 排除发送消息的 用户
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChan <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.CheckUserChannel:
			// 检查 用户是否已经进入过
			if _, ok := b.Users[nickname]; ok {
				b.CheckUserCanInChannel <- false
			} else {
				b.CheckUserCanInChannel <- true
			}
		case <-b.RequestUsersChannel:
			userList := make([]*User, 0, len(b.Users))
			for _, user := range b.Users {
				userList = append(userList, user)
			}

			b.UsersChannel <- userList
		}
	}
}

// 用户进入
func (b *Broadcaster) UserEntering(user *User) {
	b.EnteringChannel <- user
}

// 用户离开
func (b *Broadcaster) UserLeaving(user *User) {
	b.LeavingChannel <- user
}

func (b *Broadcaster) Broadcast(msg *Message) {
	b.MessageChannel <- msg
}

func (b *Broadcaster) IsEnterRoom(nickname string) bool {
	b.CheckUserChannel <- nickname
	return <-b.CheckUserCanInChannel
}

func (b *Broadcaster) GetUserList() []*User {
	b.RequestUsersChannel <- struct{}{}
	return <-b.UsersChannel
}
