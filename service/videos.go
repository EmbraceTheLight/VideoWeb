package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	fmt.Println("UserID:", UserID)
	fmt.Println("Title:", Title)
	fmt.Printf("%+v", Tags)
	fmt.Println("description:", Description)
	fmt.Printf("%#v", FH)
	fname := FH.Filename

	//TODO:检查视频后缀名
	ext := path.Ext(fname)
	println("ext:", ext)
	if _, ok := define.VideoExtCheck[ext]; !ok {
		Utilities.SendJsonMsg(c, 4017, "视频格式错误或不支持此视频格式")
		return
	}

	uploadFile, err := FH.Open()
	if err != nil {
		Utilities.SendJsonMsg(c, 5010, "打开文件失败"+err.Error())
		return
	}
	defer uploadFile.Close()

	data, err := io.ReadAll(uploadFile)
	if err != nil {
		Utilities.SendJsonMsg(c, 5011, "读取文件内容失败"+err.Error())
		return
	}

	//TODO:若不存在相关目录，则创建一个
	var b strings.Builder
	b.WriteString(define.VideoSavePath)
	b.WriteString(UserID)
	videoDirPath := b.String()
	println(videoDirPath)
	err = os.MkdirAll(videoDirPath, os.ModePerm)
	if err != nil {
		Utilities.SendJsonMsg(c, 5011, "读取文件内容失败"+err.Error())
		return
	}

	//TODO:在对应目录下创建并写入文件
	Time := time.Now().Format("2006-01-02T150405") //利用当前时间生成文件名，避免文件名重复的情况
	b.WriteString("\\")
	b.WriteString(Time)
	b.WriteString(ext) //拼接路径、文件名以及文件后缀名
	videoFilePath := b.String()
	println(videoFilePath)
	vf, err := os.Create(videoFilePath) //创建文件
	if err != nil {
		err = os.Remove(videoFilePath)
		println("err in Creating File", err.Error())
		Utilities.SendJsonMsg(c, 5011, "创建文件失败"+err.Error())
		return
	}
	_, err = vf.Write(data)
	if err != nil {
		Utilities.SendJsonMsg(c, 5011, "写入文件失败"+err.Error())
		return
	}
	_ = vf.Close()
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
	}
	err = DAO.DB.Model(&EntitySets.Video{}).Create(&video).Error
	if err != nil {
		defer func() {
			err := os.Remove(videoFilePath)
			if err == nil {
				fmt.Println("删除视频成功！")
			} else {
				fmt.Println("删除视频失败:", err.Error())
			}
		}()
		println("err in Removing File:", err.Error())
		tx.Rollback()
		Utilities.SendJsonMsg(c, 5016, "上传视频失败"+err.Error())
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
			defer func() {
				err := os.Remove(videoFilePath)
				if err == nil {
					fmt.Println("删除视频成功！")
				} else {
					fmt.Println("删除视频失败:", err.Error())
				}
			}()
			println("err in Removing File:", err.Error())
			tx.Rollback()
			Utilities.SendJsonMsg(c, 5016, "上传视频失败"+err.Error())
			return
		}
	}

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
		Utilities.SendJsonMsg(c, 5017, "获取视频信息失败")
		return
	}
	//TODO:从硬盘中删除对应视频信息
	err = os.RemoveAll(del.Path)
	if err != nil {
		Utilities.SendJsonMsg(c, 5018, "删除用户视频失败:"+err.Error())
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
		Utilities.SendJsonMsg(c, 5018, "删除用户视频失败:"+err.Error())
		tx.Rollback()
		return
	}
	//TODO:从数据库中删除与视频绑定的Tag信息
	err = tx.Delete(&EntitySets.Tags{}, "VID=?", del.VideoID).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5018, "删除用户视频失败:"+err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	Utilities.SendJsonMsg(c, 200, "删除用户视频成功")
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
		Utilities.SendJsonMsg(c, 5017, "获取视频信息失败"+err.Error())
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

	c.Header("Transfer-Encoding", "chunked")
	c.File(videoInfo.Path)

	// 使用 http.ServeContent 来处理文件下载
	//http.ServeContent(c.Writer, c.Request, "", time.Now(), file)
	Utilities.SendJsonMsg(c, 200, "下载文件成功")
}
