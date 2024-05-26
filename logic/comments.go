package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"github.com/gin-gonic/gin"
)

// AddCommentToVideo 添加对视频的评论
func AddCommentToVideo(c *gin.Context, comment *EntitySets.Comments) error {
	var err error

	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "AddCommentToVideo")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 插入评论记录
	err = EntitySets.InsertCommentRecord(tx, comment)
	if err != nil {
		return err
	}
	// 更新视频热度
	err = helper.UpdateVideoFieldForUpdate(comment.VID, "hot", define.AddHotEachComment, tx.Set("gorm:query_option", "FOR UPDATE"))
	tx.Commit()
	return nil
}
