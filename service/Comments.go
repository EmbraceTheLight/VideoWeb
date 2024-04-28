package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/Utilities/WebSocket"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"log"
)

// CommentToVideo
// @Tags Comment API
// @summary 用户评论视频
// @Accept json
// @Produce json
// @Param VideoID query string true "用户要评论的视频ID"
// @Param UserID query string true "用户ID"
// @Param CommentContent formData string true "用户要评论的内容"
// @Router /comment/ToVideo [post]
func CommentToVideo(c *gin.Context) {
	VID := c.Query("VideoID")
	UID := c.Query("UserID")
	Content := c.PostForm("CommentContent")
	CommentID := logic.GetUUID()

	Country, City := logic.GetUserIpInfo(c)
	if Country == "" {
		Country = "未知地区"
	}
	var comment = EntitySets.Comments{
		MyModel:   define.MyModel{},
		CommentID: CommentID,
		UID:       UID,
		To:        "",
		VID:       VID,
		Content:   Content,
		IPAddress: Country + " " + City,
	}
	err := DAO.DB.Create(&comment).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Comments::CommentToVideo", define.CreateCommentToVideoFailed, "创建用户评论（To视频）记录失败")
		return
	}

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
			tx.Rollback()
		}
	}()

	//根据视频ID获得视频Up主ID
	var VideoUpID string
	err = tx.Model(&EntitySets.Video{}).Where("videoID=?", VID).Pluck("UID", &VideoUpID).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Comments::CommentToVideo", define.QueryUserError, "获取用户信息失败"+err.Error())
		tx.Rollback()
		return
	}
	//使用websocket通知被评论的视频up主(如果该用户在线)，并把“被评论”这一事件作为msg写入数据库，
	//这样即使视频up主当时未在线，也能通过检索数据库的方式得知自己有新消息
	conn, ok := WebSocket.Hub.UserConnections[UID]
	liker, _ := logic.GetUserNameByID(UID)
	msg := &define.Message{
		Title: liker + "点赞了你的视频",
		Body:  "",
	}
	if ok { //TODO:Up主当前在线,待完善
		conn = conn

	}

	err = logic.CreateMessage(msg)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Comments::CommentToVideo", define.CreateMessageFailed, "创建Message失败:"+err.Error())
	}
	tx.Commit()

	Utilities.SendJsonMsg(c, 200, "发送评论成功")
}

// CommentToOtherUser
// @Tags Comment API
// @summary 用户评论其他用户
// @Accept json
// @Produce json
// @Param VideoID query string true "用户要评论的视频ID"
// @Param UserID query string true "用户ID"
// @Param UserID query string true "用户要评论的评论ID"
// @Param CommentContent formData string true "用户要评论的内容"
// @Router /comment/ToUser [post]
func CommentToOtherUser(c *gin.Context) {

}
