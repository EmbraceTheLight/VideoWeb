package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type VideoHistory struct {
	define.MyModel
	UID int64 `json:"UID" gorm:"column:UID;type:bigint;primaryKey"`
	VID int64 `json:"VID" gorm:"column:VID;type:bigint;primaryKey"`
}

func (v *VideoHistory) TableName() string {
	return "VideoHistory"
}

func (v *VideoHistory) Create(DB *gorm.DB) error {
	result := DB.Create(&v)
	return result.Error
}
