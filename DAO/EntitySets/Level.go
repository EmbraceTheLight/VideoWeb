package DAO

import (
	"gorm.io/gorm"
)

type Level struct {
	UID       string `json:"UID" gorm:"column:uid;type:char(36);primaryKey"`
	UserLevel int    `json:"userLevel" gorm:"column:UserLevel;type:int;default:0"` //当前用户等级
	Exp       int    `json:"exp" gorm:"column:Exp;type:int;default:0"`             //当前用户在该等级已获得的经验值
	NextEXP   int    `json:"nextExp" gorm:"column:NextExp;type:int;default:2"`     //升至下一级所需总经验
}

func (l *Level) TableName() string {
	return "Level"
}

func InsertLevelRecords(tx *gorm.DB, l *Level) error {
	result := tx.Create(&l)
	return result.Error
}
