package commentCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/cache"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// MakeCommentInfos saves comments of a user or a video info to redis,prefix is user_id/video_id,use commentCache.
func MakeCommentInfos(ctx context.Context, comments ...*EntitySets.Comments) error {
	cmts := make([]*CommentCache, len(comments))
	HashMap := make([]cache.HashMap, len(comments))

	for i, comment := range comments {
		cmts[i] = new(CommentCache)
		cmts[i].comments = make(map[string]any)
		cmts[i].key = strconv.FormatInt(comment.CommentID, 10)
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

	// 缓存对应评论信息到redis中
	err := cache.HSets(ctx, cache.CommentExpireTime, HashMap...)
	if err != nil {
		return fmt.Errorf("MakeCommentInfo::%w", err)
	}

	return nil
}

// AddCommentInfo adds a comment info to redis, and update the user's comment list
func AddCommentInfo(ctx context.Context, videoID, userID int64, comment *EntitySets.Comments) (err error) {
	err = addCommentInfo(ctx, videoID, comment) //缓存到视频评论信息
	if err != nil {
		return fmt.Errorf("commentCache.AddCommentInfo::%w", err)
	}

	err = addCommentInfo(ctx, userID, comment) //缓存到用户评论信息
	if err != nil {
		return fmt.Errorf("commentCache.AddCommentInfo::%w", err)
	}

	return nil
}

// GetCommentsInfo gets many comments info from redis
func GetCommentsInfo(ctx context.Context, commentIDs []int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)
	cmds := make([]*redis.MapStringStringCmd, len(commentIDs))

	pipe := DAO.RDB.Pipeline()
	for i, commentID := range commentIDs {
		key := strconv.FormatInt(commentID, 10)
		cmds[i] = pipe.HGetAll(ctx, key)
		pipe.Expire(ctx, key, cache.CommentExpireTime)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetCommentsInfo:%w", err)
	}

	for i, cmd := range cmds {
		var commentInfo = make(map[string]string)
		commentInfo, err = cmd.Result()
		if err != nil || len(commentInfo) == 0 {
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			//len(commentInfo) == 0，代表缓存。失效处理缓存失效的问题。
			var comments = &EntitySets.Comments{}
			comments, err = EntitySets.GetCommentByCommentID(DAO.DB, commentIDs[i])
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			err = MakeCommentInfos(ctx, comments)
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}

			commentInfo, err = cache.HGetAll(ctx, strconv.FormatInt(commentIDs[i], 10), cache.CommentExpireTime)
			if err != nil {
				return nil, fmt.Errorf("CommentCache.GetUserCommentInfo: %w", err)
			}
		}
		res = append(res, commentInfo)
	}
	return
}

// MapStringString2Comments 将map[string]string改成[]*EntitySets.CommentSummary
func MapStringString2Comments(comments ...map[string]string) (res []*EntitySets.CommentSummary) {
	for _, c := range comments {
		var cmt = new(EntitySets.CommentSummary)
		cmt.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", c["created_at"])
		cmt.CommentID = Utilities.String2Int64(c["comment_id"])
		cmt.UID = Utilities.String2Int64(c["user_id"])
		cmt.To = Utilities.String2Int64(c["to"])
		cmt.VID = Utilities.String2Int64(c["video_id"])
		cmt.Content = c["content"]
		cmt.Likes = Utilities.String2Uint32(c["likes"])
		cmt.Dislikes = Utilities.String2Uint32(c["dislikes"])
		cmt.IPAddress = c["ip_address"]
		res = append(res, cmt)
	}
	return
}
