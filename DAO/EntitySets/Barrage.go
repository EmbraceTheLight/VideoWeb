package DAO

import (
	"gorm.io/gorm"
)

// Barrage 弹幕表信息
type Barrage struct {
	gorm.Model

	Hour    uint8  `json:"hour" gorm:"column:hour;type:char(2)"`
	Minute  uint8  `json:"minute" gorm:"column:minute;type:char(2)"`
	Second  uint8  `json:"second" gorm:"column:second;type:char(2)"`
	UID     string `json:"uid" gorm:"column:UID;type:char(36);index:INDEX_UID"`
	VID     string `json:"vid" gorm:"column:VID;type:char(36);index:INDEX_VID"`
	Content string `json:"content" gorm:"column:content;type:text"`
	Color   string `json:"color" gorm:"column:color;type:char(8)"`
}

func InsertBarrageRecord(DB *gorm.DB, b *Barrage) error {
	err := DB.Create(&b).Error
	return err
}

func (b *Barrage) TableName() string {
	return "Barrages"
}
