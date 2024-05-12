package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"gorm.io/gorm"
)

// UpdateUserFieldForUpdate 更新用户某个数值字段(悲观锁)
// 注意:如果上层逻辑只需要更新一个数据,则传入tx为nil,该函数自动开启事务进行处理
// 否则,函数调用者(位于logic层)需要自行传入tx,并在函数结束后提交或回滚事务
func UpdateUserFieldForUpdate(UserID int64, field string, change int, tx *gorm.DB) error {
	if tx != nil {
		tx = DAO.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
		err := EntitySets.UpdateUserField(tx, UserID, field, change)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		err := EntitySets.UpdateUserField(tx, UserID, field, change)
		if err != nil {
			return err
		}
	}

	return nil
}
