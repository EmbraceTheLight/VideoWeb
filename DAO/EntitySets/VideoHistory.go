package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type VideoHistory struct {
	define.MyModel
	UID int64 `json:"UID" gorm:"column:user_id;type:bigint;primaryKey"`
	VID int64 `json:"VID" gorm:"column:video_id;type:bigint;primaryKey"`
}

func (v *VideoHistory) TableName() string {
	return "video_history"
}

func InsertVideoHistoryRecord(db *gorm.DB, vh *VideoHistory) error {
	return db.Model(&VideoHistory{}).Create(vh).Error
}

// DeleteVideoHistoryRecord 根据视频ID和用户ID删除一条视频历史记录
func DeleteVideoHistoryRecord(db *gorm.DB, uid int64, vid int64) error {
	return db.Model(&VideoHistory{}).Where("user_id = ? AND video_id = ?", uid, vid).Delete(&VideoHistory{}).Error
}

// DeleteAllVideoHistoryRecords 删除用户的所有视频历史记录
func DeleteAllVideoHistoryRecords(db *gorm.DB, uid int64) error {
	return db.Model(&VideoHistory{}).Where("user_id = ?", uid).Delete(&VideoHistory{}).Error
}
