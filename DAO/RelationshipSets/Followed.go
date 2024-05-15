package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// UserFollowed 用户粉丝列表
type UserFollowed struct {
	define.MyModel

	UID int64 `json:"UID" gorm:"primaryKey;column:user_id;type:bigint;index:index_UserID"`         //用户ID
	FID int64 `json:"FID" gorm:"primaryKey;column:followed_id;type:bigint;index:index_FollowedID"` //粉丝的ID
	//Fans *EntitySets.User `gorm:"foreignKey:UserID;references:FID"`                           //粉丝详细信息
}

func InsertFollowedRecord(db *gorm.DB, ufd *UserFollowed) error {
	return db.Model(&UserFollowed{}).Create(&ufd).Error
}

func DeleteFollowedRecord(DB *gorm.DB, ufd *UserFollowed) error {
	return DB.Model(&UserFollowed{}).Where("user_id=? AND followed_id=?", ufd.UID, ufd.FID).Delete(&ufd).Error
}

func DeleteFollowedRecordsByUserID(db *gorm.DB, UserID int64) error {
	return db.Model(UserFollowed{}).Where("user_id=?", UserID).Delete(&UserFollowed{}).Error
}

func (ufd *UserFollowed) TableName() string {
	return "user_followed"
}

// GetFollowedByFollowedID 通过用户ID来获取该用户的粉丝列表
func GetFollowedByFollowedID(DB *gorm.DB, id string) ([]*UserFollowed, error) {
	var followed []*UserFollowed
	err := DB.Model(&UserFollowed{}).Where("user_id = ?", id).Find(&followed).Error
	if err != nil {
		return nil, err
	}
	return followed, nil
}
