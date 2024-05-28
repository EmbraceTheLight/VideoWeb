package service

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostComment
// @Tags Comment API
// @summary 用户评论视频
// @Accept multipart/form-data
// @Produce json
// @Param VideoID path string true "用户要评论的视频ID"
// @Param To query string false "用户要评论的对象ID"
// @Param Authorization header string true "token"
// @Param CommentContent formData string true "用户要评论的内容"
// @Success 200 {string}  json "{"code":"200","msg":"发送评论成功"}"
// @Router /Comment/{VideoID}/Comment [post]
func PostComment(c *gin.Context) {
	Utilities.AddFuncName(c, "service::Comments::PostComment")
	videoID := Utilities.String2Int64(c.Param("VideoID"))
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	userID := UID
	Content := c.PostForm("CommentContent")
	To := Utilities.String2Int64(c.DefaultQuery("To", "-1"))

	Country, City := logic.GetUserIpInfo(c)
	if Country == "" {
		Country = "未知地区"
	}
	var comment = &EntitySets.Comments{
		MyModel:   define.MyModel{},
		CommentID: logic.GetUUID(),
		UID:       userID,
		To:        To,
		VID:       videoID,
		Content:   Content,
		IPAddress: Country + " " + City,
	}
	err := logic.AddCommentToVideo(c, comment)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	//根据视频ID获得视频Up主ID
	//up, err := EntitySets.GetVideoInfoByID(tx, videoID)
	//if err != nil {
	//	Utilities.SendErrMsg(c, "service::Comments::PostComment", define.QueryUserError, "获取用户信息失败:"+err.Error())
	//	tx.Rollback()
	//	return
	//}
	/*使用websocket通知被评论的视频up主(如果该用户在线)，并把“被评论”这一事件作为msg写入数据库，
	这样即使视频up主当时未在线，也能通过检索数据库的方式得知自己有新消息*/

	Utilities.SendJsonMsg(c, http.StatusOK, "发送评论成功")
}

// LikeOrDislikeComment
// @Tags Comment API
// @summary 用户点赞评论
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param IsLike query bool true "是点赞还是点踩"
// @Param CommentID query string true "用户要点赞/点踩的评论ID"
// @Success 200 {string}  json "{"code":"200","msg":"操作成功"}"
// @Router /Comment/{VideoID}/LikeOrDislikeComment [put]
func LikeOrDislikeComment(c *gin.Context) {
	Utilities.AddFuncName(c, "service::Comments::LikeOrDislikeComment")
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	commentID := Utilities.String2Int64(c.Query("CommentID"))
	videoID := Utilities.String2Int64(c.Param("VideoID"))
	Like := c.Query("isLike")
	var isLike bool
	if Like == "true" {
		isLike = true
	} else {
		isLike = false
	}

	err := logic.LikeOrDislikeComment(c, UID, commentID, videoID, isLike)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "操作成功")
}

// UndoLikeOrDislikeComment
// @Tags Comment API
// @summary 用户撤销点赞/点踩评论
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param IsLike query bool true "是点赞还是点踩"
// @Param CommentID query string true "用户撤销点赞/点踩的评论ID"
// @Success 200 {string}  json "{"code":"200","msg":"发送评论成功"}"
// @Router /Comment/{VideoID}/UndoLikeOrDislikeComment [put]
func UndoLikeOrDislikeComment(c *gin.Context) {
	Utilities.AddFuncName(c, "service::Comments::LikeOrDislikeComment")
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	commentID := Utilities.String2Int64(c.Query("CommentID"))
	videoID := Utilities.String2Int64(c.Param("VideoID"))
	Like := c.Query("isLike")
	var isLike bool
	if Like == "true" {
		isLike = true
	} else {
		isLike = false
	}

	err := logic.UndoLikeOrDislikeComment(c, UID, commentID, videoID, isLike)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "点赞成功")
}
