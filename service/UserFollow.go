package service

import (
	"VideoWeb/DAO"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FollowOtherUser
// @Tags User API
// @summary 关注其他用户
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param FID query string true "要关注的用户ID"
// @Param FollowListID query string true "关注列表ID"
// @Router /User/Fans/Follows [post]
func FollowOtherUser(c *gin.Context) {
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	FID := Utilities.String2Int64(c.Query("FID"))
	FollowListID := Utilities.String2Int64(c.Query("FollowListID"))
	err := logic.FollowOtherUser(c, FollowListID, UID, FID)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "关注成功")
}

// UnFollowOtherUser
// @Tags User API
// @summary 取消关注其他用户
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param FID query string true "要取消关注的用户ID"
// @Param FollowListID query string true "要取消关注的用户ID"
// @Router /User/Fans/Unfollows [delete]
func UnFollowOtherUser(c *gin.Context) {
	u, _ := c.Get("user")
	UID := u.(*logic.UserClaims).UserId
	FID := Utilities.String2Int64(c.Query("FID"))
	FollowListID := Utilities.String2Int64(c.Query("FollowListID"))

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*更新用户的关注列表*/
	followsRecord := &RelationshipSets.UserFollows{
		FollowListID: FollowListID,
		UID:          UID,
		FID:          FID,
	}
	err := RelationshipSets.DeleteFollowsRecord(tx, followsRecord)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::UserFollow::UnFollowOtherUser", define.UnfollowUserFailed, "取消关注用户失败"+err.Error())
		return
	}

	/*更新被关注用户的被关注（粉丝）列表*/
	followedRecord := &RelationshipSets.UserFollowed{
		UID: FID,
		FID: UID,
	}
	err = RelationshipSets.DeleteFollowedRecord(tx, followedRecord)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::UserFollow::UnFollowOtherUser", define.UnfollowUserFailed, "取消关注用户失败"+err.Error())
		return
	}
	tx.Commit()

	Utilities.SendJsonMsg(c, http.StatusOK, "取消关注成功")
}
