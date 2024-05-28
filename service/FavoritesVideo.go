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
	UID := logic.GetUserID(u)
	FID := Utilities.String2Int64(c.Query("FID"))

	err := logic.UpdateVideoFavorite(c, VID, FID, UID, 1)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "收藏成功")
}

// UnBookmarkVideo
// @Tags Video API
// @summary 取消收藏视频
// @Accept json
// @Produce json
// @Param ID path string true "视频ID"
// @Param Authorization header string true "token"
// @Success 200 {string}  json "{"code":"200","msg":"取消收藏成功"}"
// @Router /video/{ID}/Bookmark [delete]
func UnBookmarkVideo(c *gin.Context) {
	Utilities.AddFuncName(c, "Service::Video::BookmarkVideo")
	VID := Utilities.String2Int64(c.Param("ID"))
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	var FID int64 = 0 // 取消收藏时，用不到收藏夹，其ID置为0
	err := logic.UpdateVideoFavorite(c, VID, FID, UID, -1)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "取消收藏成功")
}
