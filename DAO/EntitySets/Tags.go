package DAO

import "gorm.io/gorm"

type Tags struct {
	Tag string `json:"Tag" gorm:"column:tag;type:varchar(15);primaryKey"`
	VID int64  `json:"VID" gorm:"column:video_id;type:bigint;primaryKey"`
}

func (t *Tags) TableName() string {
	return "tags"
}

func InsertTags(db *gorm.DB, tags []*Tags) error {
	return db.Model(&Tags{}).Create(tags).Error
}

// DeleteTagRecords 删除视频标签记录
func DeleteTagRecords(db *gorm.DB, videoID int64) error {
	err := db.Model(&Tags{}).Delete(&Tags{}, "video_id=?", videoID).Error
	return err
}
