package DAO

import (
	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	UID          int64  `json:"UID" gorm:"column:UID;type:bigint;primaryKey"`
	SearchString string `json:"searchString" gorm:"column:searchString;type:varchar(100);"` //用户搜索字符串
}

func (s *SearchHistory) TableName() string {
	return "search_history"
}

func InsertSearchRecord(db *gorm.DB, sh *SearchHistory) error {
	return db.Model(&SearchHistory{}).Create(sh).Error
}
