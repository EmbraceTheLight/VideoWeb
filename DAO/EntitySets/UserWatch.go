package DAO

import (
	"VideoWeb/define"
	"fmt"
	"gorm.io/gorm"
)

type UserWatch struct {
	define.MyModel
	UID int64 `json:"UID" gorm:"column:user_id;type:bigint;primaryKey"`
	VID int64 `json:"VID" gorm:"column:video_id;type:bigint;primaryKey"`
}

func (v *UserWatch) GetScore() float64 {
	return float64(v.UpdatedAt.UnixMilli())
}

func (v *UserWatch) GetValue() any {
	return v.VID
}

func (v *UserWatch) TableName() string {
	return "user_watch"
}

// InsertVideoHistoryRecord 插入或更新VideoHistory记录
func InsertVideoHistoryRecord(db *gorm.DB, vh *UserWatch) error {
	err := db.Where("user_id = ? AND video_id = ?", vh.UID, vh.VID).Omit("created_at").Save(vh).Error
	if err != nil {
		return fmt.Errorf("UserWatch.InsertVideoHistoryRecord: %w", err)
	}
	return nil
}

// DeleteVideoHistoryRecord 根据视频ID和用户ID删除一条视频历史记录
func DeleteVideoHistoryRecord(db *gorm.DB, uid int64, vid int64) error {
	return db.Model(&UserWatch{}).Where("user_id = ? AND video_id = ?", uid, vid).Delete(&UserWatch{}).Error
}

// DeleteAllVideoHistoryRecords 删除用户的所有视频历史记录
func DeleteAllVideoHistoryRecords(db *gorm.DB, uid int64) error {
	return db.Model(&UserWatch{}).Where("user_id = ?", uid).Delete(&UserWatch{}).Error
}
