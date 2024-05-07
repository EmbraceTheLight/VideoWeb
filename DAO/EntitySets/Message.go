package DAO

import (
	"VideoWeb/define"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

//type Message struct {
//	define.MyModel
//
//	MsgID    string `json:"msgID" gorm:"column:MsgID;type:char(36);primaryKey"` //消息ID
//	UID      string `json:"UID" gorm:"column:UID;type:char(36)"`                //消息接收者ID
//	BeRead   int    `json:"beRead" gorm:"column:BeRead;type:tinyint;default:1"` //是否已读标记，1表示未读，-1表示已读
//	Title    string `json:"title" gorm:"column:title;type:varchar(100)"`        //消息标题
//	Content  string `json:"content" gorm:"column:Content;type:text"`            //消息内容
//	SenderID string `json:"SenderID" gorm:"column:SenderID;type:char(36);"`     //消息发送者ID
//}
//
//func (m *Message) TableName() string {
//	return "Messages"
//}
//
//// GetMsgBoxByID 通过用户ID来获取该用户的未读消息列表
//func GetMsgBoxByID(DB *gorm.DB, id string) ([]*Message, error) {
//	var message []*Message
//	err := DB.Where("UID = ? AND BeRead = ?", id, 1).Find(&message).Error
//	if err != nil {
//		return message, nil
//	}
//	return message, nil
//}

type Message struct {
	MessageID      string         `bson:"_id"`        //消息ID
	RoomID         string         `bson:"RoomID"`     //房间ID
	UserID         string         `bson:"UserID"`     //发送该消息的用户ID
	MessageContent define.Message `bson:"Data"`       //消息内容
	CreatedAt      int64          `bson:"created_at"` //创建时间
}

func (m *Message) CollectionName() string {
	return "Message"
}

func InsertMessage(mongodb *mongo.Database, message *Message) error {
	_, err := mongodb.Collection(message.CollectionName()).InsertOne(context.TODO(), message)
	return err
}
