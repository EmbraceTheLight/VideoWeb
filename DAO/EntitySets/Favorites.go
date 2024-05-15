package DAO

import (
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/define"
	"gorm.io/gorm"
)

type Favorites struct {
	define.MyModel

	FavoriteID  int64                             `json:"FavoriteID" gorm:"column:favorite_id;type:bigint;primaryKey"`  //收藏夹ID
	UID         int64                             `json:"UID" gorm:"column:user_id;type:bigint;index:idx_uid"`          //收藏夹拥有者ID
	FName       string                            `json:"FName" gorm:"column:favorite_name;type:varchar(255);not null"` //收藏夹名称
	Description string                            `json:"Description" gorm:"column:description;type:text"`              //收藏夹描述
	IsPrivate   int8                              `json:"IsPrivate" gorm:"column:is_private;type:tinyint;default:1"`    //是否私密,1表示公开，-1表示私密
	Videos      []*RelationshipSets.FavoriteVideo `gorm:"foreignKey:FavoriteID;references:FavoriteID"`                  // 收藏夹包含的视频表
}

func (f *Favorites) TableName() string {
	return "favorites"
}

// GetFavoriteRecordByFavoriteID 通过收藏夹ID来获取该用户的收藏夹列表
func GetFavoriteRecordByFavoriteID(db *gorm.DB, fid int64) (*Favorites, error) {
	ret := new(Favorites)
	err := db.Model(&Favorites{}).Where("favorite_id = ?", fid).First(ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetFavoriteRecordByUserID 通过用户ID来获取该用户的收藏夹列表
func GetFavoriteRecordByUserID(db *gorm.DB, id int64) ([]*Favorites, error) {
	var favorites []*Favorites
	err := db.Where("user_id = ?", id).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// GetFavoriteRecordByNameUserID 根据收藏夹名称和用户ID获取收藏夹
func GetFavoriteRecordByNameUserID(db *gorm.DB, name string, uid int64) (*Favorites, error) {
	ret := new(Favorites)
	err := db.Model(&Favorites{}).Where("favorite_name = ? and user_id = ?", name, uid).First(ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// InsertFavoritesRecords 插入收藏夹记录
func InsertFavoritesRecords(db *gorm.DB, f *Favorites) error {
	return db.Model(&Favorites{}).Create(&f).Error
}

// DeleteFavoritesRecordsByNameUserID 根据收藏夹名称和用户ID删除收藏夹
func DeleteFavoritesRecordsByNameUserID(db *gorm.DB, name string, uid int64) error {
	return db.Model(&Favorites{}).Where("favorite_name =? AND user_id =?", name, uid).Delete(&Favorites{}).Error
}

// DeleteFavoritesRecord 删除某个收藏夹
func DeleteFavoritesRecord(db *gorm.DB, f *Favorites) error {
	return db.Model(&Favorites{}).Delete(f).Error
}

// DeleteFavoritesRecordsByUserID 删除该用户的所有收藏夹
func DeleteFavoritesRecordsByUserID(db *gorm.DB, uid int64) error {
	return db.Debug().Model(&Favorites{}).Where("user_id=?", uid).Delete(&Favorites{}).Error
}

// SaveFavoritesRecords 保存收藏夹记录
func SaveFavoritesRecords(db *gorm.DB, f *Favorites) error {
	return db.Model(&Favorites{}).Where("favorite_id =?", f.FavoriteID).Save(&f).Error
}
