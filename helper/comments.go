package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"gorm.io/gorm"
)

// GetCommentReplies 根据To字段获取当前评论的所有回复评论
func GetCommentReplies(videoID, to int64) (ret []*EntitySets.CommentSummary, err error) {
	err = DAO.DB.Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `to` = ?", videoID, to).Order("likes DESC").Find(&ret).Error
	for _, comment := range ret {
		replies, err := GetCommentReplies(videoID, comment.CommentID)
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
		err = DAO.DB.Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `to` = ?", videoID, -1).
			Order("likes DESC").Offset(int(Page)).Limit(int(CommentsNumbers)).Find(&ret).Error
	} else if order == "newest" {
		err = DAO.DB.Model(&EntitySets.Comments{}).Where("`video_id` = ? AND `to` = ?", videoID, -1).
			Order("created_at DESC").Offset(int(Page)).Limit(int(CommentsNumbers)).Find(&ret).Error
	}
	return
}

// UpdateComment 更新评论的点赞/踩数
func UpdateComment(commentID int64, isLike, isUndo bool, tx *gorm.DB) error {
	var err error
	if !isUndo { //不是撤销，增加点赞/踩数
		if isLike {
			err = EntitySets.UpdateCommentField(tx, commentID, "likes", +1)
		} else {
			err = EntitySets.UpdateCommentField(tx, commentID, "dislikes", +1)
		}
	} else {
		if isLike {
			err = EntitySets.UpdateCommentField(tx, commentID, "likes", -1)
		} else {
			err = EntitySets.UpdateCommentField(tx, commentID, "dislikes", -1)
		}
	}
	return err
}

// UpdateUserCommentRecord 根据用户的点赞/踩操作，插入/删除用户点赞/踩状态
func UpdateUserCommentRecord(uid, cid, vid int64, isLike, isUndo bool, tx *gorm.DB) error {
	var err error
	if isUndo { //若是撤销操作，则删除用户点赞/踩记录，否则插入
		if isLike {
			err = RelationshipSets.DeleteUserLikedCommentRecord(tx, uid, cid)
		} else {
			err = RelationshipSets.DeleteUserDislikedCommentRecord(tx, uid, cid)
		}
	} else {
		if isLike {
			err = RelationshipSets.InsertUserLikedCommentRecord(tx, uid, vid, cid)
		} else {
			err = RelationshipSets.InsertUserDislikedCommentRecord(tx, uid, vid, cid)
		}
	}
	return err
}

// GetUserCommentRecords 获取用户对评论的点赞/踩状态(只有赞/踩过的评论才有状态记录)
func GetUserCommentRecords(uid, vid int64, tx *gorm.DB) (likes, dislikes map[int64]bool, err error) {
	l, err := RelationshipSets.GetUserLikedCommentRecordByUidVid(tx, uid, vid)
	dl, err := RelationshipSets.GetUserDislikedCommentRecordByUidVid(tx, uid, vid)
	likes = make(map[int64]bool)
	dislikes = make(map[int64]bool)
	for _, like := range l {
		likes[like.CID] = true
	}
	for _, dislike := range dl {
		dislikes[dislike.CID] = true
	}
	return
}

// UpdateCommentsStatus 更新CommentSummary的like/dislike状态
func UpdateCommentsStatus(likes, dislikes map[int64]bool, comments []*EntitySets.CommentSummary) {
	for _, cs := range comments {
		if likes[cs.CommentID] {
			cs.Like = true
		} else if dislikes[cs.CommentID] {
			cs.Dislike = true
		}
		if len(cs.Replies) > 0 {
			UpdateCommentsStatus(likes, dislikes, cs.Replies)
		}
	}
}
