package commentCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/cache"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// MakeCommentInfos saves comments of a user or a video info to redis,prefix is user_id/video_id,use commentCache.
func MakeCommentInfos(ctx context.Context, prefix int64, comments ...*EntitySets.Comments) error {
	cmts := make([]*CommentCache, len(comments))
	HashMap := make([]cache.HashMap, len(comments))
	for i, comment := range comments {
		cmts[i] = new(CommentCache)
		cmts[i].comments = make(map[string]any)
		cmts[i].key = strconv.FormatInt(prefix, 10) + strconv.FormatInt(comment.CommentID, 10)
		cmts[i].comments["created_at"] = comment.CreatedAt.Format("2006-01-02 15:04:05")
		cmts[i].comments["comment_id"] = comment.CommentID
		cmts[i].comments["user_id"] = comment.UID
		cmts[i].comments["to"] = comment.To
		cmts[i].comments["video_id"] = comment.VID
		cmts[i].comments["content"] = comment.Content
		cmts[i].comments["likes"] = comment.Likes
		cmts[i].comments["dislikes"] = comment.Dislikes
		cmts[i].comments["ip_address"] = comment.IPAddress
		HashMap[i] = cmts[i]
	}

	// 缓存对应用户-评论信息到redis中
	err := cache.HSets(ctx, cache.CommentExpireTime, HashMap...)
	if err != nil {
		return fmt.Errorf("MakeCommentInfo:%w", err)
	}

	return nil
}

// AddCommentInfo adds a comment info to redis, and update the user's comment list
func AddCommentInfo(ctx context.Context, id int64, comment *EntitySets.Comments) (err error) {
	err = MakeCommentInfos(ctx, id, comment)
	if err != nil {
		return fmt.Errorf("commentCache.AddCommentInfo::%w", err)
	}
	err = cache.SAddWithRetry(
		ctx,
		strconv.FormatInt(id, 10)+"_comments",
		cache.DefaultTry,
		cache.DefaultSleep,
		cache.CommentExpireTime,
		comment.CommentID,
	)
	if err != nil {
		return fmt.Errorf("commentCache.AddCommentInfo::%w", err)
	}
	return nil
}

// GetUserCommentsInfo gets many comments info from redis
func GetUserCommentsInfo(ctx context.Context, userID int64, commentIDs []int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)
	cmds := make([]*redis.MapStringStringCmd, len(commentIDs))
	pipe := DAO.RDB.Pipeline()
	for i, commentID := range commentIDs {
		key := strconv.FormatInt(userID, 10) + strconv.FormatInt(commentID, 10)
		cmds[i] = pipe.HGetAll(ctx, key)
		pipe.Expire(ctx, key, cache.CommentExpireTime)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetUserCommentsInfo:%w", err)
	}

	for i, cmd := range cmds {
		var commentInfo = make(map[string]string)
		commentInfo, err = cmd.Result()
		if err != nil || len(commentInfo) == 0 {
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			//处理缓存失效的问题
			var comments = &EntitySets.Comments{}
			comments, err = EntitySets.GetCommentByCommentID(DAO.DB, commentIDs[i])
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			err = MakeCommentInfos(ctx, userID, comments)
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			commentInfo, err = cache.HGetAll(
				ctx,
				strconv.FormatInt(userID, 10)+strconv.FormatInt(commentIDs[i], 10),
				cache.CommentExpireTime,
			)
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}
		}
		res = append(res, commentInfo)
	}
	return
}

// MapStringStringToComment 将map[string]string改成[]*EntitySets.CommentSummary
func MapStringStringToComment(comment map[string]string) (cmt *EntitySets.CommentSummary) {
	cmt = new(EntitySets.CommentSummary)
	cmt.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", comment["created_at"])
	cmt.CommentID = Utilities.String2Int64(comment["comment_id"])
	cmt.UID = Utilities.String2Int64(comment["user_id"])
	cmt.To = Utilities.String2Int64(comment["to"])
	cmt.VID = Utilities.String2Int64(comment["video_id"])
	cmt.Content = comment["content"]
	cmt.Likes = Utilities.String2Uint32(comment["likes"])
	cmt.Dislikes = Utilities.String2Uint32(comment["dislikes"])
	cmt.IPAddress = comment["ip_address"]
	return
}

// MakeUserLikedComments 制作用户针对某视频的点赞评论缓存，缓存的都为点赞的评论ID
func MakeUserLikedComments(ctx context.Context, userID int64, videoID int64) (err error) {
	var likedCommentIDs []any
	err = DAO.DB.Model(&RelationshipSets.UserLikedComments{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).Select("comment_id").Find(&likedCommentIDs).Error
	if err != nil {
		return fmt.Errorf("CommentCache.MakeUserLikedComments: %w", err)
	}

	key := strconv.FormatInt(userID, 10) + strconv.FormatInt(videoID, 10) + "_liked_comments"
	if len(likedCommentIDs) != 0 {
		err = cache.SAddWithRetry(
			ctx, key,
			cache.DefaultTry, cache.DefaultSleep,
			cache.CommentExpireTime,
			likedCommentIDs...,
		)
		if err != nil {
			return fmt.Errorf("CommentCache.MakeUserLikedComments::%w", err)
		}
	}
	return nil
}

// GetUserLikedComments gets user liked comments of a video from redis
func GetUserLikedComments(ctx context.Context, userID int64, videoID int64) (res []string, err error) {
	key := strconv.FormatInt(userID, 10) + strconv.FormatInt(videoID, 10) + "_liked_comments"
	res, err = cache.SMembers(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("GetUserLikedComments::%w", err)
	}

	if len(res) == 0 {
		err = MakeUserLikedComments(ctx, userID, videoID)
		if err != nil {
			return nil, fmt.Errorf("CommentCache.GetUserLikedComments::%w", err)
		}
	}
	return
}
