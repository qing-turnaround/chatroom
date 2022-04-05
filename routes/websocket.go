package routes

import (
	"chatroom/logic"
	"log"
	"net/http"

	"go.uber.org/zap"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	// Accept 从客户端接受 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		zap.L().Error("websocket Accept error: ", zap.Error(err))
		return
	}

	// 1. 新用户进来，构建该用户的实例
	token := r.FormValue("token")       // 接收 token
	nickname := r.FormValue("nickname") // 接收昵称信息
	// 检查昵称的合法
	if l := len(nickname); l < 2 || l > 20 {
		zap.L().Error("nickname must be at least 2 characters and more than 20 characters：",
			zap.Error(err))
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度为：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	// 检查昵称是否已经存在
	if !logic.Broadcasters.IsEnterRoom(nickname) {
		zap.L().Error("昵称已经存在："+nickname, zap.Error(nil))
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("昵称已经存在！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists")
		return
	}

	// 实例
	userToken := logic.NewUser(conn, token, nickname, r.RemoteAddr)

	// 2. 开启给新用户发送消息的 goroutine
	go userToken.SendMessage(r.Context())

	// 3. 新用户进入，给新用户发送欢迎信息
	userToken.MessageChan <- logic.NewWelcomeMessage(userToken)
	// 避免 token 泄露
	tmpUser := *userToken
	user := &tmpUser
	user.Token = ""
	// 告知其他用户，新用户进入了聊天室
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcasters.Broadcast(msg)

	// 4. 将 该用户加入广播器的用户列表中
	logic.Broadcasters.UserEntering(user)
	log.Println("user:", nickname, "join chatroom!")

	// 5. 接收用户信息
	err = user.ReceiveMessage(r.Context())

	// 6. 用户离开
	logic.Broadcasters.UserLeaving(user)
	// 广播用户离开
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcasters.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 关闭 Conn
	if err == nil {
		// 正常关闭
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		// 关闭并记录错误
		zap.L().Error("Read from client error", zap.Error(err))
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}
