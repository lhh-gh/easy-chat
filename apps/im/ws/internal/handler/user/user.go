package user

import (
	"github/lhh-gh/easy-chat/apps/im/ws/internal/svc"
	"github/lhh-gh/easy-chat/apps/im/ws/websocket"
)

func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 获取所有用户
		uids := srv.GetUsers()
		// 获取连接的用户
		u := srv.GetUsers(conn)
		// 发送消息
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
