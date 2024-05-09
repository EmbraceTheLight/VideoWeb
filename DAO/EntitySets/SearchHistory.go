package DAO

import (
	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	UID          int64  `json:"UID" gorm:"column:UID;type:int64;primaryKey"`
	SearchString string `json:"searchString" gorm:"column:searchString;type:varchar(100);"` //用户搜索字符串
}

func (s *SearchHistory) TableName() string {
	return "SearchHistory"
}

func (s *SearchHistory) Create(DB *gorm.DB) error {
	result := DB.Create(&s)
	return result.Error
}
