package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// ParseTime 解析时间字符串，返回分钟和秒
func ParseTime(time string) (min, sec int) {
	se := strings.Split(time, ":")
	min, _ = strconv.Atoi(se[0])
	sec, _ = strconv.Atoi(se[1])
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
func CreateVideoRecord(tx *gorm.DB, c *gin.Context, videoFilePath string, fileSize int64) (VID string, err error) {
	t, err := GetVideoDuration(videoFilePath)
	if err != nil {
		return VID, err
	}
	UserID := c.Query("userID")
	videoTime := Utilities.SecondToTime(t)
	Title := c.PostForm("title")
	Description := c.PostForm("description")
	Class := c.PostForm("class")
	VID = GetUUID()
	video := &EntitySets.Video{
		MyModel:     define.MyModel{},
		VideoID:     VID,
		UID:         UserID,
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
	if ext != ".mp4" {
		err := helper.Other2MP4(videoFilePath)
		if err != nil {
			return err
		}
	}
	outPutFilePath := path.Dir(videoFilePath)
	fmt.Println("outPutFilePath:", outPutFilePath)
	ffmpegArgs := []string{
		"-i", outPutFilePath + "/converted.mp4",
		"-c", "copy",
		"-f", "dash",
		"-segment_time", "5",
		outPutFilePath + "/output.mpd", // 分段文件名模板
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

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*从数据库中删除视频信息*/
	err = tx.Where("VideoID=?", del.VideoID).Delete(&EntitySets.Video{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	/*从数据库中删除与视频绑定的Tag信息*/
	err = tx.Delete(&EntitySets.Tags{}, "VID=?", del.VideoID).Error
	if err != nil {
		tx.Rollback()
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
