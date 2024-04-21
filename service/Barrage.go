package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddBarrage
// @Summary 添加弹幕
// @Description 添加弹幕
// @Tags Barrage API
// @Accept  multipart/form-data
// @Produce  json
// @param UserID query string true "用户ID"
// @param Time query string true "发送弹幕时的视频时间"
// @Param Content formData string true "弹幕数据"
// @Success 200 {string}  json "{"code":"200","msg":"添加弹幕成功"}"
// @Router /video/{VideoID}/AddBarrage [post]
func AddBarrage(c *gin.Context) {
	VID := c.Param("VideoID")
	UID := c.Query("UserID")
	t := c.Query("Time")
	content := c.PostForm("Content")
	fmt.Println(content)
	minute, second := logic.ParseTime(t)
	barrage := &EntitySets.Barrage{
		Model:   gorm.Model{},
		UID:     UID,
		VID:     VID,
		Content: content,
		Minute:  minute,
		Second:  second,
	}
	err := barrage.Create(DAO.DB)
	if err != nil {
		Utilities.SendJsonMsg(c, define.AddBarrageFailed, "添加弹幕失败:"+err.Error())
		Utilities.WriteErrLog("AddBarrage", "添加弹幕失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, 200, "添加弹幕成功")
}
