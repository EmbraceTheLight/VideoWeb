package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// UserFollowed 用户粉丝列表
type UserFollowed struct {
	define.MyModel

	UID string `json:"UID" gorm:"primaryKey;column:UID;type:char(36);index:index_UserID"`     //用户ID
	FID string `json:"FID" gorm:"primaryKey;column:FID;type:char(36);index:index_FollowedID"` //粉丝的ID
	//Fans *EntitySets.User `gorm:"foreignKey:UserID;references:FID"`                           //粉丝详细信息
}

func (ufd *UserFollowed) Create(DB *gorm.DB) error {
	result := DB.Create(&ufd)
	return result.Error
}
func (ufd *UserFollowed) Delete(DB *gorm.DB) error {
	result := DB.Where("UID=? AND FID=?", ufd.UID, ufd.FID).Delete(&ufd)
	return result.Error
}

func (ufd *UserFollowed) TableName() string {
	return "UserFollowed"
}
