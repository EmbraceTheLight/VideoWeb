package DAO

import (
	"VideoWeb/define"
)

type Message struct {
	define.MyModel

	MsgID    string `json:"msgID" gorm:"column:MsgID;type:char(36);primaryKey"` //消息ID
	UID      string `json:"UID" gorm:"column:UID;type:char(36)"`                //消息接收者ID
	BeRead   int    `json:"beRead" gorm:"column:BeRead;type:tinyint;default:1"` //是否已读标记，1表示未读，-1表示已读
	Title    string `json:"title" gorm:"column:title;type:varchar(100)"`        //消息标题
	Content  string `json:"content" gorm:"column:Content;type:text"`            //消息内容
	SenderID string `json:"SenderID" gorm:"column:SenderID;type:char(36);"`     //消息发送者ID
}

func (m *Message) TableName() string {
	return "Messages"
}
