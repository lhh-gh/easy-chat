package logic

import (
	"context"
	"github/lhh-gh/easy-chat/apps/im/immodels"
	"github/lhh-gh/easy-chat/apps/im/ws/internal/svc"
	"github/lhh-gh/easy-chat/apps/im/ws/websocket"
	"github/lhh-gh/easy-chat/apps/im/ws/ws"
	"github/lhh-gh/easy-chat/pkg/wuid"

	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

/*
data：私聊数据
userId：发送人
*/
func (c *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MsgType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}

	err := c.svc.ChatLogModel.Insert(c.ctx, &chatLog)

	return err
}
