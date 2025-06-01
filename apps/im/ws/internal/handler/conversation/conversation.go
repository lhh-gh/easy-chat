package conversation

import (
	"context"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/logic"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/svc"
	"github/lhh-gh/easy-chat/apps/im/ws/websocket"
	"github/lhh-gh/easy-chat/apps/im/ws/ws"
	"github/lhh-gh/easy-chat/pkg/constants"

	"time"

	"github.com/mitchellh/mapstructure"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 私聊
		var data ws.Chat
		// 解析客户端传递过来的数据
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		switch data.ChatType {
		case constants.SingleChatType:
			// mongo中记录消息数据
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}
			// 发送消息
			srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				SendTime:       time.Now().UnixMilli(),
				Msg:            data.Msg,
			}), data.RecvId)
		}
	}
}
