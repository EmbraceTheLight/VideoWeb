package helper

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"fmt"
	"gorm.io/gorm"
	"os/exec"
	"path"
)

// Other2MP4 将其他格式的视频转换为MP4格式
func Other2MP4(videoPath string) error {
	outputPath := path.Dir(videoPath)
	//ffmpegArgs := []string{
	//	"-hwaccel_output_format", "cuda", //设置Nvidia GPU硬件加速
	//	"-c:v", "h264_cuvid", //设置解码器
	//	"-i", videoPath,
	//	"-c:v", "h264_nvenc", //设置编码器
	//	"-c:a", "aac",
	//	"-strict", "normal",
	//	outputPath + "/converted.mp4",
	//}
	ffmpegArgs := []string{
		"-i", videoPath,
		//"-c:v", "libx264", // 使用 CPU 进行编码
		"-c:v", "copy",
		"-c:a", "copy",
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
		tx.Set("gorm:query_option", "FOR UPDATE")
		//添加行级锁(悲观)
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

// DeleteCommentWithStatus 删除评论及其状态（如用户点赞等信息）
func DeleteCommentWithStatus(videoID int64, tx *gorm.DB) error {
	var commentsIDs []int64
	err := tx.Model(EntitySets.Comments{}).Where("video_id =?", videoID).Select("comment_id").Find(&commentsIDs).Error
	if err != nil {
		return fmt.Errorf("helper/video.DeleteCommentWithStatus: %w", err)
	}

	for _, id := range commentsIDs {
		err = RelationshipSets.DeleteUserLikedCommentRecordByCommentID(tx, id)
		if err != nil {
			return fmt.Errorf("helper/video.DeleteCommentWithStatus: %w", err)
		}
	}

	err = EntitySets.DeleteCommentRecordsByVideoID(tx, videoID)
	if err != nil {
		return fmt.Errorf("helper/video.DeleteCommentWithStatus: %w", err)
	}
	//TODO: Redis中删除评论状态记录

	return nil
}

// GetVideosByClass 根据分类ID获取视频信息
//func GetVideosByClass(class string) ([]*EntitySets.Video, error) {
//
//}
