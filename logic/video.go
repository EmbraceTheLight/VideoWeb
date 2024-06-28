package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/Utilities/logf"
	"VideoWeb/cache/commentCache"
	"VideoWeb/cache/videoCache"
	"VideoWeb/define"
	"VideoWeb/helper"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ParseRange 解析range头的start和end位置，若start或end不存在，则返回对应值为-1
func ParseRange(StartEnd string) (start, end int64) {
	se := strings.Split(StartEnd, "-")
	switch {
	case StartEnd[0] == '-':
		start = -1
		end, _ = strconv.ParseInt(se[1], 10, 64)
	case StartEnd[len(StartEnd)-1] == '-':
		start, _ = strconv.ParseInt(se[0], 10, 64)
		end = -1
	default:
		start, _ = strconv.ParseInt(se[0], 10, 64)
		end, _ = strconv.ParseInt(se[1], 10, 64)
	}
	return
}

// GetVideoDuration 获取视频时长
func GetVideoDuration(VideoPath string) (duration int64, err error) {
	cmd := exec.Command(define.FFProbe, "-i", VideoPath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	outStr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err:", err)
		return 0, err
	}

	tmp, _ := strconv.ParseFloat(strings.TrimRight(string(outStr), "\r\n"), 64) //注意去除字符串末尾的\r\n
	duration = Utilities.RoundOff(tmp)
	return
}

// CreateVideoRecord 创建视频记录
func CreateVideoRecord(tx *gorm.DB, c *gin.Context, UserID int64, videoFilePath string, fileSize int64) (VID int64, err error) {
	VID = GetUUID()
	UID := UserID
	var userInfo *EntitySets.User
	err = tx.Model(&EntitySets.User{}).Where("user_id = ?", UserID).First(&userInfo).Error
	if err != nil {
		return VID, err
	}
	video := &EntitySets.Video{
		MyModel:  define.MyModel{},
		VideoID:  VID,
		UID:      UID,
		UserName: userInfo.UserName,
		Path:     videoFilePath,
		Size:     fileSize,
	}
	err = EntitySets.InsertVideoRecord(tx, video)
	return VID, err
}

// MakeDASHSegments 生成DASH分段文件
func MakeDASHSegments(videoFilePath string) error {
	ext := path.Ext(videoFilePath)
	outputFilePath := path.Dir(videoFilePath) //得到输出文件得到父目录名

	var inputFileName = videoFilePath
	//处理上传的文件不是mp4格式的情况
	if ext != ".mp4" {
		err := helper.Other2MP4(videoFilePath)
		inputFileName = path.Join(outputFilePath, "converted.mp4")
		defer func() { //删除临时生成的视频文件，节约磁盘空间
			err := os.Remove(inputFileName)
			if err != nil {
				logf.WriteErrLog("logic::MakeDASHSegments", fmt.Sprintf("删除%s生成的.mp4临时文件失败:%s", videoFilePath, err.Error()))
			}
		}()
		if err != nil {
			return err
		}
	}

	// 调用ffmpeg命令行工具生成分段文件
	//fmt.Println("inputFilePath:", inputFileName)
	ffmpegArgs := []string{
		"-i", inputFileName,
		"-c", "copy",
		"-f", "dash",
		"-segment_time", "5",
		outputFilePath + "/output.mpd", // 分段文件名模板
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	err := cmd.Run()
	return err
}

// MakePreviews 生成视频预览图
//func MakePreviews(videoFilePath string) error {
//	outputDir := path.Dir(videoFilePath) + "/tmp" //缩略图存放路径
//	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
//		return fmt.Errorf("MakePreviews failed: %w", err)
//	}
//
//	// 使用 FFmpeg 提取关键帧
//	ffmpegArgs := []string{
//		"-i", videoFilePath,
//		"-vf", fmt.Sprintf("select='not(mod(n,%d))'", 20*25),
//		"-q:v", "2",
//		"-vsync", "vfr",
//		filepath.Join(outputDir, "keyframe_%04d.jpg"),
//	}
//	cmd := exec.Command("ffmpeg", ffmpegArgs...)
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	// 运行命令
//	if err := cmd.Run(); err != nil {
//		return fmt.Errorf("MakePreviews failed: %w", err)
//	}
//
//	return nil
//}

// DeleteVideo 删除视频辅助函数
func DeleteVideo(del *EntitySets.Video) error {
	/*从硬盘中删除对应视频信息*/
	err := os.RemoveAll(path.Dir(del.Path))
	if err != nil {
		return err
	}

	/*从数据库中删除视频相关信息*/
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//从数据库中删除视频信息
	err = EntitySets.DeleteVideoInfoByVideoID(tx, del.VideoID)
	if err != nil {
		return err
	}

	//从数据库中删除与视频绑定的Tag信息
	err = EntitySets.DeleteTagRecords(tx, del.VideoID)
	if err != nil {
		return err
	}

	//从数据库中删除与视频绑定的弹幕信息
	err = EntitySets.DeleteBarrageRecordsByVideoID(tx, del.VideoID)
	if err != nil {
		return err
	}

	//从数据库中删除与视频绑定的评论信息
	err = EntitySets.DeleteCommentRecordsByVideoID(tx, del.VideoID)
	err = helper.DeleteCommentWithStatus(del.VideoID, tx)
	if err != nil {
		return err
	}

	//从数据库中删除与视频绑定的所有用户点赞、投币等信息
	err = RelationshipSets.DeleteUserVideoRecordsByVideoID(tx, del.VideoID)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// CheckFileIsExist 检查视频文件是否存在
func CheckFileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

// OpenAndReadFile 打开并读取文件,返回读取到的文件内容
func OpenAndReadFile(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	defer f.Close()
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	return data, err
}

// AddVideoViewCount 增加视频观看次数
func AddVideoViewCount(c *gin.Context, videoID int64) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "AddVideoViewCount")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
	err = helper.UpdateVideoFieldForUpdate(videoID, "cnt_views", 1, tx)
	if err != nil {
		return err
	}
	err = helper.UpdateVideoFieldForUpdate(videoID, "hot", define.AddHotEachView, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// UpdateVideoLikeStatus 更新视频点赞状态,若已经点赞，则为取消点赞，反之则为点赞，更新对应状态与数据
func UpdateVideoLikeStatus(c *gin.Context, UserID, AuthorID, VideoID int64, isLiked bool) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "UpdateVideoLikeStatus")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)

	//更新视频点赞数以及视频热度
	if isLiked {
		err = helper.UpdateVideoFieldForUpdate(VideoID, "likes", -1, tx)
		if err != nil {
			return err
		}

		//更新视频热度
		err = helper.UpdateVideoFieldForUpdate(VideoID, "hot", -define.AddHotEachLike, tx)
		if err != nil {
			return err
		}

		//更新被点赞视频UP主的点赞数
		err = helper.UpdateUserFieldForUpdate(AuthorID, "cnt_likes", -1, tx)
		if err != nil {
			return err
		}

	} else {
		err = helper.UpdateVideoFieldForUpdate(VideoID, "likes", 1, tx)
		if err != nil {
			return err
		}

		err = helper.UpdateVideoFieldForUpdate(VideoID, "hot", define.AddHotEachLike, tx)
		if err != nil {
			return err
		}

		err = helper.UpdateUserFieldForUpdate(AuthorID, "cnt_likes", 1, tx)
		if err != nil {
			return err
		}

	}
	//更新用户点赞状态:当前状态取反
	err = helper.UpdateUserVideoFieldForUpdate(UserID, VideoID, "is_like", !isLiked, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateShells 更新视频投币数,更新对应状态与数据
func UpdateShells(c *gin.Context, videoInfo *EntitySets.Video, TSUID int64, throws int) error {
	var err error
	/*修改贝壳币*/
	//为视频添加贝壳
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "UpdateShells")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
	//为视频添加贝壳
	err = helper.UpdateVideoFieldForUpdate(videoInfo.VideoID, "shells", throws, tx)
	if err != nil {
		return err
	}
	//为作者添加贝壳
	err = helper.UpdateUserFieldForUpdate(videoInfo.UID, "shells", throws, tx)
	if err != nil {
		return err
	}
	//减少投贝壳用户的贝壳数量
	err = helper.UpdateUserFieldForUpdate(TSUID, "shells", -throws, tx)
	if err != nil {
		return err
	}
	//增加视频热度
	err = helper.UpdateVideoFieldForUpdate(videoInfo.VideoID, "hot", throws*define.AddHotEachShell, tx)
	if err != nil {
		return err
	}
	//增加投贝壳者经验
	err = AddExpForThrowShells(c, TSUID, throws, tx)
	if err != nil {
		return err
	}
	//增加作者经验
	err = AddExpForGainShells(c, videoInfo.UID, throws, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateVideoFavorite 更新视频收藏次数相关数据
func UpdateVideoFavorite(c *gin.Context, videoID, fid, uid int64, change int) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "UpdateVideoFavorite")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*收藏视频*/
	//更新用户收藏记录
	if change == 1 { //收藏
		newFV := &RelationshipSets.FavoriteVideo{
			UserID:     uid,
			FavoriteID: fid,
			VideoID:    videoID,
		}
		err = RelationshipSets.InsertFavoriteVideoRecord(tx, newFV)
	} else if change == -1 { //取消收藏
		err = RelationshipSets.DeleteFavoriteVideoRecordByUserIDVideoID(tx, uid, videoID)
	}
	if err != nil {
		return err
	}

	//更新Video收藏数
	println(change)
	err = helper.UpdateVideoFieldForUpdate(videoID, "cnt_favorites", change, tx.Set("gorm:query_option", "FOR UPDATE"))
	if err != nil {
		return err
	}

	//更新UserVideo用户状态
	if change == 1 { //收藏
		err = helper.UpdateUserVideoFieldForUpdate(uid, videoID, "is_favor", true, tx)
	} else if change == -1 { //取消收藏
		err = helper.UpdateUserVideoFieldForUpdate(uid, videoID, "is_favor", false, tx)
	}
	if err != nil {
		return err
	}

	//更新视频热度
	if change == 1 { //收藏
		err = helper.UpdateVideoFieldForUpdate(videoID, "hot", define.AddHotEachFavorite, tx)
	} else if change == -1 { //取消收藏
		err = helper.UpdateVideoFieldForUpdate(videoID, "hot", -define.AddHotEachFavorite, tx)
	}
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// GetVideoListByClass 根据分类及其热度获取视频列表
func GetVideoListByClass(c *gin.Context, class string) (videos []*EntitySets.VideoSummary, err error) {
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "GetVideoListByClass")
		}
	}()

	if class == "recommend" {

		err = DAO.DB.Model(&EntitySets.Video{}).
			Order("hot desc").Limit(define.DefaultSize).Omit("class").Find(&videos).Error
	} else {
		//helper.Get
		err = DAO.DB.Model(&EntitySets.Video{}).Where("class=?", class).
			Order("hot desc").Limit(define.DefaultSize).Find(&videos).Error
	}
	return
}

// GetVideoCommentsList 获取视频评论列表
func GetVideoCommentsList(c *gin.Context, videoID, userID int64, order string, offset, commentsNumbers int) (ret []*EntitySets.CommentSummary, err error) {
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "GetVideoCommentsList")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//找出所有评论
	mapComments, err := videoCache.GetAllVideoCommentsInfo(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if len(mapComments) == 0 {
		return nil, nil
	}
	//if err != nil {
	//	return nil, err
	//}
	//likedCommentIDs := Utilities.Strings2Int64s(tmpLikedCommentIDs)

	//将map转为EntitySets.CommentSummary切片
	var comments []*EntitySets.CommentSummary
	for _, v := range mapComments {
		comments = append(comments, commentCache.MapStringStringToComment(v))
	}

	//对评论按照`to`字段升序排序
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].To < comments[j].To
	})

	//获取根评论
	rootComments := make([]*EntitySets.CommentSummary, 0)

	//TODO:使用二分查找优化该逻辑
	var idx int
	for i, v := range comments {
		if v.To != -1 {
			break
		}
		rootComments = append(rootComments, v)
		idx = i
	}
	var start, end int
	start = offset
	end = start + commentsNumbers
	if end > idx {
		end = idx + 1
	}

	//获取用户点赞的评论
	//tmpLikedCommentIDs, err := cache.SInter(
	//	ctx,
	//	strconv.FormatInt(userID, 10)+strconv.FormatInt(videoID, 10)+"_liked_comments",
	//	strconv.FormatInt(videoID, 10)+"_comments",
	//)

	//递归获取每个根评论的回复列表
	var replies []*EntitySets.CommentSummary
	for _, comment := range comments[start:end] {
		replies = helper.GetCommentReplies(videoID, comment.CommentID, order, comments[idx+1:])
		comment.Replies = replies
		ret = append(ret, comment)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Likes > ret[j].Likes
	})
	//TODO:获取用户对这些评论的点赞/点踩信息

	//
	//获取用户对这些评论的点赞/点踩信息
	//likes, dislikes, err := helper.GetUserCommentRecords(userID, videoID, DAO.DB)
	//if err != nil {
	//	return nil, err
	//}
	//
	//遍历获得的评论，递归更新点赞/点踩信息
	//helper.UpdateCommentsStatus(likes, dislikes, ret)
	return
}

// GetVideosByKey 获取视频列表
func GetVideosByKey(c *gin.Context, key, order string, offset, videoNums int) (videos []*EntitySets.VideoSummary, err error) {
	defer Utilities.DeferFunc(c, err, "GetVideosByKey")
	switch order {
	case "default":
		err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
			Where("MATCH(title,description) AGAINST (? IN BOOLEAN MODE)", key).Order("hot desc").Limit(videoNums).Find(&videos).Error
		//err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
		//	Where("title LIKE ? OR description LIKE ?", "%"+key+"%", "%"+key+"%").Order("hot desc").Limit(videoNums).Find(&videos).Error
	case "mostPlay":
		err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
			Where("MATCH(title,description) AGAINST (? IN BOOLEAN MODE)", key).Order("cnt_views desc").Limit(videoNums).Find(&videos).Error
		//err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
		//	Where("title LIKE ? OR description LIKE ?", "%"+key+"%", "%"+key+"%").Order("cnt_views desc").Limit(videoNums).Find(&videos).Error
	case "newest":
		err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
			Where("MATCH(title,description) AGAINST (? IN BOOLEAN MODE)", key).Order("created_at desc").Limit(videoNums).Find(&videos).Error
		//err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
		//	Where("title LIKE ? OR description LIKE ?", "%"+key+"%", "%"+key+"%").Order("created_at desc").Limit(videoNums).Find(&videos).Error
	case "mostBarrage":
		err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
			Where("MATCH(title,description) AGAINST (? IN BOOLEAN MODE)", key).Order("cnt_barrages desc").Limit(videoNums).Find(&videos).Error
		//err = DAO.DB.Model(&EntitySets.Video{}).Offset(offset).
		//	Where("title LIKE ? OR description LIKE ?", "%"+key+"%", "%"+key+"%").Order("cnt_barrages desc").Limit(videoNums).Find(&videos).Error
	case "mostFavorite":
		err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
			Where("MATCH(title,description) AGAINST (? IN BOOLEAN MODE)", key).Order("cnt_favorites desc").Limit(videoNums).Find(&videos).Error
		//err = DAO.DB.Debug().Model(&EntitySets.Video{}).Offset(offset).
		//	Where("title LIKE ? OR description LIKE ?", "%"+key+"%", "%"+key+"%").Order("cnt_favorites desc").Limit(videoNums).Find(&videos).Error
	}
	return
}
