package DAO

import (
	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	UID          string `json:"UID" gorm:"column:UID;type:char(36);primaryKey"`
	SearchString string `json:"searchString" gorm:"column:searchString;type:varchar(100);"` //用户搜索字符串
}

func (s *SearchHistory) TableName() string {
	return "SearchHistory"
}
