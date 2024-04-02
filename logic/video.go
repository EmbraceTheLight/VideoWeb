package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
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
