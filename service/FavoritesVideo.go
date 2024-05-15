package service

import (
	"VideoWeb/Utilities"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BookmarkVideo
// @Tags Video API
// @summary 收藏视频
// @Accept json
// @Produce json
// @Param ID path string true "视频ID"
// @Param Authorization header string true "token"
// @Param FID query string true "收藏夹ID"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /video/{ID}/Bookmark [post]
func BookmarkVideo(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Video::BookmarkVideo")
	VID := Utilities.String2Int64(c.Param("ID"))
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	FID := Utilities.String2Int64(c.Query("FID"))

	err := logic.UpdateVideoFavorite(c, VID, FID, UID, 1)
	if err != nil {
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "收藏成功")
}
