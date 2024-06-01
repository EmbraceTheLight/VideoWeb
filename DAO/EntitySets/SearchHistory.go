package DAO

import (
	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	UID          int64  `json:"UID" gorm:"column:user_id;type:bigint;primaryKey"`
	SearchString string `json:"searchString" gorm:"column:search_string;type:varchar(100);"` //用户搜索字符串
}

func (s *SearchHistory) TableName() string {
	return "search_history"
}

func InsertSearchRecord(db *gorm.DB, UID int64, searchString string) error {
	sh := &SearchHistory{
		UID:          UID,
		SearchString: searchString,
	}
	return db.Model(&SearchHistory{}).Create(sh).Error
}

// DeleteOneSearchRecord 根据记录ID删除单条搜索记录
func DeleteOneSearchRecord(db *gorm.DB, ID uint) error {
	return db.Model(&SearchHistory{}).Where("id = ?", ID).Delete(&SearchHistory{}).Error
}

// DeleteAllSearchRecord 根据用户ID删除用户的所有搜索记录
func DeleteAllSearchRecord(db *gorm.DB, uid int64) error {
	return db.Model(&SearchHistory{}).Where("user_id = ?", uid).Delete(&SearchHistory{}).Error
}
