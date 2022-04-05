package logic

import (
	"time"

	"nhooyr.io/websocket"
)

type User struct {
	UID         int             `json:"uid"`      // 用户 ID
	NickName    string          `json:"nickname"` // 用户昵称
	EnterAt     time.Time       `json:"enter_at"` // 进入时间
	Addr        string          `json:"addr"`
	Token       string          `json:"token"`
	MessageChan chan *Message   `json:"-"`
	Conn        *websocket.Conn // 用户连接
	IsNew       bool            // 是否是新用户
}

// 给用户发送的消息
type Message struct {
	// 哪个用户发送的消息
	User *User `json:"user"`
	// 消息类型
	Type int `json:"type"`
	// 内容
	Content string `json:"content"`
	// 消息创建时间
	MsgTime        time.Time `json:"msg_time"`
	ClientSendTime time.Time `json:"client_send_time"`

	// 消息 @ 了谁（一次可以艾特多个人）
	Ats []string `json:"ats"`
}

type Broadcaster struct {
	// 所有聊天室用户
	Users map[string]*User

	// 进入
	EnteringChannel chan *User
	// 离开
	LeavingChannel chan *User
	// 广播消息
	MessageChannel chan *Message

	// 判断该昵称用户是否可进入聊天室（重复与否）：true 能，false 不能
	CheckUserChannel      chan string
	CheckUserCanInChannel chan bool

	// 获取用户列表
	RequestUsersChannel chan struct{}
	UsersChannel        chan []*User
}
