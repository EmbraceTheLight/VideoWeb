package userCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/cache"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// MakeUserCache 创建一个新的UserCache
func MakeUserCache() *UserCache {
	return &UserCache{
		Ub:  UserBasic{Userinfo: make(map[string]any)},
		Uc:  UserComments{UserComments: make([]int64, 0)},
		Ufs: UserFollows{UserFollows: make([]int64, 0)},
		Ufd: UserFollowed{UserFollowed: make([]int64, 0)},
		Uv:  UserVideo{UserVideo: make([]int64, 0)},
		Us:  UserSearch{UserSearch: make([]string, 0)},
		Uw:  UserWatch{UserWatch: make([]int64, 0)},
	}
}

// MakeUserinfo 给定一个用户信息，创建该用户的缓存信息
func (user *UserCache) MakeUserinfo(ctx context.Context, userID int64) (err error) {
	var eg errgroup.Group
	if DAO.RDB.TTL(ctx, strconv.FormatInt(userID, 10)).Val() > 0 {
		return nil
	}

	// user 基本信息
	eg.Go(func() error {
		return user.Ub.makeBasicInfo(ctx, userID)
	})

	// user 粉丝信息
	eg.Go(func() error { return user.Ufd.makeFollowedInfo(ctx, userID) })

	// user 关注的人的信息
	eg.Go(func() error { return user.Ufs.makeFollowsInfo(ctx, userID) })

	// user 历史评论信息
	eg.Go(func() error { return user.Uc.makeCommentsInfo(ctx, userID) })

	// user 历史观看记录信息
	eg.Go(func() error { return user.Uw.makeUserWatch(ctx, userID) })

	// user 历史搜索记录信息
	eg.Go(func() error { return user.Us.makeUserSearch(ctx, userID) })

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("MakeUserinfo::%w", err)
	}
	return nil
}

func DeleteUserCache(ctx context.Context, userID int64) (err error) {
	return nil
}

// GetUserBasicInfo 获得redis中某个用户的基本信息
func GetUserBasicInfo(ctx context.Context, userID int64) (res map[string]string, err error) {
	key := strconv.FormatInt(userID, 10)

	if checkUserCache(ctx, key, userID) != nil {
		return nil, fmt.Errorf("userCache.GetUserBasicInfo::%w", err)
	}

	res, err = cache.HGetAll(ctx, key, cache.UserExpireTime)
	if err != nil {
		return nil, fmt.Errorf("userCache.GetUserBasicInfo::%w", err)
	}

	return res, nil
}

// GetUsersBasicInfo Gets many users' information return find result and an error.
func GetUsersBasicInfo(ctx context.Context, userIDs []int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)

	// 存储每个 HGetAll 命令的结果
	cmds := make([]*redis.MapStringStringCmd, len(userIDs))

	pipe := DAO.RDB.Pipeline()
	for i, userID := range userIDs {
		key := strconv.FormatInt(userID, 10)
		cmds[i] = pipe.HGetAll(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("userCache.GetUsersBasicInfo::%w", err)
	}

	// 遍历每个命令的结果，如果为空，说明缓存的数据可能过期，则从数据库中获取用户信息
	for i, cmd := range cmds {
		var userInfo = make(map[string]string)
		userInfo, err = cmd.Result()
		if err != nil || len(userInfo) == 0 {
			if err != nil {
				return nil, fmt.Errorf("userCache.GetUserBasicInfo::%w", err)
			}

			key := strconv.FormatInt(userIDs[i], 10)
			uc := MakeUserCache()
			err = uc.MakeUserinfo(ctx, userIDs[i])
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
					err2 := cache.SetNilHash(ctx, key)
					if err2 != nil {
						return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::RDB.Pipelined: %w", err2)
					}
				} else {
					return nil, fmt.Errorf("userCache.GetUserBasicInfo: %w", err)
				}
			}

			userInfo, err = cache.HGetAll(ctx, key, cache.UserExpireTime)
			if err != nil {
				return nil, fmt.Errorf("userCache.GetUsersBasicInfo: %w", err)
			}
		}
		res = append(res, userInfo)
	}
	return res, nil
}

// MapStringString2Users 将map[string]string转换为*EntitySets.User结构切片
func MapStringString2Users(userBasic ...map[string]string) (res []*EntitySets.User) {
	for _, basic := range userBasic {
		var user = new(EntitySets.User)
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", basic["created_at"])
		user.UserID = Utilities.String2Int64(basic["user_id"])
		user.UserName = basic["user_name"]
		user.Password = basic["password"]
		user.Email = basic["email"]
		user.Signature = basic["signature"]
		user.Shells = Utilities.String2Uint32(basic["shells"])
		user.IsAdmin = Utilities.String2Int(basic["is_admin"])
		user.CntMsgNotRead = int32(Utilities.String2Int(basic["cnt_msg_not_read"]))
		user.CntLikes = Utilities.String2Uint32(basic["cnt_likes"])
		user.UserLevel.UserLevel = uint16(Utilities.String2Uint32(basic["user_level"]))
		user.AvatarPath = basic["user_avatar"]
	}

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

	key := strconv.FormatInt(userID, 10) + strconv.FormatInt(videoID, 10) + cache.ULCSfx
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
	key := strconv.FormatInt(userID, 10) + strconv.FormatInt(videoID, 10) + cache.ULCSfx

	if DAO.RDB.TTL(ctx, key).Val() < 0 {
		err = MakeUserLikedComments(ctx, userID, videoID)
		if err != nil {
			return nil, fmt.Errorf("CommentCache.GetUserLikedComments::%w", err)
		}
	}

	res, err = cache.SMembers(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("CommentCache.GetUserLikedComments::%w", err)
	}

	return
}

// UpdateUserInfoFields updates the information of a video in cache.
func UpdateUserInfoFields(ctx context.Context, userID int64, fields map[string]any) (err error) {
	if DAO.RDB.TTL(ctx, strconv.FormatInt(userID, 10)).Val() < 0 {
		return nil
	}

	key := strconv.FormatInt(userID, 10)
	if len(fields) == 0 {
		return nil
	}

	pipe := DAO.RDB.Pipeline()
	for k, v := range fields {
		pipe.HSet(ctx, key, k, v)
	}
	pipe.Expire(ctx, key, cache.UserExpireTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("UserCache.UpdateUserInfoFields::%w", err)
	}

	return nil
}
