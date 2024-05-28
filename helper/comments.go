package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
)

// GetCommentReplies 根据To字段获取当前评论的所有回复评论
func GetCommentReplies(videoID, to int64) (ret []*EntitySets.CommentSummary, err error) {
	err = DAO.DB.Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `to` = ?", videoID, to).Order("likes DESC").Find(&ret).Error
	for _, comment := range ret {
		replies, err := GetCommentReplies(videoID, comment.UID)
		if err != nil {
			return nil, err
		}
		comment.Replies = replies
	}
	return
}

// GetRootCommentsSummariesByVideoID 获得Root评论，即该评论不是回复其他评论的评论
func GetRootCommentsSummariesByVideoID(videoID int64, order string, Page, CommentsNumbers int64) (ret []*EntitySets.CommentSummary, err error) {
	if order == "default" || order == "likes" {
		err = DAO.DB.Debug().Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `to` = ?", videoID, -1).
			Order("likes DESC").Offset(int(Page)).Limit(int(CommentsNumbers)).Find(&ret).Error
	} else if order == "newest" {
		err = DAO.DB.Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `To` = ?", videoID, -1).
			Order("created_at DESC").Offset(int(Page)).Limit(int(CommentsNumbers)).Find(&ret).Error
	}
	return
}
