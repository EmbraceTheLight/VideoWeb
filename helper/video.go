package helper

import (
	"fmt"
	"os/exec"
	"path"
)

// Other2MP4 将其他格式的视频转换为MP4格式
func Other2MP4(videoPath string) error {
	outputPath := path.Dir(videoPath)
	ffmpegArgs := []string{
		"-i", videoPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "normal",
		outputPath + "/converted.mp4",
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	err := cmd.Run()
	fmt.Println("转换MP4成功！")

	return err
}
