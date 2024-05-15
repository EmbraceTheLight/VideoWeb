package DAO

import (
	"gorm.io/gorm"
)

type Level struct {
	UserLevel uint16 `json:"userLevel" gorm:"column:user_level;type:smallint unsigned;default:0"` //当前用户等级
	Exp       uint16 `json:"exp" gorm:"column:cur_exp;type:smallint unsigned;default:0"`          //当前用户在该等级已获得的经验值
	NextEXP   uint16 `json:"nextExp" gorm:"column:next_exp;type:smallint unsigned;default:2"`     //升至下一级所需总经验
	UID       int64  `json:"UID" gorm:"column:user_id;type:bigint ;primaryKey"`
}

func (l *Level) TableName() string {
	return "user_level"
}

func InsertLevelRecords(db *gorm.DB, l *Level) error {
	return db.Model(&Level{}).Create(&l).Error
}

func DeleteLevelRecordByUserID(db *gorm.DB, uid int64) error {
	return db.Model(&Level{}).Where("user_id = ?", uid).Delete(&Level{}).Error
}
