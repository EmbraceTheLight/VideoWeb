package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"github.com/gin-gonic/gin"
)

// AddVideoBarrage 添加视频弹幕
func AddVideoBarrage(c *gin.Context, VID int64, barrage *EntitySets.Barrage) error {
	var err error
	Utilities.AddFuncName(c, "AddVideoBarrage")

	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx = tx.Set("gorm:query_option", "FOR UPDATE")
	/*更新视频弹幕信息*/
	//添加弹幕记录
	err = EntitySets.InsertBarrageRecord(tx, barrage)
	if err != nil {
		return err
	}
	//更新视频弹幕数
	err = helper.UpdateVideoFieldForUpdate(VID, "cnt_barrage", 1, tx)
	if err != nil {
		return err
	}
	//更新视频热度
	err = helper.UpdateVideoFieldForUpdate(VID, "hot", define.AddHotEachBarrage, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
