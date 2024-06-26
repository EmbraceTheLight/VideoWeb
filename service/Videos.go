package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/Utilities/logf"
	"VideoWeb/cache/videoCache"
	"VideoWeb/define"
	"VideoWeb/logic"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

// UploadVideoFile
// @Tags Video API
// @summary 用户上传视频
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param uploadVideo formData file true "视频"
// @Router /video/VideoFile [post]
func UploadVideoFile(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::videos::UploadVideoFile")
	var err error

	u, _ := c.Get("user")
	UserID := logic.GetUserID(u)
	//检查视频后缀名
	FH, _ := c.FormFile("uploadVideo")
	if err := Utilities.CheckVideoExt(FH.Filename); err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.ErrorVideoFormat, err.Error())
		return
	}

	//创建对应目录
	videoPath, err := Utilities.Mkdir(strconv.FormatInt(UserID, 10))
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.ReadFileFailed, "创建视频目录失败:"+err.Error())
		return
	}

	//在对应目录下创建并写入文件
	baseVideoName := path.Base(videoPath)                              //获得刚才创建的目录名作为文件名
	videoFileName := videoPath + baseVideoName + path.Ext(FH.Filename) //拼接视频文件名
	defer func() {
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
	//写入视频源数据
	err = Utilities.WriteToNewFile(FH, videoFileName)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "上传视频失败:"+err.Error())
		return
	}
	//写入分段视频数据
	err = logic.MakeDASHSegments(videoFileName)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "创建DASH视频切片失败:"+err.Error())
		return
	}

	//将视频数据插入数据库
	VID, err := logic.CreateVideoRecord(DAO.DB, c, UserID, videoFileName, FH.Size)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "上传视频失败:"+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"ID":   VID,
	})

}

// UploadVideoInfo
// @Tags Video API
// @Param ID path string true "视频ID"
// @Param title formData string true "视频标题"
// @Param Authorization header string true "token"
// @Param videoCover formData file true "视频封面"
// @Param isUpload query bool true "是否上传"
// @Param tags formData []string true "视频标签"  collectionFormat(multi)
// @Param class formData string true "视频分类"  Enums(娱乐,教育,科技,知识,健康,旅行,探险,美食,时尚,音乐,舞蹈,体育,健身,历史,文化,游戏,电影,搞笑,资讯)
// @Param  description formData string false "视频描述"
// @Router /video/{ID}/VideoInfo [post]
func UploadVideoInfo(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::videos::UploadVideoInfo")
	var err error

	u, _ := c.Get("user")
	userID := logic.GetUserID(u)
	videoID := Utilities.String2Int64(c.Param("ID"))
	isUpload := c.Query("isUpload")

	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, videoID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoInfo", 5000, err.Error())
		return
	}

	if isUpload == "false" {
		err = os.RemoveAll(path.Dir(videoInfo.Path))
		if err != nil {
			Utilities.SendErrMsg(c, "service::Videos::UploadVideoInfo", 5000, "取消视频上传失败:"+err.Error())
			return
		}
		err = EntitySets.DeleteVideoInfoByVideoID(DAO.DB, videoID)
		if err != nil {
			Utilities.SendErrMsg(c, "service::Videos::UploadVideoInfo", 5000, "取消视频上传失败:"+err.Error())
			return
		}
		c.String(http.StatusOK, "取消视频上传成功")
		return
	}

	/*开启事务,将对应数据插入数据库*/
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//将视频标签数据插入数据库
	Tags := c.PostFormArray("tags")
	if len(Tags) != 0 {
		tags := make([]*EntitySets.Tags, len(Tags))
		for i, tag := range Tags {
			tags[i] = &EntitySets.Tags{
				Tag: tag,
				VID: videoID,
			}
		}
		err = EntitySets.InsertTags(tx, tags)
		if err != nil {
			Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "上传视频失败:"+err.Error())
			return
		}
	}

	/*插入视频封面图片信息*/
	//检查封面后缀名
	Cover, _ := c.FormFile("videoCover")
	if err = Utilities.CheckPicExt(Cover.Filename); err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.ImageFormatError, err.Error())
		return
	}

	//打开并读取视频封面文件
	coverPath := path.Join(path.Dir(videoInfo.Path), "cover"+path.Ext(Cover.Filename))
	err = Utilities.WriteToNewFile(Cover, coverPath)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile::Utilities.OpenAndReadFile", define.OpenFileFailed, "打开或读取文件失败:"+err.Error())
		return
	}

	/*更新视频用户指定的有关信息*/
	// 获得hh:mm:ss格式的视频时长
	duration, err := logic.GetVideoDuration(videoInfo.Path)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile::Utilities.OpenAndReadFile", 5000, "获取视频时长失败:"+err.Error())
		return
	}
	durationStr, _ := Utilities.SecondToTime(duration)

	//更新视频记录信息
	err = tx.Model(&EntitySets.Video{}).Where("video_id = ?", videoID).Updates(map[string]any{
		"title":       c.PostForm("title"),
		"description": c.PostForm("description"),
		"class":       c.PostForm("class"),
		"cover_path":  coverPath,
		"duration":    durationStr}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "更新视频信息失败:"+err.Error())
		return
	}

	//更新用户经验值
	err = logic.AddExpForUploadVideo(c, userID, tx)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::UploadVideoFile", define.UploadVideoFailed, "更新用户经验值失败")
		return
	}

	tx.Commit()
	Utilities.SendJsonMsg(c, http.StatusOK, "上传视频成功")

	/*将信息更新到redis中*/
	vc := videoCache.MakeVideoCache()
	err = vc.MakeVideoInfo(context.Background(), videoID)
	if err != nil {
		logf.WriteErrLog("service::Videos::UploadVideoInfo", err.Error())
	}
}

// DeleteVideo
// @Tags Video API
// @summary 用户删除视频
// @Accept multipart/form-data
// @Produce json
// @Param ID path string true "视频ID"
// @Router /video/{ID}/Delete [Delete]
func DeleteVideo(c *gin.Context) {
	VID := Utilities.String2Int64(c.Param("ID"))
	var del, err = EntitySets.GetVideoInfoByID(DAO.DB, VID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DeleteVideo", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
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
// @Param ID path string true "用户要下载的视频ID"
// @Router /video/{ID}/Download [get]
func DownloadVideo(c *gin.Context) {
	VID := c.Param("ID")
	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, Utilities.String2Int64(VID))
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::DownloadVideo", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
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
// @Param ID path string true "要传输的视频ID"
// @Router /video/{ID}/StreamTransmission [get]
func StreamTransmission(c *gin.Context) {
	VID := c.Param("ID")
	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, Utilities.String2Int64(VID))
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::StreamTransmission", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
		return
	}
	fmt.Println("videoInfo:", videoInfo)
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
// @Param filePath query string true "视频所在路径"
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
// @Param ID path string true "要获取的视频ID"
// @Param Authorization header string true "token"
// @Router /video/{ID}/VideoDetail [get]
func GetVideoInfo(c *gin.Context) {
	c.Set("funcName", "Service::Videos::GetVideoInfo")
	VID := Utilities.String2Int64(c.Param("ID"))
	var UID int64
	var videoInfo = new(EntitySets.Video)
	err := DAO.DB.Where("video_id=?", VID).
		Preload("Tags").Preload("Barrages").First(&videoInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::GetVideoInfo", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
		return
	}

	//更新UserVideo表：若没有对应记录则插入
	u, _ := c.Get("user")
	if u != nil {
		UID = u.(*logic.UserClaims).UserId
		err = logic.InsertUserVideoIfNotExist(UID, VID)
		if err != nil {
			Utilities.SendErrMsg(c, "service::Videos::GetVideoInfo", define.GetVideoInfoFailed, "获取视频信息失败:"+err.Error())
			return
		}
	}

	//该视频观看次数+1
	err = logic.AddVideoViewCount(c, videoInfo.VideoID)
	if err != nil {
		return
	}

	//TODO:添加历史观看记录

	//查找对应的UserVideo记录
	uv, err := RelationshipSets.GetUserVideoRecord(DAO.DB, UID, videoInfo.VideoID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     videoInfo,
		"status":   uv,
		"basePath": path.Dir(videoInfo.Path),
	})
}

// LikeOrUndoLike
// @Tags Video API
// @summary 用户点赞/取消点赞视频
// @Accept json
// @Produce json
// @Param ID path string true "要获取的视频ID"
// @Param Authorization header string true "token"
// @Router /video/{ID}/LikeOrUndoLike [put]
func LikeOrUndoLike(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Videos::LikeOrUndoLike")
	u, _ := c.Get("user")
	UserID := logic.GetUserID(u)
	VideoID := Utilities.String2Int64(c.Param("ID"))

	/*查找对应的UserVideo记录*/
	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, VideoID)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}

	/*检查用户是否已经点赞*/
	//查找对应的UserVideo记录
	uv, err := RelationshipSets.GetUserVideoRecord(DAO.DB, UserID, videoInfo.VideoID)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	//修改用户点赞状态与视频点赞数
	err = logic.UpdateVideoLikeStatus(c, UserID, videoInfo.UID, videoInfo.VideoID, uv.IsLike)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "点赞成功")
}

// ThrowShell
// @Tags Video API
// @summary 用户为视频扔贝壳，由前端负责检查贝壳数量是否足够投喂
// @Accept json
// @Produce json
// @Param ID path string true "要获取的视频ID"
// @Param Authorization header string true "token"
// @Param shells query int true "投贝壳的贝壳数量"
// @Router /video/{ID}/ThrowShell [put]
func ThrowShell(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Videos::ThrowShell")
	VID := Utilities.String2Int64(c.Param("ID"))
	u, _ := c.Get("user")
	TSUID := logic.GetUserID(u)
	tmpShells := c.Query("shells")

	videoInfo, err := EntitySets.GetVideoInfoByID(DAO.DB, VID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::ThrowShell", define.GetVideoInfoFailed, "投贝壳失败:"+err.Error())
		return
	}

	shells, _ := strconv.Atoi(tmpShells)
	//检查用户贝壳数量是否足够
	userInfo, err := EntitySets.GetUserInfoByID(DAO.DB, TSUID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Videos::ThrowShell", 5000, "获取用户信息失败:"+err.Error())
		return
	}
	if uint32(shells) > userInfo.Shells {
		Utilities.SendErrMsg(c, "service::Videos::ThrowShell", 4000, "贝壳数量不足")
		return
	}
	//更新视频，用户贝壳数量以及经验值相关信息
	err = logic.UpdateShells(c, videoInfo, TSUID, shells)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "投贝壳成功")
}

// GetVideoList
// @Tags Video API
// @summary 主页根据分类获取视频列表
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param class query string false "根据分类给出视频列表"
// @Router /video/VideoList [get]
func GetVideoList(c *gin.Context) {
	class := c.DefaultQuery("class", "recommend")
	videoInfo, err := logic.GetVideoListByClass(c, class)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": videoInfo,
	})
}

// GetVideoComments
// @Tags Comment API
// @summary 获取视频评论列表
// @Accept json
// @Produce json
// @Param ID path string true "视频ID"
// @Param Authorization header string true "token"
// @Param commentNums query int true "请求的评论数量"
// @Param offset query int true "评论的偏移量"
// @Param order query string false "评论排序方式:default,likes:按点赞数量排序;newest:按最新发布排序"
// @Success 200 {string}  json "{"code":"200","msg":"发送评论成功"}"
// @Router /video/{ID}/Comments [get]
func GetVideoComments(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Videos::GetVideoComments")
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	VID := Utilities.String2Int64(c.Param("ID"))
	Offset := Utilities.String2Int(c.Query("offset"))
	CommentNums := Utilities.String2Int(c.Query("commentNums"))
	order := c.DefaultQuery("order", "default")

	comments, err := logic.GetVideoCommentsList(c, VID, UID, order, Offset, CommentNums)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取评论成功",
		"data": comments,
	})
}

// SearchVideos
// @Tags Video API
// @summary 根据关键字搜索视频
// @Accept json
// @Produce json
// @Param videoNums query int true "搜索的视频数量"
// @Param Authorization header string true "token"
// @Param offset query int true "搜索视频的偏移量"
// @Param key query string true "关键字"
// @Param sortOrder query string false "视频的排序方式,default:按热度排序,newest:按最新发布排序,mostPlay:按播放量排序,mostBarrage:按弹幕数量排序,mostFavorite:按收藏数量排序"
// @Success 200 {string}  json "{"code":"200","data": videoInfo,"msg":"搜索视频成功"}"
// @Router /video/SearchVideos [get]
func SearchVideos(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Videos::SearchVideos")
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	offset := Utilities.String2Int(c.Query("offset"))
	videoNums := Utilities.String2Int(c.Query("videoNums"))
	key := c.Query("key")
	sortOrder := c.DefaultQuery("sortOrder", "default")

	videoInfo, err := logic.GetVideosByKey(c, key, sortOrder, offset, videoNums)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}

	//添加搜索记录
	err = EntitySets.InsertSearchRecord(DAO.DB, UID, key)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": videoInfo,
		"msg":  "搜索视频成功",
	})
}

//func Test(c *gin.Context) {
//	_, err := EntitySets.GetVideoInfoByID(DAO.DB, 1)
//	c.String(200, err.Error())
//}
