package routes

import (
	"chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	// 开启广播消息
	go logic.Broadcasters.Start()
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)

}
