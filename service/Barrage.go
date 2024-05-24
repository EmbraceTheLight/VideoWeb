package service

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddBarrage
// @Summary 添加弹幕
// @Description 添加弹幕
// @Tags Barrage API
// @Accept  multipart/form-data
// @Produce  json
// @Param ID path string true "视频ID"
// @Param Authorization header string true "token"
// @param Color query string true "弹幕颜色"
// @param Time query string true "发送弹幕时的视频时间"
// @Param Content formData string true "弹幕数据"
// @Success 200 {string}  json "{"code":"200","msg":"添加弹幕成功"}"
// @Router /video/{ID}/AddBarrage [post]
func AddBarrage(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Barrage::AddBarrage")
	var err error
	VID := Utilities.String2Int64(c.Param("ID"))
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	color := c.Query("Color")
	seconds := c.Query("Time")
	content := c.PostForm("Content")

	second, _ := strconv.Atoi(seconds)
	_, t := Utilities.SecondToTime(int64(second))
	barrage := &EntitySets.Barrage{
		MyModel: define.MyModel{},
		UID:     UID,
		VID:     VID,
		Content: content,
		Hour:    uint8(t[0]),
		Minute:  uint8(t[1]),
		Second:  uint8(t[2]),
		Color:   color,
	}

	err = logic.AddVideoBarrage(c, VID, barrage)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "添加弹幕成功")
}
