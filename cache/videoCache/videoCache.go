package videoCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/Utilities/logf"
	"VideoWeb/cache"
	"VideoWeb/cache/userCache"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// MakeVideoCache creates a VideoCache object that contains all the necessary information of a video.
func MakeVideoCache() *VideoCache {
	vc := &VideoCache{}
	vc.initVideoCache()
	return vc
}

// MakeVideoInfo creates a VideoInfo object that contains basic information,barrages,and tags.
func (video *VideoCache) MakeVideoInfo(ctx context.Context, videoID int64) (err error) {
	var eg errgroup.Group

	// Make basic information
	eg.Go(func() error { return video.VBasic.makeBasicInfo(ctx, videoID) })

	// Make barrage information
	eg.Go(func() error { return video.VBarrages.makeBarrageInfo(ctx, videoID) })

	// Make tag information
	eg.Go(func() error { return video.VTags.makeTagInfo(ctx, videoID) })

	// Make comment information
	eg.Go(func() error { return video.VComments.makeCommentsInfo(ctx, videoID) })

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("VideoCache.MakeVideoInfo::%w", err)
	}

	return nil
}

// DeleteVideoCache deletes all the information of a video from cache.
func DeleteVideoCache(ctx context.Context, videoID int64) (err error) {
	key := strconv.FormatInt(videoID, 10)
	var eg errgroup.Group

	// Delete basic information
	eg.Go(func() error { return DAO.RDB.Del(ctx, key).Err() })

	//Delete barrage information
	eg.Go(func() error { return DAO.RDB.Del(ctx, key+cache.BarrageSfx).Err() })

	//Delete tag information
	eg.Go(func() error { return DAO.RDB.Del(ctx, key+cache.TagSfx).Err() })

	//Delete comment information
	eg.Go(func() error { return DAO.RDB.Del(ctx, key+cache.CommentSfx).Err() })

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("VideoCache.DeleteVideoCache::%w", err)
	}
	return nil
}

// DelayDoubleDelete delete the information about a video with delay
func DelayDoubleDelete(ctx context.Context, videoID int64) {
	time.Sleep(cache.DefaultSleep)
	err := DeleteVideoCache(ctx, videoID)
	if err != nil {
		logf.WriteErrLog("videoCache.DeleteVideoCache", err.Error())
	}
}

// GetBarragesInfo Gets the detailed information of a barrage from cache.
func GetBarragesInfo(ctx context.Context, videoID int64) (barrages []map[string]string, err error) {
	key := strconv.FormatInt(videoID, 10) + cache.BarrageSfx

	if checkVideoCache(ctx, key, videoID) != nil {
		return nil, fmt.Errorf("VideoCache.GetBarragesInfo::%w", err)
	}

	barrageIDs, err := cache.SMembers(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetBarragesInfo::%w", err)
	}

	barrages, err = cache.GetInfos(ctx, barrageIDs...)
	if err != nil {
		err = fmt.Errorf("VideoCache.GetBarragesInfos::%w", err)
	}
	return
}

// GetTagsInfo Gets the detailed information of a tag from cache.
func GetTagsInfo(ctx context.Context, videoID int64) (tags []string, err error) {
	key := strconv.FormatInt(videoID, 10) + cache.TagSfx

	if checkVideoCache(ctx, key, videoID) != nil {
		return nil, fmt.Errorf("VideoCache.GetTagsInfo::%w", err)
	}

	tags, err = cache.SMembers(ctx, key)
	if err != nil {
		err = fmt.Errorf("VideoCache.GetTagsInfo::%w", err)
	}

	return
}

// GetAllVideoCommentsInfo Gets the detailed information of comments of a video from cache.
func GetAllVideoCommentsInfo(ctx context.Context, videoID int64) (comments []map[string]string, err error) {
	key := strconv.FormatInt(videoID, 10) + cache.CommentSfx

	if checkVideoCache(ctx, key, videoID) != nil {
		return nil, fmt.Errorf("VideoCache.GetAllVideoCommentsInfo::%w", err)
	}

	commentIDs, err := cache.SMembers(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetAllVideoCommentsInfo::%w", err)
	}

	comments, err = cache.GetInfos(ctx, commentIDs...)
	if err != nil {
		err = fmt.Errorf("VideoCache.GetAllVideoCommentsInfo::%w", err)
	}
	return
}

// GetUserLikedCommentsInfo Gets the detailed information of comments liked by a user in a video from cache.
func GetUserLikedCommentsInfo(ctx context.Context, videoID int64, userID int64) (res []map[string]string, err error) {
	keyVideo := strconv.FormatInt(videoID, 10) + cache.CommentSfx

	if checkVideoCache(ctx, keyVideo, videoID) != nil {
		return nil, fmt.Errorf("VideoCache.GetUserLikedCommentsInfo::%w", err)
	}

	keyUser := strconv.FormatInt(userID, 10) + strconv.FormatInt(videoID, 10) + cache.ULCSfx
	if DAO.RDB.TTL(ctx, keyUser).Val() < 0 {
		err = userCache.MakeUserLikedComments(ctx, videoID, userID)
		if err != nil {
			return nil, fmt.Errorf("VideoCache.MakeVideoInfo::%w", err)
		}
	}

	likedCommentsIDs, err := cache.SInter(
		ctx,
		keyUser,
		keyVideo,
	)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetUserLikedCommentsInfo::%w", err)
	}

	res, err = getSpecificVideoCommentsInfo(ctx, likedCommentsIDs...)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetUserLikedCommentsInfo::%w", err)
	}
	return
}

// GetVideoBasicInfo gets basic information of a video from cache.
func GetVideoBasicInfo(ctx context.Context, videoID int64) (videoBasic map[string]string, err error) {
	videoBasic, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10), cache.VideoExpireTime)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
	}

	if len(videoBasic) == 0 { // Video not found in cache, get from database and set cache
		vc := &VideoCache{}
		vc.initVideoCache()
		err = vc.MakeVideoInfo(ctx, videoID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
				key := strconv.FormatInt(videoID, 10)
				err2 := cache.SetNilHash(ctx, key)
				if err2 != nil {
					return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err2)
				}
			}
			return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
		}

		videoBasic, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10), cache.VideoExpireTime)
		if err != nil {
			return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
		}
	}
	return videoBasic, nil
}

// GetVideosBasicInfo gets basic information of many videos from cache.
func GetVideosBasicInfo(ctx context.Context, videoIDs ...int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)

	cmds := make([]*redis.MapStringStringCmd, len(videoIDs))

	pipe := DAO.RDB.Pipeline()
	for i, videoID := range videoIDs {
		cmds[i] = pipe.HGetAll(ctx, strconv.FormatInt(videoID, 10))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo: %w", err)
	}

	for i, cmd := range cmds {
		var videoInfo = make(map[string]string)
		videoInfo, err = cmd.Result()
		if err != nil || len(videoInfo) == 0 {
			if err != nil {
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::cmd.Result(): %w", err)
			}

			vinfo := &VideoCache{}
			vinfo.initVideoCache()
			key := strconv.FormatInt(videoIDs[i], 10)

			err = vinfo.MakeVideoInfo(ctx, videoIDs[i])
			if err != nil {
				if errors.Is(errors.Unwrap(err), gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
					err2 := cache.SetNilHash(ctx, key)
					if err2 != nil {
						return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::RDB.Pipelined:%w", err2)
					}
				}
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::%w", err)
			}

			videoInfo, err = cache.HGetAll(ctx, key, cache.VideoExpireTime)
			if err != nil {
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::%w", err)
			}
		}
		res = append(res, videoInfo)
	}
	return res, nil
}

// MapStringString2Videos maps a slice of map[string]string to a slice of *EntitySets.Video.
func MapStringString2Videos(videoInfos ...map[string]string) (res []*EntitySets.Video) {
	for _, videoInfo := range videoInfos {
		var video = new(EntitySets.Video)
		video.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", videoInfo["created_at"])
		video.Title = videoInfo["title"]
		video.Description = videoInfo["description"]
		video.Class = videoInfo["class"]
		video.Path = videoInfo["path"]
		video.VideoID = Utilities.String2Int64(videoInfo["video_id"])
		video.UID = Utilities.String2Int64(videoInfo["user_id"])
		video.UserName = videoInfo["user_name"]
		video.Likes = Utilities.String2Uint32(videoInfo["likes"])
		video.Shells = Utilities.String2Uint32(videoInfo["shells"])
		video.Hot = Utilities.String2Uint32(videoInfo["hot"])
		video.CntBarrages = Utilities.String2Uint32(videoInfo["cnt_barrages"])
		video.CntShares = Utilities.String2Uint32(videoInfo["cnt_shares"])
		video.CntFavorites = Utilities.String2Uint32(videoInfo["cnt_favorites"])
		video.CntViews = Utilities.String2Uint32(videoInfo["cnt_views"])
		video.Duration = videoInfo["duration"]
		video.Size = Utilities.String2Int64(videoInfo["size"])
		video.CoverPath = videoInfo["cover_path"]
		res = append(res, video)
	}

	return
}

// MakeAllVideosZSet makes a sorted set of all videos.
func MakeAllVideosZSet(ctx context.Context) (videoZSetInfos []*VideoZSetInfo, err error) {
	videoZSetInfos = make([]*VideoZSetInfo, 0)
	err = DAO.DB.Model(&EntitySets.Video{}).Find(&videoZSetInfos).Error
	if err != nil {
		return nil, fmt.Errorf("VideoCache.MakeAllVideosZSet: %w", err)
	}

	err = SaveVideosZSet(ctx, cache.VideoZSetKey, videoZSetInfos...)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.MakeAllVideosZSet::SaveVideosZSet: %w", err)
	}
	return
}

// SaveVideosZSet saves the information of video(s) to a sorted set by videos' hot.
func SaveVideosZSet(ctx context.Context, key string, infos ...*VideoZSetInfo) (err error) {
	pipe := DAO.RDB.Pipeline()
	for _, v := range infos {
		pipe.ZAdd(ctx, key, redis.Z{Score: float64(v.Hot), Member: v.VideoID})
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("VideoCache.SaveVideosZSet::RDB.Pipelined:%w", err)
	}
	return nil
}

// GetVideoZSetInfo gets the information of a video from a sorted set.
func GetVideoZSetInfo(ctx context.Context, key string, start, end int64) (videoIDs []int64, err error) {
	cmd := DAO.RDB.ZRevRangeWithScores(ctx, key, start, end)
	tmp, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideoZSetInfo::%w", err)
	}

	videoIDs = make([]int64, len(tmp))
	for i, v := range tmp {
		videoIDs[i] = Utilities.String2Int64(v.Member.(string))
	}
	return

}

// UpdateVideoInfoFields updates the information of a video in cache.
func UpdateVideoInfoFields(ctx context.Context, videoID int64, fields map[string]any) (err error) {
	key := strconv.FormatInt(videoID, 10)
	if len(fields) == 0 {
		return nil
	}

	pipe := DAO.RDB.Pipeline()
	for k, v := range fields {
		pipe.HSet(ctx, key, k, v)
	}
	pipe.Expire(ctx, key, cache.VideoExpireTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("VideoCache.UpdateVideoInfoFields: %w", err)
	}

	return nil
}
