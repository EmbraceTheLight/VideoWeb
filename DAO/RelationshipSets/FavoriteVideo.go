package DAO

import (
	"VideoWeb/define"
)

type FavoriteVideo struct {
	define.MyModel

	FavoriteID int64 `json:"favoriteID" gorm:"primaryKey;column:favoriteID;type:bigint;index:index_FID"`
	VideoID    int64 `json:"videoID" gorm:"primaryKey;column:videoID;type:bigint;index:index_VID"`
}

func (f *FavoriteVideo) TableName() string {
	return "FavoriteVideo"
}
