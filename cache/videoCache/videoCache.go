package videoCache

import (
	"VideoWeb/DAO"
	"VideoWeb/cache"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
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

// GetBarragesInfo Gets the detailed information of a barrage from cache.
func GetBarragesInfo(ctx context.Context, videoID int64) (barrages []map[string]string, err error) {
	barrageIDs, err := cache.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_barrages")
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetBarragesInfo::%w", err)
	}

	barrages, err = cache.GetInfos(ctx, videoID, barrageIDs...)
	if err != nil {
		err = fmt.Errorf("VideoCache.GetBarragesInfos::%w", err)
	}
	return
}

// GetTagsInfo Gets the detailed information of a tag from cache.
func GetTagsInfo(ctx context.Context, videoID int64) (tags []string, err error) {
	tags, err = cache.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_tags")
	if err != nil {
		err = fmt.Errorf("VideoCache.GetTagsInfo::%w", err)
	}
	return
}

// GetVideoCommentsInfo Gets the detailed information of a comment from cache.
func GetVideoCommentsInfo(ctx context.Context, videoID int64) (comments []map[string]string, err error) {
	commentIDs, err := cache.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_comments")
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideoCommentsInfo::%w", err)
	}

	comments, err = cache.GetInfos(ctx, videoID, commentIDs...)
	if err != nil {
		err = fmt.Errorf("VideoCache.GetVideoCommentsInfo::%w", err)
	}
	return
}

// GetVideoBasicInfo gets basic information of a video from cache.
func GetVideoBasicInfo(ctx context.Context, videoID int64) (videoBasic map[string]string, err error) {
	videoBasic, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10))
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
	}

	if len(videoBasic) == 0 { // Video not found in cache, get from database and set cache
		vbasic := &VideoBasic{VideoInfo: make(map[string]any)}
		err = vbasic.makeBasicInfo(ctx, videoID)
		if err != nil {
			if errors.Is(errors.Unwrap(err), gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
				key := strconv.FormatInt(videoID, 10)
				err2 := cache.SetNilHash(ctx, key)
				if err2 != nil {
					return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::RDB.Pipelined:%w", err2)
				}
			}
			return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
		}

		videoBasic, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10))
		if err != nil {
			return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
		}
	}
	return videoBasic, nil
}

// GetVideosBasicInfo gets basic information of many videos from cache.
func GetVideosBasicInfo(ctx context.Context, videoIDs []int64) (res []map[string]string, err error) {
	res = make([]map[string]string, 0)

	cmds := make([]*redis.MapStringStringCmd, len(videoIDs))

	pipe := DAO.RDB.Pipeline()
	for i, videoID := range videoIDs {
		cmds[i] = pipe.HGetAll(ctx, strconv.FormatInt(videoID, 10))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::%w", err)
	}

	for i, cmd := range cmds {
		var videoInfo = make(map[string]string)
		videoInfo, err = cmd.Result()
		if err != nil || len(videoInfo) == 0 {
			if err != nil {
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::cmd.Result():%w", err)
			}

			vbasic := &VideoBasic{VideoInfo: make(map[string]any)}
			err = vbasic.makeBasicInfo(ctx, videoIDs[i])
			if err != nil {
				if errors.Is(errors.Unwrap(err), gorm.ErrRecordNotFound) { // Video not found in database, set cache to not found
					key := strconv.FormatInt(videoIDs[i], 10)
					err2 := cache.SetNilHash(ctx, key)
					if err2 != nil {
						return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::RDB.Pipelined:%w", err2)
					}
				}
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::%w", err)
			}

			videoInfo, err = cache.HGetAll(ctx, strconv.FormatInt(videoIDs[i], 10))
			if err != nil {
				return nil, fmt.Errorf("VideoCache.GetVideosBasicInfo::%w", err)
			}
		}
		res = append(res, videoInfo)
	}
	return res, nil
}
