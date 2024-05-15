package DAO

import (
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"gorm.io/gorm"
)

type FollowList struct {
	ListID   int64                           `json:"id" gorm:"column:list_id;primary_key;type:bigint;"`          //关注列表ID
	UID      int64                           `json:"userID" gorm:"column:user_id;primary_key;type:bigint;"`      //关注列表拥有者ID
	ListName string                          `json:"listName" gorm:"column:list_name;type:varchar(20);not null"` //关注列表名称
	Users    []*RelationshipSets.UserFollows `gorm:"foreignKey:FollowListID;references:ListID"`
}

func (*FollowList) TableName() string {
	return "follow_list"
}

func InsertFollowList(db *gorm.DB, fl *FollowList) error {
	return db.Model(&FollowList{}).Create(fl).Error
}

// DeleteFollowListByUserID 根据用户ID删除用户的关注列表信息
func DeleteFollowListByUserID(db *gorm.DB, uid int64) error {
	return db.Model(&FollowList{}).Where("user_id =?", uid).Delete(&FollowList{}).Error
}
