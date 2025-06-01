package msgTransfer

import (
	"context"
	"fmt"

	"github/lhh-gh/easy-chat/apps/task/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

// 实现消费消息的接口
func (m *MsgChatTransfer) Consume(ctx context.Context, key string, value string) error {
	fmt.Println("key：", key, " value：", value)
	return nil
}

//// 必须包含 ctx 参数
//func (m *MsgChatTransfer) Consume(ctx context.Context, key string, value string) error {
//	// 示例消费逻辑
//	if err := m.svc.SomeService.HandleMessage(ctx, value); err != nil {
//		return err
//	}
//	return nil
//}
