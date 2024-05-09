package DAO

import (
	"gorm.io/gorm"
)

// UserFollows 用户关注列表
type UserFollows struct {
	gorm.Model
	GroupName string `json:"groupName" gorm:"column:groupName;type:varchar(15);"`     //关注分组的名称，不超过15个字
	UID       int64  `json:"UID" gorm:"column:UID;type:bigint;index:index_UID"`       //用户ID
	FID       int64  `json:"FID" gorm:"column:FID;type:bigint;index:index_FollowsID"` //关注的人的ID
}

func (ufs *UserFollows) TableName() string {
	return "UserFollows"
}

func InsertFollowsRecord(db *gorm.DB, ufs *UserFollows) error {
	result := db.Model(UserFollows{}).Create(&ufs)
	return result.Error
}

func DeleteFollowsRecord(db *gorm.DB, ufs *UserFollows) error {
	result := db.Model(UserFollows{}).Where("UID=? AND FID=?", ufs.UID, ufs.FID).Delete(&ufs)
	return result.Error
}

// GetFollowsByUserID 通过用户ID来获取该用户的关注列表
func GetFollowsByUserID(db *gorm.DB, id string) ([]*UserFollows, error) {
	var follows = make([]*UserFollows, 0)
	err := db.Model(UserFollows{}).Where("UID = ?", id).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}
