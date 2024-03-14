package DAO

import (
	"VideoWeb/define"
)

type FavoriteVideo struct {
	define.MyModel

	FavoriteID string `json:"favoriteID" gorm:"primaryKey;column:favoriteID;type:char(36);index:index_FID"`
	VideoID    string `json:"videoID" gorm:"primaryKey;column:videoID;type:char(36);index:index_VID"`
}

func (f *FavoriteVideo) TableName() string {
	return "FavoriteVideo"
}
