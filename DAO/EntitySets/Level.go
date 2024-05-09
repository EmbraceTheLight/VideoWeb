package DAO

import (
	"gorm.io/gorm"
)

type Level struct {
	UserLevel uint16 `json:"userLevel" gorm:"column:UserLevel;type:smallint unsigned;default:0"` //当前用户等级
	Exp       uint16 `json:"exp" gorm:"column:Exp;type:smallint unsigned;default:0"`             //当前用户在该等级已获得的经验值
	NextEXP   uint16 `json:"nextExp" gorm:"column:NextExp;type:smallint unsigned;default:2"`     //升至下一级所需总经验
	UID       int64  `json:"UID" gorm:"column:uid;type:bigint ;primaryKey"`
}

func (l *Level) TableName() string {
	return "Level"
}

func InsertLevelRecords(db *gorm.DB, l *Level) error {
	result := db.Create(&l)
	return result.Error
}
