package videoCache

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities/logf"
	"VideoWeb/cache"
	"VideoWeb/cache/commentCache"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func (vbasic *VideoBasic) makeBasicInfo(ctx context.Context, videoID int64) error {
	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, videoID)
	if err != nil {
		return fmt.Errorf("VideoBasic->makeBasicInfo::EntitySets.GetVideoInfoByID error: %w", err)
	}
	vbasic.VideoInfo["created_at"] = videoInfo.CreatedAt.Format("2006-01-02 15:04:05")
	vbasic.VideoInfo["title"] = videoInfo.Title
	vbasic.VideoInfo["description"] = videoInfo.Description
	vbasic.VideoInfo["class"] = videoInfo.Class
	vbasic.VideoInfo["path"] = videoInfo.Path
	vbasic.VideoInfo["video_id"] = videoInfo.VideoID
	vbasic.VideoInfo["user_id"] = videoInfo.UID
	vbasic.VideoInfo["user_name"] = videoInfo.UserName
	vbasic.VideoInfo["likes"] = videoInfo.Likes
	vbasic.VideoInfo["shells"] = videoInfo.Shells
	vbasic.VideoInfo["hot"] = videoInfo.Hot
	vbasic.VideoInfo["cnt_barrages"] = videoInfo.CntBarrages
	vbasic.VideoInfo["cnt_shares"] = videoInfo.CntShares
	vbasic.VideoInfo["cnt_favorites"] = videoInfo.CntFavorites
	vbasic.VideoInfo["cnt_views"] = videoInfo.CntViews
	vbasic.VideoInfo["duration"] = videoInfo.Duration
	vbasic.VideoInfo["size"] = videoInfo.Size
	vbasic.VideoInfo["cover_path"] = videoInfo.CoverPath

	err = cache.HSetWithRetry(
		ctx,
		strconv.FormatInt(videoInfo.VideoID, 10),
		cache.DefaultTry, cache.DefaultSleep, cache.VideoExpireTime,
		vbasic.VideoInfo,
	)
	if err != nil {
		return fmt.Errorf("VideoBasic->makeBasicInfo::cache.HSetWithRetry error: %w", err)
	}
	return nil
}

func (vb *VideoBarrages) makeBarrageInfo(ctx context.Context, videoID int64) error {
	var barrages []*EntitySets.Barrage
	err := DAO.DB.Model(&EntitySets.Barrage{}).
		Where("video_id = ?", videoID).
		Find(&barrages).Error
	if err != nil {
		return fmt.Errorf("VideoBarrages->makeBarrageInfo: %w", err)
	}
	var barrageIDs = make([]any, len(barrages))
	for i, b := range barrages {
		barrageIDs[i] = b.BID
	}

	if len(barrageIDs) != 0 {
		err = cache.SAddWithRetry(
			ctx,
			strconv.FormatInt(videoID, 10)+"_barrages",
			cache.DefaultTry, cache.DefaultSleep, cache.VideoExpireTime,
			barrageIDs...,
		)
		if err != nil {
			return fmt.Errorf("VideoBarrages->MakeBarrageInfo::cache.SAddWithRetry: %w", err)
		}

		err = makeBarragesInfos(ctx, videoID, barrages...)
		if err != nil {
			return fmt.Errorf("VideoBarrages->MakeBarrageInfo::makeBarragesInfos: %w", err)
		}

		//TODO:此处应添加“将barrage个数缓存入VideoBasic中”的逻辑，而在表EntitySets.Video中应该去掉CntBarrages字段。在完成缓存系统后进行修改。
	}
	return nil
}

func (vt *VideoTags) makeTagInfo(ctx context.Context, videoID int64) error {
	var tags []any
	err := DAO.DB.Model(&EntitySets.Tags{}).
		Where("video_id = ?", videoID).
		Select("tag").
		Find(&tags).Error
	if err != nil {
		return fmt.Errorf("VideoTag->makeTagInfo: %w", err)
	}

	if len(tags) != 0 {
		err = cache.SAddWithRetry(
			ctx,
			strconv.FormatInt(videoID, 10)+"_tags",
			cache.DefaultTry,
			cache.DefaultSleep,
			cache.VideoExpireTime,
			tags...,
		)
		if err != nil {
			return fmt.Errorf("VideoTag->makeTagInfo::cache.SAddWithRetry: %w", err)
		}
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

// getSpecificVideoCommentsInfo Gets the detailed information of specific comments of a video from cache.
func getSpecificVideoCommentsInfo(ctx context.Context, videoID int64, commentID ...string) (comments []map[string]string, err error) {
	pipe := DAO.RDB.Pipeline()
	prefix := strconv.FormatInt(videoID, 10)
	cmds := make([]*redis.MapStringStringCmd, len(commentID))
	for i, id := range commentID {
		cmds[i] = pipe.HGetAll(ctx, prefix+id)
		pipe.Expire(ctx, prefix+id, cache.CommentExpireTime)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("VideoCache.getSpecificVideoCommentsInfo::pipe.Exec(): %w", err)
	}
	for _, cmd := range cmds {
		comments = append(comments, cmd.Val())
	}
	return
}

func (vc *VideoComments) makeCommentsInfo(ctx context.Context, videoID int64) error {
	var comments []*EntitySets.Comments
	err := DAO.DB.Model(&EntitySets.Comments{}).
		Where("video_id = ?", videoID).
		Find(&comments).Error

	if err != nil {
		return fmt.Errorf("VideoComment->makeCommentsInfo: %w", err)
	}

	if len(comments) != 0 {
		commentIDs := make([]string, len(comments))
		for i, c := range comments {
			commentIDs[i] = strconv.FormatInt(c.CommentID, 10)
		}
		err = commentCache.MakeCommentInfos(ctx, videoID, comments...)

		err = cache.SAddWithRetry(
			ctx,
			strconv.FormatInt(videoID, 10)+"_comments",
			cache.DefaultTry, cache.DefaultSleep,
			cache.VideoExpireTime,
			commentIDs,
		)
		if err != nil {
			return fmt.Errorf("VideoComment->makeCommentsInfo::cache.SAddWithRetry: %w", err)
		}
	}
	return nil
}
func set(set []string, v any) {
	set = append(set, v.(string))
}
func getOne(ctx context.Context, funcName, key, member string) string {
	err := cache.SIsMember(ctx, key, member)
	if err != nil {
		logf.WriteErrLog(funcName, err.Error())
		return ""
	}
	return member
}
func del(sli []string, v any) {
	for i, val := range sli {
		if val == v.(string) {
			sli = append(sli[:i], sli[i+1:]...)
			break
		}
	}
}
