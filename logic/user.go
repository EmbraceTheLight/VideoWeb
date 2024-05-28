package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"path"
	"unicode/utf8"
)

// GetUserID 获取用户ID
func GetUserID(u any) int64 {
	if u != nil {
		return u.(*UserClaims).UserId
	} else {
		return 0
	}
}

func InsertInitRecords(defaultFavorites, privateFavorites *EntitySets.Favorites,
	userLevel *EntitySets.Level,
	defaultFollowList *EntitySets.FollowList,
	newUser *EntitySets.User) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = EntitySets.InsertFollowList(tx, defaultFollowList)
	if err != nil {
		return err
	}

	err = EntitySets.InsertUserRecord(tx, newUser)
	if err != nil {
		return err
	}

	err = EntitySets.InsertFavoritesRecords(tx, defaultFavorites)
	if err != nil {
		return err
	}

	err = EntitySets.InsertFavoritesRecords(tx, privateFavorites)
	if err != nil {
		return err
	}

	err = EntitySets.InsertLevelRecords(tx, userLevel)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// GetUserIpInfo 获取并返回用户所在国家和地区
func GetUserIpInfo(c *gin.Context) (Country, City string) {
	UserIP := c.ClientIP()
	fmt.Println(UserIP)
	info, err := Utilities.GetIPInfo(UserIP)
	if err != nil {
		fmt.Println("获取用户IP失败:", err)
		return "", ""
	}
	return info.Country, info.City
}

// ComparePassword  比较用户输入的密码与数据库中的密码
func ComparePassword(userPassword, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword)); err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// ModifyPassword 用户修改辅助函数
func ModifyPassword(id int64, newPassword, repeatPassword string) (int, error) {
	switch {
	case len(newPassword) < 6: //密码长度小于6位
		return 4002, errors.New("密码长度不能小于6位，请重新输入密码")
	case newPassword != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		return 4003, errors.New("第一次输入的密码与第二次输入的密码不一致，请重新输入")
	}
	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return define.PasswordEncryptionError, errors.New("密码加密错误")
	}

	//更新用户密码
	err = EntitySets.UpdateUserStringField(DAO.DB, id, "password", string(hashedPassword))
	if err != nil {
		return define.ModifyPasswordFailed, errors.New("修改密码失败")
	}
	return http.StatusOK, nil
}

// CheckRegisterInfo 检查注册信息是否正确。
func CheckRegisterInfo(checkInfo *EntitySets.User, repeatPassword string) error {
	var countName int64
	err := DAO.DB.Model(&EntitySets.User{}).Where("user_name=?", checkInfo.UserName).Count(&countName).Error

	if err != nil {
		return errors.New("查询用户信息失败")
	}
	switch {
	case countName > 0: //已有同名用户
		return errors.New("用户名已存在")
	case len(checkInfo.Password) < 6: //密码长度小于6位
		return errors.New("密码长度小于6位")
	case checkInfo.Password != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		return errors.New("第一次输入的密码与第二次输入的密码不一致")
	case utf8.RuneCountInString(checkInfo.Signature) > 25:
		return errors.New("个性签名过长")
	}
	return nil
}

// RemoveUserResource 删除用户资源
func RemoveUserResource(userID string) error {
	err := os.RemoveAll(path.Join(define.VideoSavePath, userID))
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserInfoInDB 删除用户在数据库中的所有相关信息
func DeleteUserInfoInDB(c *gin.Context, uid int64) error {
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			Utilities.AddFuncName(c, "DeleteUserInfoInDB")
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	/*删除用户的收藏夹信息*/
	err = helper.DeleteUserFavoriteRecords(uid, tx)
	if err != nil {
		return err
	}
	///*删除用户上传的视频信息*/
	//err = EntitySets.DeleteVideoInfoByUserID(tx, uid)
	//if err != nil {
	//	return err
	//}
	/*删除用户的关注列表信息*/
	err = helper.DeleteFollowListRecords(uid, tx)
	if err != nil {
		return err
	}
	/*删除被关注用户的对应粉丝列表信息*/
	err = RelationshipSets.DeleteFollowedRecordsByUserID(tx, uid)
	if err != nil {
		return err
	}
	/*删除用户对应等级信息*/
	err = EntitySets.DeleteLevelRecordByUserID(tx, uid)
	if err != nil {
		return err
	}
	/*删除用户的搜索历史信息和观看历史信息*/
	err = EntitySets.DeleteAllSearchRecord(tx, uid)
	if err != nil {
		return err
	}
	err = EntitySets.DeleteAllVideoHistoryRecords(tx, uid)
	if err != nil {
		return err
	}
	/*删除用户信息*/
	err = EntitySets.DeleteUserRecordByID(tx, uid)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// FollowOtherUser 关注其他用户
func FollowOtherUser(c *gin.Context, followlistID, UID, followsID int64) error {
	Utilities.AddFuncName(c, "FollowOtherUser")
	var err error
	tx := DAO.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	/*更新关注用户的关注列表*/
	followsRecord := &RelationshipSets.UserFollows{
		FollowListID: followlistID,
		UID:          UID,
		FID:          followsID,
	}
	err = RelationshipSets.InsertFollowsRecord(tx, followsRecord)
	if err != nil {
		return err
	}

	/*更新被关注用户的被关注（粉丝）列表*/
	followedRecord := &RelationshipSets.UserFollowed{
		MyModel: define.MyModel{},
		UID:     followsID,
		FID:     UID,
	}

	err = RelationshipSets.InsertFollowedRecord(tx, followedRecord)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
