package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type Level struct {
	UserLevel uint16 `json:"userLevel" gorm:"column:user_level;type:smallint unsigned;default:0"` //当前用户等级
	Exp       uint16 `json:"exp" gorm:"column:cur_exp;type:smallint unsigned;default:0"`          //当前用户在该等级已获得的经验值
	NextEXP   uint16 `json:"nextExp" gorm:"column:next_exp;type:smallint unsigned;default:50"`    //升至下一级所需总经验
	UID       int64  `json:"UID" gorm:"column:user_id;type:bigint ;primaryKey"`
}

func (level *Level) TableName() string {
	return "user_level"
}

func (level *Level) HandleLevelUp() {
	if level.UserLevel == 6 {
		return
	}

	//改变经验值
	level.Exp = level.Exp - level.NextEXP
	//用户等级+1
	level.UserLevel += 1
	//设置下一级所需经验值
	switch level.UserLevel {
	case 1:
		level.NextEXP = define.ToLevel2
	case 2:
		level.NextEXP = define.ToLevel3
	case 3:
		level.NextEXP = define.ToLevel4
	case 4:
		level.NextEXP = define.ToLevel5
	case 5:
		level.NextEXP = define.ToLevel6
	}

}

func (level *Level) AddExp(addNum int) {
	level.Exp += uint16(addNum)
	if level.Exp >= level.NextEXP {
		level.HandleLevelUp()
	}
}

// GetLevelByUserID 获取用户等级信息
func GetLevelByUserID(db *gorm.DB, userID int64) (*Level, error) {
	var level *Level
	err := db.Model(&Level{}).Where("user_id = ?", userID).First(&level).Error
	return level, err
}

// SaveLevelRecords 保存或更新用户等级信息
func SaveLevelRecords(db *gorm.DB, l *Level) error {
	return db.Model(&Level{}).Where("user_id = ?", l.UID).Save(&l).Error
}

// DeleteLevelRecordByUserID 删除用户等级信息
func DeleteLevelRecordByUserID(db *gorm.DB, uid int64) error {
	return db.Model(&Level{}).Where("user_id = ?", uid).Delete(&Level{}).Error
}

// GetLevelRecordByUserID 获取用户等级信息
func GetLevelRecordByUserID(db *gorm.DB, uid int64) (l *Level, err error) {
	err = db.Model(&Level{}).Where("user_id = ?", uid).First(&l).Error
	return
}
