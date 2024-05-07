package DAO

import (
	"gorm.io/gorm"
)

// UserFollows 用户关注列表
type UserFollows struct {
	gorm.Model
	GroupName string `json:"groupName" gorm:"column:groupName;type:varchar(15);"`       //关注分组的名称，不超过15个字
	UID       string `json:"UID" gorm:"column:UID;type:char(36);index:index_UID"`       //用户ID
	FID       string `json:"FID" gorm:"column:FID;type:char(36);index:index_FollowsID"` //关注的人的ID
}

func (ufs *UserFollows) TableName() string {
	return "UserFollows"
}

func InsertFollowsRecord(DB *gorm.DB, ufs *UserFollows) error {
	result := DB.Create(&ufs)
	return result.Error
}

func DeleteFollowsRecord(DB *gorm.DB, ufs *UserFollows) error {
	result := DB.Where("UID=? AND FID=?", ufs.UID, ufs.FID).Delete(&ufs)
	return result.Error
}

// GetFollowsByUserID 通过用户ID来获取该用户的关注列表
func GetFollowsByUserID(DB *gorm.DB, id string) ([]*UserFollows, error) {
	var follows = make([]*UserFollows, 0)
	err := DB.Where("UID = ?", id).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}
