package service

import (
	"VideoWeb/DAO"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FollowOtherUser
// @Tags User API
// @summary 关注其他用户
// @Accept json
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param FID query string true "要关注的用户ID"
// @Router /user/fans/follows [post]
func FollowOtherUser(c *gin.Context) {
	UID := c.Query("userID")
	FID := c.Query("FID")

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//TODO:更新关注用户的关注列表
	followsRecord := &RelationshipSets.UserFollows{
		Model:     gorm.Model{},
		GroupName: "默认关注分组",
		UID:       UID,
		FID:       FID,
	}
	err := RelationshipSets.InsertFollowsRecord(tx, followsRecord)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::UserFollow::FollowOtherUser", define.FollowUserFailed, "关注用户失败"+err.Error())
		return
	}

	//TODO:更新被关注用户的被关注（粉丝）列表
	followedRecord := &RelationshipSets.UserFollowed{
		MyModel: define.MyModel{},
		UID:     FID,
		FID:     UID,
	}

	err = RelationshipSets.InsertFollowedRecord(tx, followedRecord)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::UserFollow::FollowOtherUser", define.FollowUserFailed, "关注用户失败"+err.Error())
		return
	}
	tx.Commit()

	Utilities.SendJsonMsg(c, 200, "关注成功")
}

// UnFollowOtherUser
// @Tags User API
// @summary 取消关注其他用户
// @Accept json
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param FID query string true "要取消关注的用户ID"
// @Router /user/fans/unfollows [delete]
func UnFollowOtherUser(c *gin.Context) {
	UID := c.Query("userID")
	FID := c.Query("FID")

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//TODO:更新用户的关注列表
	followsRecord := &RelationshipSets.UserFollows{
		UID: UID,
		FID: FID,
	}
	err := RelationshipSets.DeleteFollowsRecord(tx, followsRecord)
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, "service::UserFollow::UnFollowOtherUser", define.UnfollowUserFailed, "取消关注用户失败"+err.Error())
		return
	}

	//TODO:更新被关注用户的被关注（粉丝）列表
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
	// TODO:
	Utilities.SendJsonMsg(c, 200, "取消关注成功")
}
