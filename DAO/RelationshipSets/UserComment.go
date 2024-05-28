package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// UserLikedComments 用户点赞的评论记录
type UserLikedComments struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;"`
	VID int64 `json:"vid" gorm:"primaryKey;column:video_id;type:bigint;"`
	CID int64 `json:"cid" gorm:"primaryKey;column:comment_id;type:bigint;"`
}

func (*UserLikedComments) TableName() string {
	return "user_liked_comments"
}

func GetUserLikedCommentRecordByUidVid(db *gorm.DB, uid, vid int64) (ulc []*UserLikedComments, err error) {
	err = db.Model(&UserLikedComments{}).Where("user_id = ? AND video_id = ?", uid, vid).Find(&ulc).Error
	return
}

func InsertUserLikedCommentRecord(db *gorm.DB, uid, vid, cid int64) error {
	ulc := &UserLikedComments{
		MyModel: define.MyModel{},
		UID:     uid,
		VID:     vid,
		CID:     cid,
	}
	return db.Model(&UserLikedComments{}).Create(ulc).Error
}

func DeleteUserLikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	return db.Model(&UserLikedComments{}).Where("user_id = ? AND comment_id = ?", uid, cid).Delete(&UserLikedComments{}).Error
}

// UserDislikedComments 用户点踩的评论记录
type UserDislikedComments struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;"`
	VID int64 `json:"vid" gorm:"primaryKey;column:video_id;type:bigint;"`
	CID int64 `json:"cid" gorm:"primaryKey;column:comment_id;type:bigint;"`
}

func (*UserDislikedComments) TableName() string {
	return "user_unliked_comments"
}

func GetUserDislikedCommentRecordByUidVid(db *gorm.DB, uid, vid int64) (udlc []*UserDislikedComments, err error) {
	err = db.Model(&UserDislikedComments{}).Where("user_id = ? AND video_id = ?", uid, vid).Find(&udlc).Error
	return
}

func InsertUserDislikedCommentRecord(db *gorm.DB, uid, vid, cid int64) error {
	udlc := &UserLikedComments{
		MyModel: define.MyModel{},
		UID:     uid,
		VID:     vid,
		CID:     cid,
	}
	return db.Model(&UserDislikedComments{}).Create(udlc).Error
}

func DeleteUserDislikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	return db.Model(&UserDislikedComments{}).Where("user_id = ? AND comment_id = ?", uid, cid).Delete(&UserDislikedComments{}).Error
}
