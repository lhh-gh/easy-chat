package constants

type MsgType int

const (
	TextMsgType MsgType = iota // 文本消息
)

type ChatType int

const (
	GroupChatType  ChatType = iota + 1 // 群聊
	SingleChatType                     // 私聊
)
