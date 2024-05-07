package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// AddBarrage
// @Summary 添加弹幕
// @Description 添加弹幕
// @Tags Barrage API
// @Accept  multipart/form-data
// @Produce  json
// @param UserID query string true "用户ID"
// @param Color query string true "弹幕颜色"
// @param Time query string true "发送弹幕时的视频时间"
// @Param Content formData string true "弹幕数据"
// @Success 200 {string}  json "{"code":"200","msg":"添加弹幕成功"}"
// @Router /video/{VideoID}/AddBarrage [post]
func AddBarrage(c *gin.Context) {
	VID := c.Param("VideoID")
	UID := c.Query("UserID")
	color := c.Query("Color")
	seconds := c.Query("Time")
	content := c.PostForm("Content")

	second, _ := strconv.Atoi(seconds)
	_, t := Utilities.SecondToTime(int64(second))
	barrage := &EntitySets.Barrage{
		Model:   gorm.Model{},
		UID:     UID,
		VID:     VID,
		Content: content,
		Hour:    uint8(t[0]),
		Minute:  uint8(t[1]),
		Second:  uint8(t[2]),
		Color:   color,
	}
	err := EntitySets.InsertBarrageRecord(DAO.DB, barrage)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Barrage::AddBarrage", define.AddBarrageFailed, "添加弹幕失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, 200, "添加弹幕成功")
}
