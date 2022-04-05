package logic

import (
	"chatroom/global"
	"container/ring"
)

type offlineProcessor struct {
	n int

	// 保存所有用户最近的 n 条消息
	recentRing *ring.Ring

	// 保存某个用户离线消息（一样 n 条）
	userRing map[string]*ring.Ring
}

// 饿汉 单例模式
var OfflineProcessor = &offlineProcessor{
	n:          global.Conf.OfflineNum,
	recentRing: ring.New(global.Conf.OfflineNum),
	userRing:   make(map[string]*ring.Ring),
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}

	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Ats {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *User) {
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChan <- value.(*Message)
		}
	})

	if user.IsNew {
		return
	}

	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChan <- value.(*Message)
			}
		})

		delete(o.userRing, user.NickName)
	}
}
