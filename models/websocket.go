package models

import (
	"golang.org/x/net/websocket"
)

//type WSServer struct {
//	ListenAddr string
//}

//socket统一结构
type SocketMessage struct {
	MessageType				int					`description:"消息类型，0为建立连接，1为普通聊天，2为被挤下线" json:"messageType" `
	MessageSendTime			int64				`description:"消息发送时间" json:"messageSendTime" `
	MessageSendUid			int64				`description:"心跳发送uid" json:"messageSendUid" `
	MessageExpireTime		int64				`description:"心跳有效时间" json:"messageExpireTime" `
	MessageContent			string				`description:"消息内容，jsonString" json:"messageContent" `
	MessageSign				string				`description:"消息签名" json:"messageSign" `
	MessageToken			string				`description:"用户token" json:"messageToken" `
	MessageAppNo			int					`description:"消息来源，1为得问，2为掌律，3为学习大师" json:"messageAppNo" `
}

//聊天结构体
type UserSocketMessage struct {
	FromNickName  			string 	 			`description:"fromNickName" json:"fromNickName" `
	FromUid       			int64 	 			`description:"fromUid" json:"fromUid" `
	ToNickName         		string 	 			`description:"toNickName" json:"toNickName" `
	ToUid         			int64 	 			`description:"toUid" json:"toUid" `
	GroupId           		int64	  			`description:"groupId" json:"groupId" `
	GroupType        		int    				`description:"组类型 0:一对一 1:一对多 2:系统消息（不存库，仅识别使用） 3:客服一对一 " json:"groupType"`
	From					int	   				`description:"客服显示当前咨询者进入IM的入口，1.首页 2.首页列表（张三1231）3.次首页 4.次首页列表（张三1231）5.列表页（张三1231）6.详情页（张三1231）0.其他" json:"from" xorm:"notnull default 0"`
	Param					string 				`description:"咨询者进入IM的入口相关信息(json string)" json:"param" valid:"MaxSize(300)" `
	Content           		string  			`description:"content" json:"content" `
	ActionUrl         		string  			`description:"actionUrl" json:"actionUrl" `
	Type	           		int		  			`description:"消息内容类型 0:文本 1:图片 2:语音 3:视频" json:"type" `
	ImageWidth     			string 				`description:"图片宽度,客户端根据这个显示图片宽度" json:"imageWidth"`
	ImageHeight     		string 				`description:"图片高度,客户端根据这个显示图片高度" json:"imageHeight"`
}

//socket签名key
const SOCKET_MESSAGE_SIGN_KEY string = "wenshixiong123socketmessage"

const SOCKET_UNSENT_MESSAGE = "SocketUnsentMessage"

//连接存储
type SocketConnection struct {
	Conn				*websocket.Conn			`description:"socket连接" json:"conn"`
	ExpireTime				int64				`description:"socket连接有效截止时间" json:"expireTime"`
	Token					string				`description:"用户token" json:"token"`
	AppNo					int					`description:"app号，1为得问，2为掌律，3为学习大师" json:"appNo"`
}



