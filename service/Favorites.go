package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateFavorites
// @Tags User API
// @summary 创建收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// //@Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param FName formData string true "收藏夹名称"
// @Param IsPrivate query string true "是否私密"  Enums(公开, 私密)
// @Param Description formData string false "描述"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /user/{UserID}/favorites/create [post]
func CreateFavorites(c *gin.Context) {
	UID := c.Param("UserID")
	IsPrivate := logic.String2Int8(c.Query("IsPrivate"))
	FavoriteID := logic.GetUUID()
	FName := c.PostForm("FName")
	Description := c.PostForm("Description")

	Favorite := &EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  FavoriteID,
		UID:         Utilities.String2Int64(UID),
		FName:       FName,
		Description: Description,
		IsPrivate:   IsPrivate,
		Videos:      nil,
	}
	var count int64
	err := DAO.DB.Model(&EntitySets.Favorites{}).Where("FName", FName).Count(&count).Error
	if count != 0 {
		Utilities.SendErrMsg(c, "service::Favorites::CreateFavorites", define.SameNameFavorite, "已有同名收藏夹")
		return
	}
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::CreateFavorites", define.CreateFavoriteFailed, "创建收藏夹失败:"+err.Error())
		return
	}
	err = EntitySets.InsertFavoritesRecords(DAO.DB, Favorite)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::CreateFavorites", define.CreateFavoriteFailed, "创建收藏夹失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建收藏夹成功")
}

// DeleteFavorites
// @Tags User API
// @summary 删除收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param UserID path string true "用户ID"
// @Param FName query string true "收藏夹名称"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /user/{UserID}/favorites/delete [delete]
func DeleteFavorites(c *gin.Context) {
	UID := c.Param("UserID")
	FName := c.Query("FName")
	var del *EntitySets.Favorites
	err := DAO.DB.Model(&EntitySets.Favorites{}).Debug().Where("UID=? AND FName=?", UID, FName).Delete(&del).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::DeleteFavorites", define.DeleteFavoriteFailed, "删除收藏夹失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, 200, "删除收藏夹"+FName+"成功")
}

// ModifyFavorites
// @Tags User API
// @summary 修改收藏夹
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param UserID path string true "用户ID"
// @Param FavoriteID query string true "收藏夹ID"
// @Param newName formData string false "要修改的收藏夹名称"
// @Param IsPrivate formData string false "是否私密"  Enums(公开, 私密)
// @Param Description formData string false "描述"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /user/{UserID}/favorites/modify [put]
func ModifyFavorites(c *gin.Context) {
	UID := c.Param("UserID")
	FavoriteID := c.Query("FavoriteID")
	IsPrivate := logic.String2Int8(c.PostForm("IsPrivate"))

	newName := c.PostForm("newName")
	Description := c.PostForm("Description")

	var newFavorite EntitySets.Favorites
	var count int64
	_ = DAO.DB.Debug().Where("FavoriteID=?", FavoriteID).First(&newFavorite).Error

	err := DAO.DB.Debug().Model(&EntitySets.Favorites{}).Where("UID=? AND FName=?", UID, newName).Count(&count).Error

	if newName == "" {
		Utilities.SendErrMsg(c, "service::Favorites::ModifyFavorites", define.ProhibitFavoritesNameEmpty, "收藏夹名称不能为空")
		return
	}

	newFavorite.FName = newName
	if IsPrivate != 0 {
		newFavorite.IsPrivate = IsPrivate
	}

	newFavorite.Description = Description
	// 更新收藏夹,使用结构体更新记录，这样不会更新零值。如：不更新IsPrivate时，其字段值为0
	if Description == "" {
		err = DAO.DB.Model(EntitySets.Favorites{}).Where("FavoriteID=?", FavoriteID).Updates(&newFavorite).Update("Description", "").Error
	} else {
		err = DAO.DB.Model(EntitySets.Favorites{}).Where("FavoriteID=?", FavoriteID).Updates(&newFavorite).Error
	}
	if err != nil {
		Utilities.SendErrMsg(c, "service::Favorites::ModifyFavorites", define.ModifyFavoriteFailed, "修改收藏夹失败："+err.Error())
		return
	}

	Utilities.SendJsonMsg(c, 200, "修改收藏夹成功")
}
