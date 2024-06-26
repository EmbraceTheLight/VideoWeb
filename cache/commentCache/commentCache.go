package commentCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/cache"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
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

// GetCommentInfo 获得redis中某个缓存的评论信息
func GetCommentInfo(ctx context.Context, userID, commentID int64) (map[string]string, error) {
	mp, err := cache.HGetAll(
		ctx,
		strconv.FormatInt(userID, 10)+strconv.FormatInt(commentID, 10),
	)
	if err != nil {
		return nil, fmt.Errorf("GetCommentInfo:%w", err)
	}
	return mp, nil
}

// GetUserCommentsInfo gets many comments info from redis
func GetUserCommentsInfo(ctx context.Context, userID int64, commentIDs []int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)
	cmds := make([]*redis.MapStringStringCmd, len(commentIDs))
	pipe := DAO.RDB.Pipeline()
	for i, commentID := range commentIDs {
		key := strconv.FormatInt(userID, 10) + strconv.FormatInt(commentID, 10)
		cmds[i] = pipe.HGetAll(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetUserCommentsInfo:%w", err)
	}

	for i, cmd := range cmds {
		var commentInfo = make(map[string]string)
		commentInfo, err = cmd.Result()
		if err != nil || len(commentInfo) == 0 {
			var comments = &EntitySets.Comments{}
			comments, err = EntitySets.GetCommentByCommentID(DAO.DB, commentIDs[i])
			if err != nil {
				return nil, fmt.Errorf("GetUserCommentInfo:%w", err)
			}

			err = MakeCommentInfos(ctx, userID, comments)
			if err != nil {
				return nil, fmt.Errorf("GetUserCommentInfo:%w", err)
			}
			commentInfo, err = cache.HGetAll(ctx, strconv.FormatInt(userID, 10)+strconv.FormatInt(commentIDs[i], 10))
			if err != nil {
				return nil, fmt.Errorf("GetUserCommentInfo:%w", err)
			}
		}
		res = append(res, commentInfo)
	}
	return
}
