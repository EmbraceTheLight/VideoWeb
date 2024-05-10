package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
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
func CreateVideoRecord(tx *gorm.DB, c *gin.Context, videoFilePath string, fileSize int64) (VID int64, err error) {
	t, err := GetVideoDuration(videoFilePath)
	if err != nil {
		return VID, err
	}
	UserID := c.Param("ID")
	videoTime, _ := Utilities.SecondToTime(t)
	Title := c.PostForm("title")
	Description := c.PostForm("description")
	Class := c.PostForm("class")

	VID = GetUUID()
	UID := Utilities.String2Int64(UserID)
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
	err = tx.Model(&EntitySets.Video{}).Create(&video).Error
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
	err = EntitySets.DeleteVideoInfoByID(tx, del.VideoID)
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
	//从数据库中删除与视频绑定的收藏信息
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

func UpdateVideoFieldForUpdate(c *gin.Context, VideoID int64, field string, change int) error {
	return helper.UpdateVideoFieldForUpdate(c, VideoID, field, change)
}

func UpdateShells(c *gin.Context, videoInfo *EntitySets.Video, TSUID int64, throws int) error {
	/*修改贝壳币*/
	//为视频添加贝壳
	err := UpdateVideoFieldForUpdate(c, videoInfo.VideoID, "shells", throws)
	if err != nil {
		return err
	}
	//为作者添加贝壳
	err = UpdateUserFieldForUpdate(c, videoInfo.UID, "shells", throws)
	if err != nil {
		return err
	}
	//减少投贝壳用户的贝壳数量
	err = UpdateUserFieldForUpdate(c, TSUID, "shells", -throws)
	if err != nil {
		return err
	}
	return nil
}
