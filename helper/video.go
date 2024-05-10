package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"github.com/gin-gonic/gin"
	"os/exec"
	"path"
)

// Other2MP4 将其他格式的视频转换为MP4格式
func Other2MP4(videoPath string) error {
	outputPath := path.Dir(videoPath)
	ffmpegArgs := []string{
		"-hwaccel_output_format", "cuda", //设置Nvidia GPU硬件加速
		"-c:v", "h264_cuvid", //设置解码器
		"-i", videoPath,
		"-c:v", "h264_nvenc", //设置编码器
		"-c:a", "aac",
		"-strict", "normal",
		outputPath + "/converted.mp4",
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	err := cmd.Run()

	return err
}

// UpdateVideoFieldForUpdate 更新视频某个字段(悲观锁)
func UpdateVideoFieldForUpdate(c *gin.Context, videoID int64, field string, change int) error {
	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
	err := EntitySets.UpdateVideoField(tx, videoID, field, change)
	funcName, _ := c.Get("funcName")
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, funcName.(string), 5000, err.Error())
		return err
	}
	tx.Commit()
	return nil
}
