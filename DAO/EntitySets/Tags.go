package DAO

import "gorm.io/gorm"

type Tags struct {
	Tag string `json:"Tag" gorm:"column:tag;type:varchar(15);primaryKey"`
	VID int64  `json:"VID" gorm:"column:VID;type:bigint;primaryKey"`
}

func (t *Tags) TableName() string {
	return "Tags"
}

// DeleteTagRecords 删除视频标签记录
func DeleteTagRecords(db *gorm.DB, VID int64) error {
	err := db.Delete(&Tags{}, "VID=?", VID).Error
	return err
}
