package DAO

import (
	"VideoWeb/define"
	"fmt"
	"gorm.io/gorm"
)

type UserSearchHistory struct {
	define.MyModel
	UID          int64  `json:"UID" gorm:"column:user_id;type:bigint;primaryKey"`
	SearchString string `json:"searchString" gorm:"column:search_string;type:varchar(100);"` //用户搜索字符串
}

func (s *UserSearchHistory) TableName() string {
	return "user_search_history"
}

func (s *UserSearchHistory) GetScore() float64 {
	return float64(s.UpdatedAt.UnixMilli())
}

func (s *UserSearchHistory) GetValue() any {
	return s.SearchString
}

// InsertSearchRecord 插入或更新用户搜索记录
func InsertSearchRecord(db *gorm.DB, UID int64, searchString string) error {
	sh := &UserSearchHistory{
		UID:          UID,
		SearchString: searchString,
	}
	err := db.Where("user_id = ?", UID).Omit("created_at").Save(sh).Error
	if err != nil {
		return fmt.Errorf("UserSearch.InsertSearchRecord: %w", err)
	}
	return nil
}

// DeleteOneSearchRecord 根据记录ID删除单条搜索记录
func DeleteOneSearchRecord(db *gorm.DB, ID uint) error {
	return db.Model(&UserSearchHistory{}).Where("id = ?", ID).Delete(&UserSearchHistory{}).Error
}

// DeleteAllSearchRecord 根据用户ID删除用户的所有搜索记录
func DeleteAllSearchRecord(db *gorm.DB, uid int64) error {
	return db.Model(&UserSearchHistory{}).Where("user_id = ?", uid).Delete(&UserSearchHistory{}).Error
}
