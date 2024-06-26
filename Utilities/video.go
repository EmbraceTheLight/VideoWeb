package Utilities

import (
	"VideoWeb/define"
	"fmt"
	"os"
	"strings"
	"time"
)

// RoundOff 四舍五入函数
func RoundOff(f float64) int64 {
	var decimal = f - float64(int(f))
	if decimal >= 0.5 {
		return int64(f) + 1
	}
	return int64(f)
}

// SecondToTime 秒数转时间字符串以及时分秒数组，可以按需取用两个返回值中的任何一个
func SecondToTime(second int64) (string, []int64) {
	hour := second / 3600
	minute := (second % 3600) / 60
	second = second % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second), []int64{hour, minute, second}
}

// Mkdir 利用用户ID和当前时间来创建视频对应目录
func Mkdir(uid string) (path string, err error) {
	var b strings.Builder
	curTime := time.Now().Format("2006-01-02T150405")
	b.WriteString(define.BaseDir)
	//TODO:可以改造这个函数，令其参数为string切片以更灵活地创建目录
	b.WriteString(uid)
	b.WriteString("/")
	b.WriteString("videos")
	b.WriteString("/")
	b.WriteString(curTime)
	b.WriteString("/")
	videoDirPath := b.String()
	err = os.MkdirAll(videoDirPath, os.ModePerm)
	return videoDirPath, err
}
