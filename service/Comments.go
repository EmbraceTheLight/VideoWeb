package service

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentToVideo
// @Tags Comment API
// @summary 用户评论视频
// @Accept json
// @Produce json
// @Param VideoID query string true "用户要评论的视频ID"
// @Param Authorization header string true "token"
// @Param CommentContent formData string true "用户要评论的内容"
// @Success 200 {string}  json "{"code":"200","msg":"发送评论成功"}"
// @Router /Comment/ToVideo [post]
func CommentToVideo(c *gin.Context) {
	VID := c.Query("VideoID")
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	Content := c.PostForm("CommentContent")
	userID := UID
	videoID := Utilities.String2Int64(VID)

	Country, City := logic.GetUserIpInfo(c)

	if Country == "" {
		Country = "未知地区"
	}
	var comment = &EntitySets.Comments{
		MyModel:   define.MyModel{},
		CommentID: logic.GetUUID(),
		UID:       userID,
		To:        -1,
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
	//	Utilities.SendErrMsg(c, "service::Comments::CommentToVideo", define.QueryUserError, "获取用户信息失败:"+err.Error())
	//	tx.Rollback()
	//	return
	//}
	/*使用websocket通知被评论的视频up主(如果该用户在线)，并把“被评论”这一事件作为msg写入数据库，
	这样即使视频up主当时未在线，也能通过检索数据库的方式得知自己有新消息*/

	Utilities.SendJsonMsg(c, http.StatusOK, "发送评论成功")
}

// CommentToOtherUser
// @Tags Comment API
// @summary 用户评论其他用户
// @Accept json
// @Produce json
// @Param VideoID query string true "用户要评论的视频ID"
// @Param Authorization header string true "token"
// @Param CommentID query string true "用户要评论的评论ID"
// @Param CommentContent formData string true "用户要评论的内容"
// @Success 200 {string}  json "{"code":"200","msg":"发送评论成功"}"
// @Router /Comment/ToUser [post]
func CommentToOtherUser(c *gin.Context) {

}
