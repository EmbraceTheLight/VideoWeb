package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/Utilities/logf"
	"VideoWeb/define"
	"VideoWeb/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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
	t, err := GetVideoDuration(videoFilePath)
	if err != nil {
		return VID, err
	}
	videoTime, _ := Utilities.SecondToTime(t)
	Title := c.PostForm("title")
	Description := c.PostForm("description")
	Class := c.PostForm("class")

	VID = GetUUID()
	UID := UserID
	video := &EntitySets.Video{
		MyModel:     define.MyModel{},
		VideoID:     VID,
		UID:         UID,
		Title:       Title,
		Description: Description,
		Class:       Class,
		Path:        videoFilePath,
		Duration:    videoTime,
		Size:        fileSize,
	}
	err = EntitySets.InsertVideoRecord(tx, video)
	return VID, err
}

func MakeDASHSegments(videoFilePath string) error {
	ext := path.Ext(videoFilePath)
	outputFilePath := path.Dir(videoFilePath) //得到输出文件得到父目录名

	var inputFileName = videoFilePath
	//处理上传的文件不是mp4格式的情况
	if ext != ".mp4" {
		err := helper.Other2MP4(videoFilePath)
		inputFileName = path.Join(outputFilePath, "converted.mp4")
		defer func() {
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
	fmt.Println("inputFilePath:", inputFileName)
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
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//从数据库中删除视频信息
	err = EntitySets.DeleteVideoInfoByVideoID(tx, del.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//从数据库中删除与视频绑定的Tag信息
	err = EntitySets.DeleteTagRecords(tx, del.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//从数据库中删除与视频绑定的弹幕信息
	err = EntitySets.DeleteBarrageRecordsByVideoID(tx, del.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//从数据库中删除与视频绑定的评论信息
	err = EntitySets.DeleteCommentRecordsByVideoID(tx, del.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}
	//TODO:从数据库中删除与视频绑定的收藏信息

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
	Utilities.AddFuncName(c, "AddVideoViewCount")
	err := helper.UpdateVideoFieldForUpdate(videoID, "cnt_views", 1, nil)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return err
	}
	return nil
}

// UpdateVideoLikeStatus 更新视频点赞状态,若已经点赞，则为取消点赞，反之则为点赞，更新对应状态与数据
func UpdateVideoLikeStatus(c *gin.Context, UserID, VideoID int64, field string, isLiked bool) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "UpdateVideoLikeStatus")
			Utilities.HandleInternalServerError(c, err)
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)

	//更新视频点赞数
	if isLiked {
		err = helper.UpdateVideoFieldForUpdate(VideoID, field, -1, tx)
		if err != nil {
			return err
		}
	} else {
		err = helper.UpdateVideoFieldForUpdate(VideoID, field, 1, tx)
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
			Utilities.HandleInternalServerError(c, err)
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
	//为视频添加贝壳
	fmt.Println("tx before update video shells:", tx)
	err = helper.UpdateVideoFieldForUpdate(videoInfo.VideoID, "shells", throws, tx)
	if err != nil {
		return err
	}
	//为作者添加贝壳
	fmt.Println("tx before update Author shells:", tx)
	err = helper.UpdateUserFieldForUpdate(videoInfo.UID, "shells", throws, tx)
	if err != nil {
		return err
	}
	//减少投贝壳用户的贝壳数量
	fmt.Println("tx before update user shells:", tx)
	err = helper.UpdateUserFieldForUpdate(TSUID, "shells", -throws, tx)
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
			Utilities.HandleInternalServerError(c, err)
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*收藏视频*/
	//增加用户收藏记录
	newFV := &RelationshipSets.FavoriteVideo{
		UserID:     uid,
		FavoriteID: fid,
		VideoID:    videoID,
	}
	err = RelationshipSets.InsertFavoriteVideoRecord(tx, newFV)
	if err != nil {
		return err
	}
	//更新Video收藏数
	err = helper.UpdateVideoFieldForUpdate(videoID, "cnt_favorites", change, tx)
	if err != nil {
		return err
	}
	//更新UserVideo用户状态
	err = helper.UpdateUserVideoFieldForUpdate(uid, videoID, "is_favor", true, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
