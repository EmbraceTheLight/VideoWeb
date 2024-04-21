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
	"time"
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
	FH, _ := c.FormFile("uploadVideo")
	Title := c.PostForm("title")
	Tags := c.PostFormArray("tags")
	Description := c.PostForm("description")
	Class := c.PostForm("class")
	Cover, _ := c.FormFile("videoCover")
	fmt.Println("CoverName:", Cover.Filename)
	fname := FH.Filename

	//TODO:检查视频后缀名
	videoExt := path.Ext(fname)
	if _, ok := define.VideoExtCheck[videoExt]; !ok {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ErrorVideoFormat, "视频格式错误或不支持此视频格式")
		return
	}

	uploadFile, err := FH.Open()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.OpenFileFailed, "打开文件失败"+err.Error())
		return
	}
	defer uploadFile.Close()

	//TODO:若不存在相关目录，则创建一个
	var b strings.Builder
	b.WriteString(define.VideoSavePath)
	b.WriteString(UserID + "/")
	videoDirPath := b.String()
	println(videoDirPath)
	err = os.MkdirAll(videoDirPath, os.ModePerm)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ReadFileFailed, "读取文件内容失败"+err.Error())
		return
	}

	//TODO:在对应目录下创建并写入文件
	Time := time.Now().Format("2006-01-02T150405") //利用当前时间生成文件名，避免文件名重复的情况
	//b.WriteString("\\")
	b.WriteString(Time)
	b.WriteString(videoExt) //拼接路径、文件名以及文件后缀名
	videoFilePath := b.String()
	fmt.Println(videoFilePath)
	//t, err := logic.GetVideoDuration(videoFilePath)
	//if err != nil {
	//	Utilities.SendErrMsg(c, define.UploadVideoFailed, "上传视频失败:"+err.Error())
	//	return
	//}
	//videoTime := Utilities.SecondToTime(t)
	defer func() { //处理发生错误的情况，以便删除已经上传的视频.注意defer注册函数的顺序，否则会因为删除未关闭的文件导致删除失败
		if err != nil {
			e := os.Remove(videoFilePath)
			if e == nil {
				fmt.Println("删除视频成功！")
			} else {
				fmt.Println("删除视频失败:", e.Error())
			}
		}
	}()

	vf, err := os.Create(videoFilePath) //创建文件
	defer vf.Close()
	if err != nil {
		println("err in Creating File", err.Error())
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败"+err.Error())
		return
	}

	_, err = io.Copy(vf, uploadFile)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", 5011, "写入文件失败"+err.Error())
		return
	}
	t, err := logic.GetVideoDuration(videoFilePath)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败:"+err.Error())
		return
	}
	videoTime := Utilities.SecondToTime(t)
	//将数据插入数据库
	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//将视频数据插入数据库
	VID := logic.GetUUID()
	video := &EntitySets.Video{
		MyModel:     define.MyModel{},
		VideoID:     VID,
		UID:         UserID,
		Title:       Title,
		Description: Description,
		Class:       Class,
		Path:        videoFilePath,
		Duration:    videoTime,
		Size:        FH.Size,
	}
	err = tx.Model(&EntitySets.Video{}).Create(&video).Error
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败"+err.Error())
		return
	}

	//将视频标签数据插入数据库
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
			Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.UploadVideoFailed, "上传视频失败"+err.Error())
			return
		}
	}

	/*插入视频封面图片信息*/
	//检查封面后缀名
	coverName := Cover.Filename
	coverExt := path.Ext(coverName)
	if _, ok := define.PicExtCheck[coverExt]; ok != true {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ImageFormatError, "图片格式错误或不支持此图片格式")
		err = errors.New("图片格式错误或不支持此图片格式")
		tx.Rollback()
		return
	}

	//打开并读取文件
	coverFile, err := Cover.Open()
	defer coverFile.Close()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.OpenFileFailed, "打开文件失败"+err.Error())
		tx.Rollback()
		return
	}
	coverData, err := io.ReadAll(coverFile)
	fmt.Println("coverData size:", Cover.Size)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.ReadFileFailed, "读取文件内容失败"+err.Error())
		tx.Rollback()
		return
	}

	err = tx.Model(&EntitySets.Video{}).Where("videoID=?", VID).Update("VideoCover", coverData).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideo", define.CreateVideoCoverFailed, "上传视频封面失败")
		tx.Rollback()
		fmt.Println("err:", err.Error())
		return
	}
	//TODO:更新用户经验值

	tx.Commit()

	Utilities.SendSuccessMsg(c, 200, "上传视频成功")

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
	//TODO:从硬盘中删除对应视频信息
	err = os.RemoveAll(del.Path)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.DeleteVideoFailed, "删除用户视频失败:"+err.Error())
		return
	}

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//TODO:从数据库中删除视频信息
	err = tx.Where("VideoID=?", del.VideoID).Delete(&EntitySets.Video{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.DeleteVideoFailed, "删除用户视频失败:"+err.Error())
		tx.Rollback()
		return
	}
	//TODO:从数据库中删除与视频绑定的Tag信息
	err = tx.Delete(&EntitySets.Tags{}, "VID=?", del.VideoID).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.DeleteVideoFailed, "删除用户视频失败:"+err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	Utilities.SendSuccessMsg(c, 200, "删除用户视频成功")
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
	http.ServeFile(c.Writer, c.Request, videoInfo.Path)
	c.Header("Content-Length", "")
	c.Header("Transfer-Encoding", "chunked")
	Utilities.SendSuccessMsg(c, 200, "下载文件成功")
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
	videoInfo, err := logic.GetVideoInfoByID(VID)
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
	Utilities.SendSuccessMsg(c, 200, "流式传输视频成功")
}

// GetVideoInfo
// @Tags Video API
// @summary 获取对应视频信息
// @Accept json
// @Produce json
// @Param VideoID query string true "视频ID"
// @Router /video/getVideoDetail [get]
func GetVideoInfo(c *gin.Context) {
	VID := c.Query("VideoID")
	var videoInfo = new(EntitySets.Video)
	err := DAO.DB.Where("VideoID=?", VID).Preload("Comments").
		Preload("Tags").Preload("Barrages").First(&videoInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::GetVideoInfo", define.GetVideoInfoFailed, "获取视频信息失败"+err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": videoInfo,
	})
}
