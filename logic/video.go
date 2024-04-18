package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetVideoInfoByID 根据视频ID获得视频信息
func GetVideoInfoByID(VID string) (*EntitySets.Video, error) {
	var info = new(EntitySets.Video)
	err := DAO.DB.Where("videoID=?", VID).First(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

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
