package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"gorm.io/gorm"
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
// 注意:如果上层逻辑只需要更新一个数据,则传入tx为nil,该函数自动开启事务进行处理
// 否则,函数调用者(位于logic层)需要自行传入tx,并在函数结束后提交或回滚事务
func UpdateVideoFieldForUpdate(videoID int64, field string, change int, tx *gorm.DB) error {
	var err error
	if tx == nil {
		tx = DAO.DB.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			}
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
		err = EntitySets.UpdateVideoField(tx, videoID, field, change)
		if err != nil {
			return err
		}
		tx.Commit()
	} else {
		err = EntitySets.UpdateVideoField(tx, videoID, field, change)
		if err != nil {
			return err
		}
	}
	return nil
}
