package DAO

import (
	"VideoWeb/define"
)

type Comments struct {
	define.MyModel

	CommentID string `json:"CommentID" gorm:"column:commentID;type:char(36);primaryKey"` //评论唯一标识
	UID       string `json:"UID" gorm:"column:UID;type:char(36)"`                        //评论者ID
	To        string `json:"To" gorm:"column:To;type:char(36)"`                          //若是回复评论，则To为被回复评论ID，否则为Null
	VID       string `json:"VID" gorm:"column:VID;type:char(36)"`                        //视频ID
	Content   string `json:"Content" gorm:"column:content;type:text"`                    //评论内容
	Likes     uint32 `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`      //评论被点赞数
	UnLikes   uint32 `json:"UnLikes" gorm:"column:unlikes;type:int unsigned;default:0"`  //评论被点踩数
	IPAddress string `json:"IPAddress" gorm:"column:IPAddress;type:varchar(15)"`         //评论者IP归属地
}
