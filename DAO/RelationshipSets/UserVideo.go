package DAO

import (
	"VideoWeb/define"
	"fmt"
	"gorm.io/gorm"
)

type UserVideo struct {
	define.MyModel
	UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;index:idx_uid_vid"`
	VID int64 `json:"vid" gorm:"primaryKey;column:video_id;type:bigint;index:idx_uid_vid"`

	IsLike    bool `json:"is_like" gorm:"column:is_like;type:tinyint(1);default:0"`
	IsDislike bool `json:"is_unlike" gorm:"column:is_dislike;type:tinyint(1);default:0"`
	IsFavor   bool `json:"is_favor" gorm:"column:is_favor;type:tinyint(1);default:0"`
}

func (*UserVideo) TableName() string {
	return "user_video"
}

// InsertUserVideoRecord 插入用户视频记录
func InsertUserVideoRecord(db *gorm.DB, uv *UserVideo) error {
	err := db.Model(&UserVideo{}).Save(uv).Error
	if err != nil {
		return fmt.Errorf("UserVideo.InsertUserVideo: %w", err)
	}
	return nil
}

// GetUserVideoRecord 根据用户ID和视频ID获取用户视频记录
func GetUserVideoRecord(db *gorm.DB, uid int64, vid int64) (*UserVideo, error) {
	var uv = new(UserVideo)
	err := db.Model(&UserVideo{}).Where("user_id =? AND video_id = ?", uid, vid).First(uv).Error
	if err != nil {
		return nil, err
	}
	return uv, nil
}

// UpdateUserVideoRecord 更新用户视频记录的三个bool状态
func UpdateUserVideoRecord(db *gorm.DB, UID, VID int64, field string, change bool) error {
	return db.Model(&UserVideo{}).Where("user_id =? AND video_id = ?", UID, VID).Update(field, change).Error
}

// DeleteUserVideoRecordsByUserID 根据用户ID删除用户所有观看过的视频记录
func DeleteUserVideoRecordsByUserID(db *gorm.DB, UserID int64) error {
	return db.Model(&UserVideo{}).Where("user_id =?", UserID).Delete(&UserVideo{}).Error
}

// DeleteUserVideoRecordsByVideoID 根据视频ID删除所有观看过该视频的用户记录
func DeleteUserVideoRecordsByVideoID(db *gorm.DB, videoID int64) error {
	return db.Model(&UserVideo{}).Where("video_id =?", videoID).Delete(&UserVideo{}).Error
}
