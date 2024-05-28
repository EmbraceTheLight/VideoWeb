package DAO

import (
	"gorm.io/gorm"
)

type FavoriteVideo struct {
	gorm.Model
	UserID     int64 `json:"userID" gorm:"column:user_id;type:bigint;index:index_UID"`
	FavoriteID int64 `json:"favoriteID" gorm:"column:favorite_id;type:bigint;index:index_FID"`
	VideoID    int64 `json:"videoID" gorm:"column:video_id;type:bigint;index:index_VID"`
}

func (f *FavoriteVideo) TableName() string {
	return "favorite_video"
}

func InsertFavoriteVideoRecord(db *gorm.DB, fv *FavoriteVideo) error {
	return db.Model(FavoriteVideo{}).Create(fv).Error
}

// DeleteFavoriteVideoRecordsByFavoriteID 根据收藏夹ID删除收藏记录，涉及到删除一个收藏夹的操作
func DeleteFavoriteVideoRecordsByFavoriteID(db *gorm.DB, fid int64) error {
	return db.Model(&FavoriteVideo{}).Where("favorite_id = ?", fid).Delete(&FavoriteVideo{}).Error
}

// DeleteFavoriteVideoRecordsByUserID 根据用户ID删除收藏记录，涉及到注销一个用户时的操作
func DeleteFavoriteVideoRecordsByUserID(db *gorm.DB, uid int64) error {
	return db.Model(&FavoriteVideo{}).Where("user_id = ?", uid).Delete(&FavoriteVideo{}).Error
}

// DeleteFavoriteVideoRecordByUserIDVideoID 根据用户ID和视频ID删除收藏记录，涉及到删除一个收藏记录的操作
func DeleteFavoriteVideoRecordByUserIDVideoID(db *gorm.DB, uid, vid int64) error {
	return db.Model(&FavoriteVideo{}).Where("user_id = ? and video_id = ?", uid, vid).Delete(&FavoriteVideo{}).Error
}
