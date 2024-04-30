package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// UploadVideo
// @Tags Video API
// @summary 用户上传视频
// @Accept multipart/form-data
// @Produce json
// @Param userID query string true "用户ID"
// @Param  title formData string true "视频标题"
// @Param uploadVideo formData file true "视频"
// @Param videoCover formData file true "视频封面"
// @Param class formData string true "视频分类"  Enums(娱乐,教育,科技,知识,健康,旅行,探险,美食,时尚,音乐,舞蹈,体育,健身,历史,文化,游戏,电影,搞笑,资讯)
// @Param tags formData []string false "视频标签"  collectionFormat(multi)
// @Param  description formData string false "视频描述"
// @Router /video/upload [post]
func UploadVideo(c *gin.Context) {
	UserID := c.Query("userID")

	/*检查视频后缀名*/
	FH, _ := c.FormFile("uploadVideo")
	if Utilities.CheckVideoExt(FH.Filename) != true {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ErrorVideoFormat, "视频格式错误或不支持此视频格式")
		return
	}

	/*创建对应目录*/
	videoPath, err := Utilities.Mkdir(UserID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ReadFileFailed, "创建视频目录失败:"+err.Error())
		return
	}
	fmt.Println("videoPath:", videoPath)

	/*在对应目录下创建并写入文件*/
	baseVideoName := path.Base(videoPath)
	videoFileName := videoPath + baseVideoName + path.Ext(FH.Filename)
	fmt.Println("videoFileName:", videoFileName)
	defer func() { //处理发生错误的情况，以便删除已经上传的视频.注意defer注册函数的顺序，否则会因为删除未关闭的文件导致删除失败
		if err != nil {
			e := os.RemoveAll(videoPath)
			if e == nil {
				fmt.Println("删除视频成功！")
			} else {
				Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.DeleteVideoFailed, "删除视频失败:"+e.Error())
			}
		}
	}()

	/*创建文件并将视频数据写入文件*/
	err = Utilities.WriteToNewFile(FH, videoFileName)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败:"+err.Error())
		return
	}
	err = logic.MakeDASHSegments(videoFileName)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "创建DASH视频切片失败:"+err.Error())
		return
	}
	/*将对应数据插入数据库*/
	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*将视频数据插入数据库*/
	VID, err := logic.CreateVideoRecord(tx, c, videoFileName, FH.Size)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败:"+err.Error())
		return
	}

	/*将视频标签数据插入数据库*/
	Tags := c.PostFormArray("tags")
	if len(Tags) != 0 {
		tags := make([]*EntitySets.Tags, len(Tags))
		for i, tag := range Tags {
			tags[i] = &EntitySets.Tags{
				Tag: tag,
				VID: VID,
			}
		}
		err = tx.Model(&EntitySets.Tags{}).Create(&tags).Error
		if err != nil {
			tx.Rollback()
			Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败:"+err.Error())
			return
		}
	}

	/*插入视频封面图片信息*/
	//检查封面后缀名
	Cover, _ := c.FormFile("videoCover")
	if Utilities.CheckPicExt(Cover.Filename) != true {
		err = errors.New("图片格式错误或不支持此图片格式") //便于上面的defer捕获错误
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ImageFormatError, err.Error())
		tx.Rollback()
		return
	}
	//打开并读取视频封面文件
	coverFile, err := Cover.Open()
	defer coverFile.Close()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.OpenFileFailed, "打开文件失败:"+err.Error())
		tx.Rollback()
		return
	}
	coverData, err := io.ReadAll(coverFile)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ReadFileFailed, "读取文件内容失败:"+err.Error())
		tx.Rollback()
		return
	}
	/*将视频封面图片信息插入数据库*/
	err = tx.Model(&EntitySets.Video{}).Where("videoID=?", VID).Update("VideoCover", coverData).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.CreateVideoCoverFailed, "上传视频封面失败")
		tx.Rollback()
		return
	}
	//TODO:更新用户经验值

	tx.Commit()
	Utilities.SendJsonMsg(c, 200, "上传视频成功")
}

// DeleteVideo
// @Tags Video API
// @summary 用户删除视频
// @Accept multipart/form-data
// @Produce json
// @Param VideoID query string true "视频ID"
// @Router /video/delete [Delete]
func DeleteVideo(c *gin.Context) {
	VID := c.Query("VideoID")
	var del = new(EntitySets.Video)
	err := DAO.DB.Model(&EntitySets.Video{}).Where("videoID=?", VID).First(&del).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.GetVideoInfoFailed, "获取视频信息失败")
		return
	}
	err = logic.DeleteVideo(del)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.DeleteVideoFailed, "删除用户视频失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "删除用户视频成功")
}

// DownloadVideo
// @Tags Video API
// @summary 用户下载视频(根据视频ID下载视频)
// @Accept json
// @Produce octet-stream
// @Param VideoID query string true "用户要下载的视频ID"
// @Router /video/download [get]
func DownloadVideo(c *gin.Context) {
	VID := c.Query("VideoID")
	videoInfo := new(EntitySets.Video)
	err := DAO.DB.Where("videoID=?", VID).First(&videoInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DownloadVideo", define.GetVideoInfoFailed, "获取视频信息失败"+err.Error())
		return
	}

	fileExt := path.Ext(videoInfo.Path)
	/*
		NOTE:要使用PathEscape而非QueryEscape
		因为PathEscape会将空格‘ ’转义为%20，加号‘+’不变，能被浏览器正确识别原文件名称
		QueryEscape会将空格‘ ’转义为‘+’，加号‘+’转义为%2B，浏览器会将原文件中的空格‘ ’识别为加号‘+’
		例：若使用QueryEscape，则：《Ori and the Blind Forest》前身介绍-->《Ori+and+the+Blind+Forest》前身介绍
	*/
	retFileName := url.PathEscape(videoInfo.Title + fileExt) //将非ASCII码和一些特殊字符转换转换为对应编码格式

	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+retFileName)
	fmt.Println("file name:", videoInfo.Title+fileExt)
	c.Header("Content-Type", "application/octet-stream")
	c.File(videoInfo.Path)
}

// StreamTransmission
// @Tags Video API
// @summary 流式传输视频
// @Accept json
// @Produce octet-stream
// @Param VideoID query string true "要传输的视频ID"
// @Router /video/StreamTransmission [get]
func StreamTransmission(c *gin.Context) {
	VID := c.Query("VideoID")
	videoInfo, err := EntitySets.GetVideoInfoByID(VID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
		return
	}

	file, err := os.Open(videoInfo.Path)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.OpenFileFailed, "打开文件失败:"+err.Error())
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.OpenFileFailed, "打开文件失败:"+err.Error())
		return
	}

	rangeHeader := c.GetHeader("Range")
	rangeParts := strings.Split(rangeHeader, "=") //分离出两个部分:Byte和start-end
	start, end := logic.ParseRange(rangeParts[1]) //分离出start和end
	if end == -1 {                                //请求的是最后一块视频
		end = int64(stat.Size()) - 1
	}
	_, err = file.Seek(start, 0) //从第0个字节定位文件指针到第start个字节处
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.OpenFileFailed, "打开文件失败:"+err.Error())
		return
	}

	c.Header("Accept-Ranges", "bytes")
	//c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
	c.Header("Content-Range", "bytes "+rangeParts[1]+"/"+fmt.Sprintf("%d", stat.Size()))
	c.Header("Transfer-Encoding", "chunked")

	_, err = io.CopyN(c.Writer, file, end-start+1)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.ReadFileFailed, "读取文件失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, 200, "流式传输视频成功")
}

// DASHStreamTransmission
// @Tags Video API
// @summary DASH流式传输视频
// @Accept json
// @Produce octet-stream
// @Param filePath query string true "要传输的视频路径"
// @Router /video/DASHStreamTransmission [get]
func DASHStreamTransmission(c *gin.Context) {
	filePath := c.Query("filePath")
	fmt.Println("filePath:", filePath)
	exist := logic.CheckFileIsExist(filePath)
	if !exist {
		Utilities.SendErrMsg(c, "service::Videos::DASHStreamTransmission", define.GetVideoInfoFailed, "请求的文件不存在")
		return
	}
	c.File(filePath)
}

// OfferMpd
// @Tags Video API
// @summary 提供DASH所需的.mpd文件
// @Accept json
// @Produce octet-stream
// @Param filePath query string true "要传输的视频ID"
// @Router /video/OfferMpd [get]
func OfferMpd(c *gin.Context) {
	basePath := c.Query("filePath")
	c.File(basePath + "/output.mpd")
}

// GetVideoInfo
// @Tags Video API
// @summary 提供视频信息详情
// @Accept json
// @Produce octet-stream
// @Param VideoID query string true "要获取的视频ID"
// @Router /video/getVideoDetail [get]
func GetVideoInfo(c *gin.Context) {
	VID := c.Query("VideoID")
	var videoInfo = new(EntitySets.Video)
	err := DAO.DB.Where("VideoID=?", VID).Preload("Comments").
		Preload("Tags").Preload("Barrages").First(&videoInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::GetVideoInfo", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code":     http.StatusOK,
		"data":     videoInfo,
		"basePath": "./" + path.Dir(videoInfo.Path),
	})
}
