package ws

import "github/lhh-gh/easy-chat/pkg/constants"

type (
	Msg struct {
		constants.MsgType `mapstructure:"msgType"`
		Content           string `mapstructure:"content"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}
)
