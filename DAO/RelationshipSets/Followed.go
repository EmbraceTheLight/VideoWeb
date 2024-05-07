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

func InsertFollowedRecord(DB *gorm.DB, ufd *UserFollowed) error {
	result := DB.Create(&ufd)
	return result.Error
}

func DeleteFollowedRecord(DB *gorm.DB, ufd *UserFollowed) error {
	result := DB.Where("UID=? AND FID=?", ufd.UID, ufd.FID).Delete(&ufd)
	return result.Error
}

func (ufd *UserFollowed) TableName() string {
	return "UserFollowed"
}

// GetFollowedByFollowedID 通过用户ID来获取该用户的粉丝列表
func GetFollowedByFollowedID(DB *gorm.DB, id string) ([]*UserFollowed, error) {
	var followed []*UserFollowed
	err := DB.Where("UID = ?", id).Find(&followed).Error
	if err != nil {
		return nil, err
	}
	return followed, nil
}
