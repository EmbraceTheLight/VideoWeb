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
	return "video_history"
}

func InsertVideoHistoryRecord(db *gorm.DB, vh *VideoHistory) error {
	return db.Model(&VideoHistory{}).Create(vh).Error
}
