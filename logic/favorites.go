package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
)

func String2Int8(IsPrivate string) int8 {
	var ret int8
	if IsPrivate == "公开" {
		ret = 1
	} else if IsPrivate == "私密" {
		ret = -1
	}
	return ret
}

// DeleteFavoritesRecordsByNameUserID 删除收藏夹及其包含的所有收藏记录
func DeleteFavoritesRecordsByNameUserID(FName string, UID int64) error {
	tx := DAO.DB.Begin()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//查询待删除的收藏夹ID
	f, err := EntitySets.GetFavoriteRecordByNameUserID(tx, FName, UID)
	if err != nil {
		return err
	}
	//删除收藏记录
	err = RelationshipSets.DeleteFavoriteVideoRecordsByFavoriteID(tx, f.FavoriteID)
	if err != nil {
		return err
	}
	//删除收藏夹
	err = EntitySets.DeleteFavoritesRecordsByNameUserID(tx, FName, UID)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
