package logic

import (
	"VideoWeb/DAO"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/helper"
	"errors"
	"gorm.io/gorm"
)

// InsertUserVideoIfNotExist 如果不存在用户-视频记录,则插入新记录
func InsertUserVideoIfNotExist(uid, vid int64) error {
	vu1, err1 := RelationshipSets.GetUserVideoRecord(DAO.DB, uid, vid)
	if vu1 == nil {
		//未找到记录,则插入新记录
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			vu2 := &RelationshipSets.UserVideo{
				UID: uid,
				VID: vid,
			}
			err2 := RelationshipSets.InsertUserVideoRecord(DAO.DB, vu2)
			if err2 != nil {
				return err2
			}
		} else { // 其他错误
			return err1
		}
	}
	return nil
}

// UpdateUserVideoIsLike 更新用户-视频记录的点赞状态
func UpdateUserVideoIsLike(uid, vid int64, isLike bool, tx *gorm.DB) error {
	return helper.UpdateUserVideoFieldForUpdate(uid, vid, "is_like", isLike, tx)
}
