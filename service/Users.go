package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"unicode/utf8"
)

// Register
// @Tags User API
// @summary 用户注册
// @Accept multipart/form-data
// @Produce json,multipart/form-data
// @Param userName formData string true "用户名"
// @Param password formData string true "用户密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Param Email formData string true "用户邮箱"
// @Param Code formData string true "验证码"
// @Param Signature formData string false "用户个性签名(至多25个字)"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /User/Register [post]
func Register(c *gin.Context) {
	userName := c.PostForm("userName")
	password := c.PostForm("password")
	repeatPassword := c.PostForm("repeatPassword")
	email := c.PostForm("Email")
	Signature := c.PostForm("Signature")
	verify := c.PostForm("Code")
	UUID := logic.GetUUID()

	newUser := EntitySets.User{
		MyModel:   define.MyModel{},
		UserID:    UUID,
		UserName:  userName,
		Password:  password,
		Email:     email,
		Signature: Signature,
	}
	err := logic.CheckRegisterInfo(&newUser, repeatPassword)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register::CheckRegisterInfoFailed", define.CheckRegisterInfoFailed, err.Error())
		return
	}

	//验证码获取及验证
	code, err := DAO.RDB.Get(c, email).Result()
	if errors.Is(err, redis.Nil) {
		Utilities.SendErrMsg(c, "service::Users::Register::RedisGet", define.CodeExpired, "验证码已过期，请重新获取验证码")
		return
	}
	if code != verify {
		Utilities.SendErrMsg(c, "service::Users::Register", define.VerificationError, "验证码输入错误，请重新输入")
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register::GenerateFromPassword", define.PasswordEncryptionError, "密码加密错误")
		return
	}
	newUser.Password = string(hashedPassword)

	//设置用户默认头像
	avatar, err := Utilities.ReadFileContent(define.PictureSavePath + "default.jpg")
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register-->Utilities.ReadFileContent", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}
	newUser.Avatar = avatar

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
			tx.Rollback()
		}
	}()

	defaultFavorites := EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  logic.GetUUID(),
		UID:         newUser.UserID,
		FName:       "默认收藏夹",
		Description: "",
		IsPrivate:   1,
		Videos:      nil,
	}
	privateFavorites := EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  logic.GetUUID(),
		UID:         newUser.UserID,
		FName:       "私密收藏夹",
		Description: "",
		IsPrivate:   -1,
		Videos:      nil,
	}
	userLevel := EntitySets.Level{
		UID: newUser.UserID,
	}

	err = EntitySets.InsertUserRecord(tx, &newUser)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}

	err = EntitySets.InsertFavoritesRecords(tx, &defaultFavorites)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}

	err = EntitySets.InsertFavoritesRecords(tx, &privateFavorites)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}

	err = EntitySets.InsertLevelRecords(tx, &userLevel)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}

	token, err := logic.CreateToken(newUser.UserID, newUser.UserName, newUser.IsAdmin)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"msg":   "注册成功",
		"token": token,
	})
	DAO.RDB.Del(c, email)
}

// Login
// @Tags User API
// @summary 用户登录
// @Accept multipart/form-data
// @Produce json,xml
// @Param Username formData string true "用户名"
// @Param password formData string true "用户密码"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /User/Login [post]
func Login(c *gin.Context) {
	var userInfo = new(EntitySets.User)
	var err error
	Username := c.PostForm("Username")
	password := c.PostForm("password")
	fmt.Println("Account:", Username)
	fmt.Println("password:", password)
	if Username == "" || password == "" {
		Utilities.SendErrMsg(c, "service::Users::Login", define.EmptyAccountOrPassword, "账号或密码为空")
		return
	}

	err = DAO.DB.Where("userName=?", Username).First(&userInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound != nil { //未找到用户信息记录
			Utilities.SendErrMsg(c, "service::Users::Login", define.AccountNotFind, "用户名不存在，请重新检查输入的账号")
			return
		}
		//其他错误
		Utilities.SendErrMsg(c, "service::Users::Login", define.ObtainUserInformationFailed, "Get UserInfo failed:"+err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
		Utilities.SendErrMsg(c, "service::Users::Login", define.ErrorPassword, "密码错误")
		return
	}
	token, err := logic.CreateToken(userInfo.UserID, userInfo.UserName, userInfo.IsAdmin)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Login", define.CreateTokenError, "CreateToken error:"+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// Logout
// @Tags User API
// @summary 用户注销
// @Accept json
// @Produce json
// @Param UserID path string true "用户ID"
// @Success 200 {string}  json "{"code":"200","msg":"注销用户成功"}"
// @Router /user/{UserID}/Logout [delete]
func Logout(c *gin.Context) {
	id := c.Param("UserID")
	favorites, err := logic.GetFavoritesByID(id)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Login", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		return
	}
	/*删除用户上传的视频信息*/
	err = os.RemoveAll(path.Join(define.VideoSavePath, id))
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		return
	}

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	/*删除用户的收藏夹信息*/
	for _, favorite := range favorites {
		err := EntitySets.DeleteFavoritesRecords(tx, favorite)
		if err != nil {
			Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
			tx.Rollback()
			return
		}
	}
	/*删除用户上传的视频信息*/
	err = DAO.DB.Where("UID=?", id).Delete(&EntitySets.Video{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	/*删除用户的关注列表信息*/
	err = DAO.DB.Where("UID=?", id).Delete(&RelationshipSets.UserFollows{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	/*删除被关注用户的对应粉丝列表信息*/
	err = DAO.DB.Where("FID=?", id).Delete(&RelationshipSets.UserFollowed{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	/*删除用户对应等级信息*/
	err = DAO.DB.Where("UID=?", id).Delete(&EntitySets.Level{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	/*删除用户信息*/
	var user *EntitySets.User
	err = DAO.DB.Debug().Where("UserID=?", id).Delete(&user).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	Utilities.SendJsonMsg(c, http.StatusOK, "注销账户成功")

}

// GetUserDetail
// @Tags User API
// @Summary 获取用户完整、详细的信息
// @Accept json
// @Produce json
// @Param UserID path string true "用户标识"
// @Success 200 {string}  json "{"code":"200","data":userInfo}"
// @Router /User/{UserID}/User-detail [get]
func GetUserDetail(c *gin.Context) {
	userID := c.Param("UserID")
	if userID == "" {
		Utilities.SendErrMsg(c, "service::Users::GetUserDetail", define.ObtainUserInformationFailed, "用户唯一标识不能为空")
		return
	}

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Omit("password").Where("UserID=?", userID).Preload("Videos").Preload("Favorites").
		Preload("Favorites.Videos").Preload("Comments").
		Preload("Follows").Preload("UserLevel").
		Preload("Followed").Preload("UserWatch").Preload("UserSearch").First(&userInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::GetUserDetail", define.ObtainUserInformationFailed, "Get User Info failed:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": userInfo,
	})
}

// ModifyUserSignature
// @Tags User API
// @summary 用户信息修改-更新用户签名
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param userSignature formData string false "用户签名,为空表示没有签名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户签名成功"}"
// @Router /User/{UserID}/ModifySignature [put]
func ModifyUserSignature(c *gin.Context) {
	id := c.Param("UserID")
	signature := c.PostForm("userSignature")
	if utf8.RuneCountInString(signature) > 25 {
		Utilities.SendErrMsg(c, "service::Users::ModifySignature", define.SignatureTooLong, "个性签名过长，请重新输入")
		return
	}

	err := DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("signature", signature).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifySignature", define.ModifySignatureFailed, "修改用户签名失败:"+err.Error())
		return
	}

	Utilities.SendJsonMsg(c, http.StatusOK, "修改用户签名成功")
}

// ModifyUserEmail
// @Tags User API
// @summary 用户信息修改-更新用户邮箱
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param userEmail formData string true "用户新邮箱"
// @Param Code formData string true "验证码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户邮箱成功"}"
// @Router /User/{UserID}/ModifyEmail [put]
func ModifyUserEmail(c *gin.Context) {
	//获取用户id,email,验证码
	id := c.Param("UserID")
	userEmail := c.PostForm("userEmail")
	verify := c.PostForm("Code")

	code, err := DAO.RDB.Get(c, userEmail).Result()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyEmail", -1, "验证码已过期，请重新获取验证码")
		return
	}
	if code != verify {
		Utilities.SendErrMsg(c, "service::Users::ModifyEmail", define.VerificationError, "验证码输入错误，请重新输入")
		return
	}

	//修改后存入数据库中
	err = DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("Email", userEmail).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyEmail", define.CodeSendFailed, "验证码发送失败"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "修改用户邮箱成功")
}

// ForgetPassword
// @Tags User API
// @summary 重置用户密码
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param userEmail formData string true "用户邮箱"
// @Param Code formData string true "验证码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"重置用户密码成功"}"
// @Router /User/{UserID}/ForgetPassword [put]
func ForgetPassword(c *gin.Context) {
	id := c.Param("UserID")
	userEmail := c.PostForm("userEmail")
	Code := c.PostForm("Code")
	newPassword := c.PostForm("newPassword")
	repeatPassword := c.PostForm("repeatPassword")

	var userInfo = new(EntitySets.User)
	verify, RDBErr := DAO.RDB.Get(c, userEmail).Result()
	err := DAO.DB.Model(&EntitySets.User{}).Where("Userid=?", id).First(&userInfo).Error
	switch {
	case err != nil:
		Utilities.SendErrMsg(c, "service::Users::ForgetPassword", define.ObtainUserInformationFailed, "获取用户信息失败")
		return
	case userEmail != userInfo.Email:
		Utilities.SendErrMsg(c, "service::Users::ForgetPassword", define.NotMatchMail, "邮箱不匹配")
		return
	case RDBErr != nil:
		Utilities.SendErrMsg(c, "service::Users::ForgetPassword", define.CodeExpired, "验证码已过期，请重新获取验证码")
		return
	case Code != verify:
		Utilities.SendErrMsg(c, "service::Users::ForgetPassword", define.VerificationError, "验证码输入错误，请重新输入")
		return
	}

	status, err := logic.ModifyPassword(id, newPassword, repeatPassword)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ForgetPassword", status, err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "重置密码成功")
	DAO.RDB.Del(c, userEmail)
}

// ModifyPassword
// @Tags User API
// @summary 用户信息修改-修改用户密码
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param password formData string true "用户密码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户密码成功"}"
// @Router /User/{UserID}/ModifyPassword [put]
func ModifyPassword(c *gin.Context) {
	id := c.Param("UserID")
	password := c.PostForm("password")
	newPassword := c.PostForm("newPassword")
	repeatPassword := c.PostForm("repeatPassword")

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Where("Userid=?", id).First(&userInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyPassword", define.ObtainUserInformationFailed, "获取用户信息失败")
		return
	}

	err = logic.ComparePassword(userInfo.Password, password)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyPassword", define.ErrorPassword, err.Error())
		return
	}

	code, err := logic.ModifyPassword(id, newPassword, repeatPassword)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyPassword", code, err.Error())
		return
	}
	Utilities.SendJsonMsg(c, code, "修改密码成功")
}

// ModifyUserName
// @Tags User API
// @summary 用户信息修改-修改用户名
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param userName formData string true "用户名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户名成功"}"
// @Router /User/{UserID}/ModifyUserName [put]
func ModifyUserName(c *gin.Context) {
	id := c.Param("UserID")
	userName := c.PostForm("userName")

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Where("Userid=?", id).First(&userInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ObtainUserInformationFailed, "获取用户信息失败")
		return
	}

	var countName int64
	err = DAO.DB.Model(&EntitySets.User{}).Where("userName=?", userName).Count(&countName).Error
	switch {
	case err != nil: //数据库查询错误
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.QueryUserError, "error while searching database:"+err.Error())
		return
	case countName > 0: //已有同名用户
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ExistUserName, "用户名已存在，请重新输入用户名")
		return
	}

	err = DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("userName", userName).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ModifyUserNameFailed, "修改用户名失败")
		return
	}

	Utilities.SendJsonMsg(c, 200, "修改用户名成功")
	DAO.RDB.Del(c, userInfo.Email)
}

// UploadUserAvatar
// @Tags User API
// @summary 用户信息修改-上传头像
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param UserID path string true "用户ID"
// @Param file formData file true "头像"
// @Success 200 {string}  json "{"code":"200","msg":"上传头像成功"}"
// @Router /User/{UserID}/Face/Upload/Avatar [put]
func UploadUserAvatar(c *gin.Context) {
	userID := c.Param("UserID")
	FH, _ := c.FormFile("file") //FH=FileHeader
	//TODO:检查文件后缀名是否为 .jpg/.jpeg/.png/.jfif
	fname := FH.Filename
	extension := path.Ext(fname)

	println("ext:", extension)
	if _, ok := define.PicExtCheck[extension]; ok != true {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.ImageFormatError, "图片格式错误或不支持此图片格式")
		return
	}

	file, err := FH.Open()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.OpenFileFailed, "打开文件失败"+err.Error())
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.ReadFileFailed, "读取文件内容失败"+err.Error())
		return
	}

	err = DAO.DB.Model(&EntitySets.User{}).Where("userID=?", userID).Update("Avatar", data).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.UploadUserAvatarFailed, "上传用户头像失败")
		fmt.Println("err:", err.Error())
		return
	}

	Utilities.SendJsonMsg(c, 200, "上传用户头像成功")
}

// AddSearchHistory
// @Tags User API
// @summary 增加搜索历史记录
// @Accept json
// @Produce json
// @Param UserID path string true "用户ID"
// @Param searchString query string true "搜索记录"
// @Router /User/{UserID}/AddSearchHistory [post]
func AddSearchHistory(c *gin.Context) {
	UID := c.Param("UserID")
	searchString := c.Query("searchString")
	SearchHistory := &EntitySets.SearchHistory{
		Model:        gorm.Model{},
		UID:          Utilities.String2Int64(UID),
		SearchString: searchString,
	}
	err := EntitySets.InsertSearchRecord(DAO.DB, SearchHistory)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::AddSearchHistory", 5019, "创建搜索历史失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建搜索历史成功")

}

// AddVideoHistory
// @Tags User API
// @summary 增加视频历史记录
// @Accept json
// @Produce json
// @Param UserID path string true "用户ID"
// @Param VID query string true "视频ID"
// @Router /User/{UserID}/AddVideoHistory [post]
func AddVideoHistory(c *gin.Context) {
	UID := c.Param("UserID")
	VID := c.Query("VID")
	VideoHistory := &EntitySets.VideoHistory{
		MyModel: define.MyModel{},
		UID:     Utilities.String2Int64(UID),
		VID:     Utilities.String2Int64(VID),
	}
	err := EntitySets.InsertVideoHistoryRecord(DAO.DB, VideoHistory)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::AddVideoHistory", 5020, "创建视频历史记录失败")
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建视频历史记录成功")
}
