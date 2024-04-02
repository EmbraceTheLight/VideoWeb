package DAO

import "gorm.io/gorm"

// Barrage 弹幕表信息
type Barrage struct {
	gorm.Model

	UID     string `json:"uid" gorm:"column:UID;type:char(36);index:INDEX_UID"`
	VID     string `json:"vid" gorm:"column:VID;type:char(36);index:INDEX_VID"`
	Content string `json:"content" gorm:"column:content;type:text"`
}
