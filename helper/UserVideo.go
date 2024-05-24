package helper

import (
	"VideoWeb/DAO"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"gorm.io/gorm"
)

// UpdateUserVideoFieldForUpdate 更新用户-视频联系表某个布尔字段
// 注意:如果上层逻辑只需要更新一个数据,则传入tx为nil,该函数自动开启事务进行处理
// 否则,函数调用者(位于logic层)需要自行传入tx,并在函数结束后提交或回滚事务
func UpdateUserVideoFieldForUpdate(UserID, VideoID int64, field string, change bool, tx *gorm.DB) error {
	if tx == nil {
		tx = DAO.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		tx = tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
		err := RelationshipSets.UpdateUserVideoRecord(tx, UserID, VideoID, field, change)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		err := RelationshipSets.UpdateUserVideoRecord(tx, UserID, VideoID, field, change)
		if err != nil {
			return err
		}
	}
	return nil
}
