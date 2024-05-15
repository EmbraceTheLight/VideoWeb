package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateFavorites
// @Tags User API
// @summary 创建收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param FName formData string true "收藏夹名称"
// @Param IsPrivate query string true "是否私密"  Enums(公开, 私密)
// @Param Description formData string false "描述"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /User/Favorites/Create [post]
func CreateFavorites(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			Utilities.SendErrMsg(c, "service::Favorites::CreateFavorites", define.CreateFavoriteFailed, "创建收藏夹失败:"+err.Error())
		}
	}()
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	IsPrivate := logic.String2Int8(c.Query("IsPrivate"))
	FavoriteID := logic.GetUUID()
	FName := c.PostForm("FName")
	Description := c.PostForm("Description")

	Favorite := &EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  FavoriteID,
		UID:         UID,
		FName:       FName,
		Description: Description,
		IsPrivate:   IsPrivate,
		Videos:      nil,
	}

	/*判断是否有同名收藏夹*/
	f, err := EntitySets.GetFavoriteRecordByNameUserID(DAO.DB, FName, UID)
	if f != nil { //有同名收藏夹
		err = errors.New("已有同名收藏夹")
		return
	}

	/*创建收藏夹*/
	err = EntitySets.InsertFavoritesRecords(DAO.DB, Favorite)
	if err != nil {
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建收藏夹成功")
}

// DeleteFavorites
// @Tags User API
// @summary 删除收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param FName query string true "收藏夹名称"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /User/Favorites/Delete [delete]
func DeleteFavorites(c *gin.Context) {
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	FName := c.Query("FName")

	err := logic.DeleteFavoritesRecordsByNameUserID(FName, UID)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::DeleteFavorites", define.DeleteFavoriteFailed, "删除收藏夹失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "删除收藏夹"+FName+"成功")
}

// ModifyFavorites
// @Tags User API
// @summary 修改收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param Authorization header string true "token"
// @Param FavoriteID query string true "收藏夹ID"
// @Param newName formData string false "要修改的收藏夹名称"
// @Param IsPrivate formData string false "是否私密"  Enums(公开, 私密)
// @Param Description formData string false "描述"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /User/Favorites/Modify [put]
func ModifyFavorites(c *gin.Context) {
	FavoriteID := Utilities.String2Int64(c.Query("FavoriteID"))
	IsPrivate := logic.String2Int8(c.PostForm("IsPrivate"))
	newName := c.PostForm("newName")
	Description := c.PostForm("Description")

	/*更新收藏夹信息*/
	oldFavorite, _ := EntitySets.GetFavoriteRecordByFavoriteID(DAO.DB, FavoriteID)
	var newFavorite = *oldFavorite
	if newName != "" { //更新收藏夹名称
		newFavorite.FName = newName
	}
	if IsPrivate != 0 {
		newFavorite.IsPrivate = IsPrivate //更新收藏夹私密状态
	}
	newFavorite.Description = Description //更新收藏夹描述
	err := EntitySets.SaveFavoritesRecords(DAO.DB, &newFavorite)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::ModifyFavorites", define.ModifyFavoriteFailed, "修改收藏夹失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "修改收藏夹成功")
}
