package userCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities/logf"
	"VideoWeb/cache"
	"VideoWeb/cache/commentCache"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

func set(set []int64, v any) {
	set = append(set, v.(int64))
}
func getOne(ctx context.Context, funcName, key, member string) string {
	err := cache.SIsMember(ctx, key, member)
	if err != nil {
		logf.WriteErrLog(funcName, err.Error())
		return ""
	}
	return member
}
func del(sli []int64, v any) {
	for i, val := range sli {
		if val == v.(int64) {
			sli = append(sli[:i], sli[i+1:]...)
			break
		}
	}
}

func (ub *UserBasic) makeBasicInfo(ctx context.Context, userID int64) error {
	userinfo, err := EntitySets.GetUserInfoByID(DAO.DB, userID)
	if err != nil {
		return fmt.Errorf("UserBasic->makeBasicInfo::EntitySets.GetUserInfoByID: %w", err)
	}

	level, err := EntitySets.GetLevelRecordByUserID(DAO.DB, userID)
	if err != nil {
		return fmt.Errorf("UserBasic->makeBasicInfo::EntitySets.GetLevelRecordByUserID: %w", err)
	}

	// user 基本信息
	ub.Userinfo["created_at"] = userinfo.CreatedAt.Format("2006-01-02 15:04:05")
	ub.Userinfo["user_id"] = userinfo.UserID
	ub.Userinfo["user_name"] = userinfo.UserName
	ub.Userinfo["password"] = userinfo.Password
	ub.Userinfo["email"] = userinfo.Email
	ub.Userinfo["signature"] = userinfo.Signature
	ub.Userinfo["shells"] = userinfo.Shells
	ub.Userinfo["is_admin"] = userinfo.IsAdmin
	ub.Userinfo["cnt_msg_not_read"] = userinfo.CntMsgNotRead
	ub.Userinfo["cnt_likes"] = userinfo.CntLikes
	ub.Userinfo["level"] = level.UserLevel
	ub.Userinfo["avatar"] = userinfo.AvatarPath

	err = cache.HSetWithRetry(
		ctx,
		strconv.FormatInt(userinfo.UserID, 10),
		cache.DefaultTry, cache.DefaultSleep, cache.UserExpireTime,
		ub.Userinfo,
	)
	if err != nil {
		return fmt.Errorf("UserBasic->makeBasicInfo::cache.HSetWithRetry: %w", err)
	}
	return nil
}

func (ufs *UserFollows) makeFollowsInfo(ctx context.Context, userID int64) error {
	var followsIDs []any
	err := DAO.DB.Model(&RelationshipSets.UserFollows{}).
		Where("user_id = ?", userID).
		Select("follow_user_id").
		Find(&followsIDs).Error
	if err != nil {
		return fmt.Errorf("UserFollows->makeFollowsInfo: %w", err)
	}

	if len(followsIDs) != 0 {
		err = cache.SAddWithRetry(
			ctx,
			strconv.FormatInt(userID, 10)+cache.FlsSfx,
			cache.DefaultTry,
			cache.DefaultSleep,
			cache.UserExpireTime,
			followsIDs...,
		)
		if err != nil {
			return fmt.Errorf("UserFollows->makeFollowsInfo: %w", err)
		}

		err = addToUserCache(ctx, strconv.FormatInt(userID, 10), "cnt_follows", len(followsIDs))
		if err != nil {
			return fmt.Errorf("UserFollows->makeFollowsInfo: %w", err)
		}
	}
	return nil
}

func (ufd *UserFollowed) makeFollowedInfo(ctx context.Context, userID int64) error {
	var followedIDs []any
	err := DAO.DB.Model(&RelationshipSets.UserFollowed{}).
		Where("user_id = ?", userID).
		Select("followed_id").
		Find(&followedIDs).Error
	if err != nil {
		return fmt.Errorf("UserFollowed->makeFollowedInfo: %w", err)
	}

	if len(followedIDs) != 0 { //if followedIDs is not empty, add it to redis cache
		ufd.key = strconv.FormatInt(userID, 10) + cache.FldSfx
		err = cache.SAddWithRetry(ctx, ufd.key, cache.DefaultTry, cache.DefaultSleep, cache.UserExpireTime, followedIDs...)
		if err != nil {
			return fmt.Errorf("UserFollowed->makeFollowedInfo::%w", err)
		}

		err = addToUserCache(ctx, strconv.FormatInt(userID, 10), "cnt_followed", len(followedIDs)) //add or update a key-value of followed count to user cache
		if err != nil {
			return fmt.Errorf("UserFollowed->makeFollowedInfo::%w", err)
		}
	}

	return nil
}

func (uc *UserComments) makeCommentsInfo(ctx context.Context, userID int64) error {
	var comments []*EntitySets.Comments
	comments, err := EntitySets.GetCommentsByUserID(DAO.DB, userID)
	if err != nil {
		return fmt.Errorf("UserComments->makeCommentsInfo: %w", err)
	}

	if len(comments) != 0 {
		commentIDs := make([]any, len(comments))
		for i, c := range comments {
			commentIDs[i] = c.CommentID
		}

		err = commentCache.MakeCommentInfos(ctx, comments...)
		if err != nil {
			return fmt.Errorf("UserComments->makeCommentsInfo: %w", err)
		}

		uc.key = strconv.FormatInt(userID, 10) + cache.CommentSfx
		err = cache.SAddWithRetry(ctx, uc.key, cache.DefaultTry, cache.DefaultSleep, cache.CommentExpireTime, commentIDs...)
		if err != nil {
			return fmt.Errorf("UserVideo->makeVideoIDInfo: %w", err)
		}
	}
	return nil
}

func (uv *UserVideo) makeVideoIDInfo(ctx context.Context, userID int64) error {
	var videoIDs []any
	err := DAO.DB.Model(&EntitySets.Video{}).
		Where("user_id = ?", userID).
		Select("video_id").
		Find(&videoIDs).Error
	if err != nil {
		return fmt.Errorf("UserVideo->makeVideoIDInfo: %w", err)
	}

	uv.key = strconv.FormatInt(userID, 10) + "_videos"
	err = cache.SAddWithRetry(ctx, uv.key, cache.DefaultTry, cache.DefaultSleep, cache.VideoExpireTime, videoIDs...)
	if err != nil {
		return fmt.Errorf("UserVideo->makeVideoIDInfo: %w", err)
	}

	return nil
}

func (us *UserSearch) makeUserSearch(ctx context.Context, userID int64) error {
	var searchStrings []*EntitySets.UserSearchHistory
	err := DAO.DB.Model(&EntitySets.UserSearchHistory{}).
		Where("user_id = ?", userID).
		Find(&searchStrings).Error
	if err != nil {
		return fmt.Errorf("UserSearch->makeUserSearch::%w", err)
	}
	if len(searchStrings) == 0 {
		return nil
	}

	var ZSetValues = make([]cache.ZSetValue, len(searchStrings))
	for i, s := range searchStrings {
		ZSetValues[i] = s
	}

	err = cache.AddToZSet(ctx, strconv.FormatInt(userID, 10)+cache.SearchSfx, cache.UserExpireTime, ZSetValues...)
	if err != nil {
		return fmt.Errorf("UserSearch->makeUserSearch::%w", err)
	}

	return nil
}

func (uw *UserWatch) makeUserWatch(ctx context.Context, userID int64) error {
	var watches []*EntitySets.UserWatch
	err := DAO.DB.Model(&EntitySets.UserWatch{}).
		Where("user_id = ?", userID).
		Find(&watches).Error
	if err != nil {
		return fmt.Errorf("UserWatch->makeUserWatch::%w", err)
	}
	if len(watches) == 0 {
		return nil
	}

	var ZSetValues = make([]cache.ZSetValue, len(watches))
	for i, w := range watches {
		ZSetValues[i] = w
	}
	err = cache.AddToZSet(ctx, strconv.FormatInt(userID, 10)+cache.WatchSfx, cache.UserExpireTime, ZSetValues...)
	if err != nil {
		return fmt.Errorf("UserWatch->makeUserWatch::%w", err)
	}

	return nil
}

// addToUserCache sets hash table in redis
func addToUserCache(ctx context.Context, key string, values ...any) (err error) {
	mp := cache.MakeMap(values...)
	if mp == nil {
		return fmt.Errorf("UserCache->addToUserCache: arguments is invalid")
	}

	err = cache.HSetWithRetry(
		ctx, key,
		cache.DefaultTry, cache.DefaultSleep, cache.UserExpireTime,
		mp,
	)
	if err != nil {
		return fmt.Errorf("UserCache->addToUserCache: %w", err)
	}
	return nil
}

// checkUserCache checks if the user cache is exist, if not, it will create it
func checkUserCache(ctx context.Context, key string, userID int64) (err error) {
	if DAO.RDB.TTL(ctx, key).Val() < 0 {
		userCache := MakeUserCache()
		err = userCache.MakeUserinfo(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err2 := cache.SetNilHash(ctx, key)
				if err2 != nil {
					return fmt.Errorf("checkUserCache::%w", err)
				}
			} else {
				return fmt.Errorf("checkUserCache::%w", err)
			}
		}
	}
	return nil
}
