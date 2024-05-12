package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type UserVideo struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:uid;type:bigint;index:idx_uid_vid"`
	VID int64 `json:"vid" gorm:"primaryKey;column:vid;type:bigint;index:idx_uid_vid"`

	IsLike   bool `json:"is_like" gorm:"column:is_like;type:tinyint(1);default:0"`
	IsUnlike bool `json:"is_unlike" gorm:"column:is_unlike;type:tinyint(1);default:0"`
	IsFavor  bool `json:"is_favor" gorm:"column:is_favor;type:tinyint(1);default:0"`
}

func (*UserVideo) TableName() string {
	return "user_video"
}

func InsertUserVideoRecord(db *gorm.DB, uv *UserVideo) error {
	return db.Model(&UserVideo{}).Create(uv).Error
}

func GetUserVideoRecord(db *gorm.DB, uid int64, vid int64) (*UserVideo, error) {
	var uv = new(UserVideo)
	err := db.Debug().Model(&UserVideo{}).Where("uid =? AND vid = ?", uid, vid).First(uv).Error
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func UpdateUserVideoRecord(db *gorm.DB, UID, VID int64, field string, change bool) error {
	return db.Model(&UserVideo{}).Where("uid =? AND vid = ?", UID, VID).Update(field, change).Error
}
