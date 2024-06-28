package userCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
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
		UV:  UserVideo{UserVideo: make([]int64, 0)},
		US:  UserSearch{UserSearch: make([]string, 0)},
		UW:  UserWatch{UserWatch: make([]int64, 0)},
	}
}

// MakeUserinfo 给定一个用户信息，创建该用户的缓存信息
func (user *UserCache) MakeUserinfo(ctx context.Context, userID int64) (err error) {
	t := time.Now()
	var eg errgroup.Group

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

	err = eg.Wait()
	fmt.Println("time elapsed:", time.Since(t))
	if err != nil {
		return fmt.Errorf("MakeUserinfo::%w", err)
	}
	return nil
}

// GetUserBasicInfo 获得redis中某个用户的基本信息
func GetUserBasicInfo(ctx context.Context, userID int64) (res map[string]string, err error) {
	res, err = cache.HGetAll(ctx, strconv.FormatInt(userID, 10), cache.UserExpireTime)
	if err != nil {
		return nil, fmt.Errorf("GetUserBasicInfo::%w", err)
	}

	if len(res) == 0 {
		ub := &UserBasic{Userinfo: make(map[string]any)}
		err = ub.makeBasicInfo(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				key := strconv.FormatInt(userID, 10)
				err2 := cache.SetNilHash(ctx, key)
				if err2 != nil {
					return nil, fmt.Errorf("GetUserBasicInfo::RDB.Pipelined: %w", err2)
				}
			}
			return nil, fmt.Errorf("userCache.GetUserBasicInfo::%w", err)
		}

		res, err = cache.HGetAll(ctx, strconv.FormatInt(userID, 10), cache.UserExpireTime)
		if err != nil {
			return nil, fmt.Errorf("userCache.GetUserBasicInfo::%w", err)
		}
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

			ub := &UserBasic{Userinfo: make(map[string]any)}
			err = ub.makeBasicInfo(ctx, userIDs[i])
			if err != nil {
				if errors.Is(errors.Unwrap(err), gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
					key := strconv.FormatInt(userIDs[i], 10)
					err2 := cache.SetNilHash(ctx, key)
					if err2 != nil {
						return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::RDB.Pipelined: %w", err2)
					}
				} else {
					return nil, fmt.Errorf("userCache.GetUserBasicInfo: %w", err)
				}
			}

			//TODO:待优化
			userInfo, err = cache.HGetAll(ctx, strconv.FormatInt(userIDs[i], 10), cache.UserExpireTime)
			if err != nil {
				return nil, fmt.Errorf("userCache.GetUsersBasicInfo: %w", err)
			}
		}
		res = append(res, userInfo)
		fmt.Println("reslen:", len(res))
	}
	return res, nil
}

// MapStringStringToUser 将map[string]string转换为*EntitySets.User结构
func MapStringStringToUser(userBasic map[string]string) (user *EntitySets.User) {
	user = new(EntitySets.User)
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", userBasic["created_at"])
	user.UserID = Utilities.String2Int64(userBasic["user_id"])
	user.UserName = userBasic["user_name"]
	user.Password = userBasic["password"]
	user.Email = userBasic["email"]
	user.Signature = userBasic["signature"]
	user.Shells = Utilities.String2Uint32(userBasic["shells"])
	user.IsAdmin = Utilities.String2Int(userBasic["is_admin"])
	user.CntMsgNotRead = int32(Utilities.String2Int(userBasic["cnt_msg_not_read"]))
	user.CntLikes = Utilities.String2Uint32(userBasic["cnt_likes"])
	user.UserLevel.UserLevel = uint16(Utilities.String2Uint32(userBasic["user_level"]))
	user.AvatarPath = userBasic["user_avatar"]
	return
}
