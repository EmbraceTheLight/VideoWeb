package DAO

import (
	"VideoWeb/define"
)

type VideoHistory struct {
	define.MyModel
	UID string `json:"UID" gorm:"column:UID;type:char(36);primaryKey"`
	VID string `json:"VID" gorm:"column:VID;type:char(36);primaryKey"`
}

func (v *VideoHistory) TableName() string {
	return "VideoHistory"
}
