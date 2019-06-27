package models

//存入MongoDB的聊天记录
type MongoDBMessage struct {
	MId					int64			`description:"消息id" json:"mId" `
	GroupId				int64			`description:"聊天组id" json:"groupId" `
	SenderUid			int64			`description:"发送人id" json:"senderUid" `
	Type				int				`description:"聊天类型，1为文本，2为语音" json:"type" `
	Content				string			`description:"聊天内容" json:"content" `
	Created           	int64  			`description:"消息时间" json:"created"`
}
