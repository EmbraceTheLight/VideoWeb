package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
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
	"net/http"
	"path"
	"unicode/utf8"
)

// Register
// @Tags User API
// @summary 用户注册
// @Accept multipart/form-data
// @Produce json
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

	newUser := &EntitySets.User{
		MyModel:   define.MyModel{},
		UserID:    UUID,
		UserName:  userName,
		Password:  password,
		Email:     email,
		Signature: Signature,
	}
	err := logic.CheckRegisterInfo(newUser, repeatPassword)
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

	defaultFavorites := &EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  logic.GetUUID(),
		UID:         newUser.UserID,
		FName:       "默认收藏夹",
		Description: "",
		IsPrivate:   1,
		Videos:      nil,
	}
	privateFavorites := &EntitySets.Favorites{
		MyModel:     define.MyModel{},
		FavoriteID:  logic.GetUUID(),
		UID:         newUser.UserID,
		FName:       "私密收藏夹",
		Description: "",
		IsPrivate:   -1,
		Videos:      nil,
	}
	userLevel := &EntitySets.Level{
		UID: newUser.UserID,
	}
	defaultFollowList := &EntitySets.FollowList{
		ListID:   logic.GetUUID(),
		UID:      newUser.UserID,
		ListName: "默认关注列表",
	}

	err = logic.InsertInitRecords(defaultFavorites, privateFavorites, userLevel, defaultFollowList, newUser)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register-->logic.InsertInitRecords", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}

	token, err := logic.CreateToken(newUser.UserID, newUser.UserName, newUser.IsAdmin)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}

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
	var err error
	Username := c.PostForm("Username")
	password := c.PostForm("password")
	fmt.Println("Account:", Username)
	fmt.Println("password:", password)
	if Username == "" || password == "" {
		Utilities.SendErrMsg(c, "service::Users::Login", define.EmptyAccountOrPassword, "账号或密码为空")
		return
	}

	userInfo, err := EntitySets.GetUserInfoByName(DAO.DB, Username)

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
// @Param Authorization header string true "token"
// @Success 200 {string}  json "{"code":"200","msg":"注销用户成功"}"
// @Router /User/Logout [delete]
func Logout(c *gin.Context) {
	Utilities.AddFuncName(c, "service::Users::Logout")
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
	///*删除用户磁盘资源，如视频、图片、音频等*/
	//err := logic.RemoveUserResource(strconv.FormatInt(id, 10))
	//if err != nil {
	//	Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
	//	return
	//}

	/*删除用户在数据库中的信息，如用户记录,关注列表,收藏夹,关注用户,搜索记录,观看记录,收藏夹收藏的信息等*/
	err := logic.DeleteUserInfoInDB(c, id)
	if err != nil {
		Utilities.HandleInternalServerError(c, err)
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "注销账户成功")

}

// GetUserDetail
// @Tags User API
// @Summary 获取用户完整、详细的信息
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {string}  json "{"code":"200","data":userInfo}"
// @Router /User/User-detail [get]
func GetUserDetail(c *gin.Context) {
	u, _ := c.Get("user")
	userID := logic.GetUserID(u)

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Omit("password").Where("user_id=?", userID).Preload("Videos").Preload("Favorites").
		Preload("Favorites.Videos").Preload("Comments").
		Preload("Follows").Preload("Follows.Users").Preload("UserLevel").
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
// @Param Authorization header string true "token"
// @Param userSignature formData string false "用户签名,为空表示没有签名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户签名成功"}"
// @Router /User/ModifySignature [put]
func ModifyUserSignature(c *gin.Context) {
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
	signature := c.PostForm("userSignature")
	if utf8.RuneCountInString(signature) > 25 {
		Utilities.SendErrMsg(c, "service::Users::ModifySignature", define.SignatureTooLong, "个性签名过长，请重新输入")
		return
	}

	err := EntitySets.UpdateUserStringField(DAO.DB, id, "signature", signature)
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
// @Param Authorization header string true "token"
// @Param userEmail formData string true "用户新邮箱"
// @Param Code formData string true "验证码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户邮箱成功"}"
// @Router /User/ModifyEmail [put]
func ModifyUserEmail(c *gin.Context) {
	//获取用户id,email,验证码
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
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
	err = EntitySets.UpdateUserStringField(DAO.DB, id, "email", userEmail)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyEmail", define.CodeSendFailed, "邮箱修改失败:"+err.Error())
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "修改用户邮箱成功")
	DAO.RDB.Del(c, userEmail)
}

// ForgetPassword
// @Tags User API
// @summary 重置用户密码
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param userEmail formData string true "用户邮箱"
// @Param Code formData string true "验证码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"重置用户密码成功"}"
// @Router /User/ForgetPassword [put]
func ForgetPassword(c *gin.Context) {
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
	userEmail := c.PostForm("userEmail")
	Code := c.PostForm("Code")
	newPassword := c.PostForm("newPassword")
	repeatPassword := c.PostForm("repeatPassword")

	verify, RDBErr := DAO.RDB.Get(c, userEmail).Result()
	userInfo, err := EntitySets.GetUserInfoByID(DAO.DB, id)
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
// @Param Authorization header string true "token"
// @Param password formData string true "用户密码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户密码成功"}"
// @Router /User/ModifyPassword [put]
func ModifyPassword(c *gin.Context) {
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
	password := c.PostForm("password")
	newPassword := c.PostForm("newPassword")
	repeatPassword := c.PostForm("repeatPassword")

	userInfo, err := EntitySets.GetUserInfoByID(DAO.DB, id)
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
// @Param Authorization header string true "token"
// @Param userName formData string true "用户名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户名成功"}"
// @Router /User/ModifyUserName [put]
func ModifyUserName(c *gin.Context) {
	u, _ := c.Get("user")
	id := logic.GetUserID(u)
	userName := c.PostForm("userName")

	userInfo, err := EntitySets.GetUserInfoByID(DAO.DB, id)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ObtainUserInformationFailed, "获取用户信息失败")
		return
	}

	user, _ := EntitySets.GetUserInfoByName(DAO.DB, userName)

	if user != nil { //已有同名用户
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ExistUserName, "用户名已存在，请重新输入用户名")
		return
	}

	err = EntitySets.UpdateUserStringField(DAO.DB, id, "user_name", userName)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::ModifyUserName", define.ModifyUserNameFailed, "修改用户名失败")
		return
	}

	Utilities.SendJsonMsg(c, http.StatusOK, "修改用户名成功")
	DAO.RDB.Del(c, userInfo.Email)
}

// UploadUserAvatar
// @Tags User API
// @summary 用户信息修改-上传头像
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param file formData file true "头像"
// @Success 200 {string}  json "{"code":"200","msg":"上传头像成功"}"
// @Router /User/Face/Upload/Avatar [put]
func UploadUserAvatar(c *gin.Context) {
	u, _ := c.Get("user")
	userID := logic.GetUserID(u)
	FH, _ := c.FormFile("file") //FH=FileHeader

	//检查文件后缀名是否为 .jpg/.jpeg/.png/.jfif
	fname := FH.Filename
	extension := path.Ext(fname)
	if _, ok := define.PicExtCheck[extension]; ok != true {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.ImageFormatError, "图片格式错误或不支持此图片格式")
		return
	}

	//读取用户头像文件内容
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

	//更新用户头像
	err = EntitySets.UpdateUserAvatar(DAO.DB, userID, data)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.UploadUserAvatarFailed, "上传用户头像失败")
		fmt.Println("err:", err.Error())
		return
	}

	Utilities.SendJsonMsg(c, http.StatusOK, "上传用户头像成功")
}

// AddSearchHistory
// @Tags User API
// @summary 增加搜索历史记录
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param searchString query string true "搜索记录"
// @Router /User/AddSearchHistory [post]
func AddSearchHistory(c *gin.Context) {
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	searchString := c.Query("searchString")
	SearchHistory := &EntitySets.SearchHistory{
		Model:        gorm.Model{},
		UID:          UID,
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
// @Param Authorization header string true "token"
// @Param VID query string true "视频ID"
// @Router /User/AddVideoHistory [post]
func AddVideoHistory(c *gin.Context) {
	u, _ := c.Get("user")
	UID := logic.GetUserID(u)
	VID := Utilities.String2Int64(c.Query("VID"))
	VideoHistory := &EntitySets.VideoHistory{
		MyModel: define.MyModel{},
		UID:     UID,
		VID:     VID,
	}
	err := EntitySets.InsertVideoHistoryRecord(DAO.DB, VideoHistory)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::AddVideoHistory", 5020, "创建视频历史记录失败")
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建视频历史记录成功")
}
