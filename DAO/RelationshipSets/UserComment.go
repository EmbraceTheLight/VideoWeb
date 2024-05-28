package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type UserLikedComments struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;index:idx_uid_cid"`
	CID int64 `json:"vid" gorm:"primaryKey;column:comment_id;type:bigint;index:idx_uid_cid"`
}

func (*UserLikedComments) TableName() string {
	return "user_liked_comments"
}

func InsertUserLikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	ulc := &UserLikedComments{
		MyModel: define.MyModel{},
		UID:     uid,
		CID:     cid,
	}
	return db.Model(&UserLikedComments{}).Create(ulc).Error
}

func DeleteUserLikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	return db.Model(&UserLikedComments{}).Where("user_id = ? AND comment_id = ?", uid, cid).Delete(&UserLikedComments{}).Error
}

type UserDislikedComments struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;index:idx_uid_cid"`
	CID int64 `json:"vid" gorm:"primaryKey;column:comment_id;type:bigint;index:idx_uid_cid"`
}

func (*UserDislikedComments) TableName() string {
	return "user_unliked_comments"
}

func InsertUserDislikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	udlc := &UserLikedComments{
		MyModel: define.MyModel{},
		UID:     uid,
		CID:     cid,
	}
	return db.Model(&UserDislikedComments{}).Create(udlc).Error
}

func DeleteUserDislikedCommentRecord(db *gorm.DB, uid, cid int64) error {
	return db.Model(&UserDislikedComments{}).Where("user_id = ? AND comment_id = ?", uid, cid).Delete(&UserDislikedComments{}).Error
}
