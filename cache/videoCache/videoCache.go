package videoCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
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

// makeBarragesInfos creates Barrage objects of a video.
func makeBarragesInfos(ctx context.Context, prefix int64, barrages ...*EntitySets.Barrage) error {
	barrageInfos := make([]*BarrageInfo, len(barrages))
	HashMap := make([]cache.HashMap, len(barrages))
	for i, b := range barrages {
		barrageInfos[i] = new(BarrageInfo)
		barrageInfos[i].barrageInfo = make(map[string]any)
		barrageInfos[i].key = strconv.FormatInt(prefix, 10) + strconv.FormatInt(b.BID, 10)
		barrageInfos[i].barrageInfo["barrage_id"] = b.BID
		barrageInfos[i].barrageInfo["user_id"] = b.UID
		barrageInfos[i].barrageInfo["video_id"] = b.VID
		barrageInfos[i].barrageInfo["content"] = b.Content
		barrageInfos[i].barrageInfo["color"] = b.Color
		HashMap[i] = barrageInfos[i]
	}
	err := cache.HSets(ctx, cache.VideoExpireTime, HashMap...)
	if err != nil {
		return fmt.Errorf("VideoCache.MakeBarragesInfos::cache.HSets: %w", err)
	}

	return nil

}

// GetVideoInfo Gets all the information of a video from cache.
func GetVideoInfo(ctx context.Context, videoID int64) (basic map[string]string, barrages, Tags, Comments []map[string]string, err error) {
	pipe := DAO.RDB.Pipeline()
	cmdMMS := new(redis.MapStringStringCmd)   //map-string-string command to Get basic information
	cmdSS := make([]*redis.StringSliceCmd, 3) //string-slice commands to get tags,barrages,comments ID

	// Get basic information
	cmdMMS = pipe.HGetAll(ctx, strconv.FormatInt(videoID, 10))
	// Get barrage IDs
	cmdSS[0] = pipe.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_barrages")
	// Get tag IDs
	cmdSS[1] = pipe.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_tags")
	// Get Comment IDs
	cmdSS[2] = pipe.SMembers(ctx, strconv.FormatInt(videoID, 10)+"_comments")

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::Pipelener.Exec: %w", err)
	}

	// Get basic information
	basic, err = cmdMMS.Result()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::cmdMMS.Result(): %w", err)
	}

	// Get barrages' information
	barrageIDs, err := cmdSS[0].Result()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::cmdSS[0].Result(): %w", err)
	}
	barrages, err = getBarragesInfo(ctx, videoID, barrageIDs...)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::getBarragesInfo::%w", err)
	}

	// Get tags' information
	tagIDs, err := cmdSS[1].Result()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::cmdSS[0].Result(): %w", err)
	}
	tags, err := getTagsInfo(ctx, videoID, tagIDs...)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::getTagsInfo::%w", err)
	}

	//Get comments' information
	commentIDs, err := cmdSS[2].Result()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::cmdSS[0].Result(): %w", err)
	}
	comments, err := getCommentsInfo(ctx, videoID, commentIDs...)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VideoCache.GetVideoInfo::getCommentsInfo::%w", err)
	}

	return
}

// GetVideoBasicInfo gets basic information of a video from cache.
func GetVideoBasicInfo(ctx context.Context, videoID int64) (res map[string]string, err error) {
	res, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10))
	if err != nil {
		return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
	}

	if len(res) == 0 { // Video not found in cache, get from database and set cache
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

		res, err = cache.HGetAll(ctx, strconv.FormatInt(videoID, 10))
		if err != nil {
			return nil, fmt.Errorf("VideoCache.GetVideoBasicInfo::%w", err)
		}
	}
	return res, nil
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
