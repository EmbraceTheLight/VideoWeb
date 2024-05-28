package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type CommentSummary struct {
	Comments
	Replies []*CommentSummary `json:"Replies" gorm:"-"` //回复这条评论的评论，可能包含多重嵌套
}

type Comments struct {
	define.MyModel

	CommentID int64  `json:"CommentID" gorm:"column:comment_id;type:bigint;primaryKey"`   //评论唯一标识
	UID       int64  `json:"UID" gorm:"column:user_id;type:bigint"`                       //评论者ID
	To        int64  `json:"To" gorm:"column:to;type:bigint;default:-1;index:idx_to_vid"` //若是回复评论，则To为被回复评论ID，否则为Null
	VID       int64  `json:"VID" gorm:"column:video_id;type:bigint;index:idx_to_vid"`     //视频ID
	Content   string `json:"Content" gorm:"column:content;type:text"`                     //评论内容
	Likes     uint32 `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`       //评论被点赞数
	UnLikes   uint32 `json:"UnLikes" gorm:"column:unlikes;type:int unsigned;default:0"`   //评论被点踩数
	IPAddress string `json:"IPAddress" gorm:"column:ip_address;type:varchar(15)"`         //评论者IP归属地
}

func (c *Comments) TableName() string { return "comments" }

func InsertCommentRecord(db *gorm.DB, c *Comments) error {
	return db.Model(&Comments{}).Create(c).Error
}

func DeleteCommentRecordsByVideoID(db *gorm.DB, videoID int64) error {
	return db.Model(&Comments{}).Delete(&Comments{}, "video_id = ?", videoID).Error
}
