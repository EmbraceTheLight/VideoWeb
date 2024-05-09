package DAO

import (
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/define"
	"gorm.io/gorm"
)

type Favorites struct {
	define.MyModel

	FavoriteID  int64                             `json:"FavoriteID" gorm:"column:FavoriteID;type:bigint;primaryKey"` //收藏夹ID
	UID         int64                             `json:"UID" gorm:"column:UID;type:bigint"`                          //收藏夹拥有者ID
	FName       string                            `json:"FName" gorm:"column:FName;type:varchar(255);not null"`       //收藏夹名称
	Description string                            `json:"Description" gorm:"column:Description;type:text"`            //收藏夹描述
	IsPrivate   int8                              `json:"IsPrivate" gorm:"column:IsPrivate;type:tinyint;default:1"`   //是否私密,1表示公开，-1表示私密
	Videos      []*RelationshipSets.FavoriteVideo `gorm:"foreignKey:FavoriteID;references:FavoriteID"`                // 收藏夹包含的视频表
}

func (f *Favorites) TableName() string {
	return "Favorites"
}

func InsertFavoritesRecords(db *gorm.DB, f *Favorites) error {
	result := db.Model(&Favorites{}).Create(&f)
	return result.Error
}

func DeleteFavoritesRecords(db *gorm.DB, f *Favorites) error {
	result := db.Model(&Favorites{}).Where("FavoriteID =?", f.FavoriteID).Delete(&f)
	return result.Error
}

func UpdateFavoritesRecords(db *gorm.DB, f *Favorites) error {
	err := db.Model(&Favorites{}).Save(&f).Error
	return err
}

func (f *Favorites) Delete(DB *gorm.DB) error {
	result := DB.Delete(&f)
	return result.Error
}
