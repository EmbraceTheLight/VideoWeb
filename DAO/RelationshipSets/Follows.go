package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// UserFollows 用户关注列表
type UserFollows struct {
	define.MyModel

	FollowListID int64 `json:"FollowListID" gorm:"primaryKey;column:followlist_id;type:bigint;index:index_UID"` //关注列表ID
	UID          int64 `json:"UID" gorm:"primaryKey;column:user_id;type:bigint;index:index_UID"`                //用户ID
	FID          int64 `json:"FID" gorm:"primaryKey;column:follow_user_id;type:bigint;index:index_FollowsID"`   //关注的人的ID
}

func (ufs *UserFollows) TableName() string {
	return "user_follows"
}

// GetUserFollows 获取用户关注的用户信息,返回的是一个UserFollows的数组,FollowListID为0时,表示获取所有关注的用户信息
func GetUserFollows(db *gorm.DB, FollowListID int64, UID int64) (ufs []*UserFollows, err error) {
	if FollowListID == 0 {
		err = db.Model(UserFollows{}).Where("user_id=?", UID).Find(&ufs).Error
	} else {
		err = db.Model(UserFollows{}).Where("followlist_id=? AND user_id=?", FollowListID, UID).Find(&ufs).Error
	}
	return
}

// InsertFollowsRecord 插入一条关注信息
func InsertFollowsRecord(db *gorm.DB, ufs *UserFollows) error {
	return db.Model(UserFollows{}).Create(&ufs).Error
}

// DeleteFollowsRecord 删除关注记录
func DeleteFollowsRecord(db *gorm.DB, ufs *UserFollows) error {
	return db.Model(UserFollows{}).Where("followlist_id=? AND follow_user_id=?", ufs.FollowListID, ufs.FID).Delete(&ufs).Error
}

// DeleteFollowsRecordsByUserID 删除所有关注用户的记录
func DeleteFollowsRecordsByUserID(db *gorm.DB, UserID int64) error {
	return db.Model(&UserFollows{}).Where("user_id=?", UserID).Delete(&UserFollows{}).Error
}
