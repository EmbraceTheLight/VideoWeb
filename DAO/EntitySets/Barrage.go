package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// Barrage 弹幕表信息
type Barrage struct {
	Hour   uint8 `json:"hour" gorm:"column:hour;type:char(2)"`
	Minute uint8 `json:"minute" gorm:"column:minute;type:char(2)"`
	Second uint8 `json:"second" gorm:"column:second;type:char(2)"`

	BID int64 `json:"bid" gorm:"column:barrage_id;type:bigint;primary_key"`
	UID int64 `json:"uid" gorm:"column:user_id;type:bigint;index:INDEX_UID"`
	VID int64 `json:"vid" gorm:"column:video_id;type:bigint;index:INDEX_VID"`
	define.MyModel

	Content string `json:"content" gorm:"column:content;type:text"`
	Color   string `json:"color" gorm:"column:color;type:char(8)"`
}

func (b *Barrage) TableName() string {
	return "barrages"
}

// InsertBarrageRecord 插入一条弹幕记录
func InsertBarrageRecord(DB *gorm.DB, b *Barrage) error {
	return DB.Model(&Barrage{}).Create(&b).Error
}

// DeleteBarrageRecordsByVideoID 删除对应视频的所有弹幕记录
func DeleteBarrageRecordsByVideoID(db *gorm.DB, VID int64) error {
	return db.Delete(&Barrage{}, "video_id = ?", VID).Error
}
