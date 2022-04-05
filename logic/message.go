package logic

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

const (
	MsgTypeNormal    = iota // 普通 用户消息
	MsgTypeWelcome          // 当前用户欢迎消息
	MsgTypeUserEnter        // 用户进入
	MsgTypeUserLeave        // 用户退出
	MsgTypeError            // 错误消息
)

// 构建 message 实例
func NewMessage(user *User, content, sentTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if sentTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(sentTime))
	}
	return message

}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " 离开了聊天室",
		MsgTime: time.Now(),
	}
}

// 欢迎用户 message
func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: fmt.Sprintf("欢迎 用户 %v 加入聊天室！", user.NickName),
		MsgTime: time.Now(),
	}
}

// 错误 message
func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
