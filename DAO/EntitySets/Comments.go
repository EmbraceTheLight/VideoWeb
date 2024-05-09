package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type Comments struct {
	define.MyModel

	CommentID int64  `json:"CommentID" gorm:"column:commentID;type:bigint;primaryKey"`  //评论唯一标识
	UID       int64  `json:"UID" gorm:"column:UID;type:bigint"`                         //评论者ID
	To        int64  `json:"To" gorm:"column:To;type:bigint"`                           //若是回复评论，则To为被回复评论ID，否则为Null
	VID       int64  `json:"VID" gorm:"column:VID;type:bigint"`                         //视频ID
	Content   string `json:"Content" gorm:"column:content;type:text"`                   //评论内容
	Likes     uint32 `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`     //评论被点赞数
	UnLikes   uint32 `json:"UnLikes" gorm:"column:unlikes;type:int unsigned;default:0"` //评论被点踩数
	IPAddress string `json:"IPAddress" gorm:"column:IPAddress;type:varchar(15)"`        //评论者IP归属地
}

func (c *Comments) TableName() string { return "Comments" }

// DeleteCommentRecordsByVideoID 删除对应视频下所有评论记录
func DeleteCommentRecordsByVideoID(db *gorm.DB, videoID int64) error {
	err := db.Delete(&Comments{}, "VID = ?", videoID).Error
	return err
}
