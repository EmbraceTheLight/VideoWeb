package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type VideoHistory struct {
	define.MyModel
	UID string `json:"UID" gorm:"column:UID;type:char(36);primaryKey"`
	VID string `json:"VID" gorm:"column:VID;type:char(36);primaryKey"`
}

func (v *VideoHistory) TableName() string {
	return "VideoHistory"
}

func (v *VideoHistory) Create(DB *gorm.DB) error {
	result := DB.Create(&v)
	return result.Error
}
