package DAO

import (
	"gorm.io/gorm"
)

// Barrage 弹幕表信息
type Barrage struct {
	gorm.Model

	UID     string `json:"uid" gorm:"column:UID;type:char(36);index:INDEX_UID"`
	VID     string `json:"vid" gorm:"column:VID;type:char(36);index:INDEX_VID"`
	Content string `json:"content" gorm:"column:content;type:text"`
	Minute  int    `json:"minute" gorm:"column:minute;type:int"`
	Second  int    `json:"second" gorm:"column:second;type:int"`
}

func (b *Barrage) Create(DB *gorm.DB) error {
	err := DB.Create(&b).Error
	return err
}

func (b *Barrage) TableName() string {
	return "Barrages"
}
