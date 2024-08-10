package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&EntitySets.Video{}, comment.VID)
	err = helper.UpdateVideoFieldForUpdate(comment.VID, "hot", define.AddHotEachComment, tx)
	tx.Commit()
	return nil
}

// LikeOrDislikeComment 点赞或点踩评论
func LikeOrDislikeComment(c *gin.Context, userID, commentID, videoID int64, isLike bool) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "LikeOrDislikeComment")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	/*点赞或点踩评论*/
	//更新评论点赞/点踩数
	tx.Set("gorm:query_option", "FOR UPDATE")
	err = helper.UpdateComment(commentID, isLike, false, tx)
	if err != nil {
		return err
	}
	//插入点赞/点踩记录
	err = helper.UpdateUserCommentRecord(userID, commentID, videoID, isLike, false, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// UndoLikeOrDislikeComment 取消点赞或取消点踩评论
func UndoLikeOrDislikeComment(c *gin.Context, userID, commentID, videoID int64, isLike bool) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "LikeOrDislikeComment")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE")

	/*点赞或点踩评论*/
	//更新评论点赞/点踩数
	err = helper.UpdateComment(commentID, isLike, true, tx)
	if err != nil {
		return err
	}

	//删除点赞/点踩记录
	err = helper.UpdateUserCommentRecord(userID, commentID, videoID, isLike, true, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
