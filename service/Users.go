package service

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
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
// @Router /user/register [post]
func Register(c *gin.Context) {
	userName := c.PostForm("userName")
	password := c.PostForm("password")
	repeatPassword := c.PostForm("repeatPassword")
	email := c.PostForm("Email")
	Signature := c.PostForm("Signature")
	verify := c.PostForm("Code")
	UUID := logic.GetUUID()

	var countName int64
	//var countEmail int64
	err := DAO.DB.Model(&EntitySets.User{}).Where("userName=?", userName).Count(&countName).Error
	switch {
	case err != nil: //数据库查询错误
		Utilities.SendErrMsg(c, "service::Users::Register", define.QueryUserError, "error while searching database:"+err.Error())

		return
	case countName > 0: //已有同名用户
		Utilities.SendErrMsg(c, "service::Users::Register", define.ExistUserName, "用户名已存在，请重新输入待注册用户名")
		return
	case len(password) < 6: //密码长度小于6位
		Utilities.SendErrMsg(c, "service::Users::Register", define.ShortPasswordLength, "密码长度不能小于6位，请重新输入密码")
		return
	case password != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		Utilities.SendErrMsg(c, "service::Users::Register", define.PasswordInconsistency, "第一次输入的密码与第二次输入的密码不一致，请重新输入")
		return
	case utf8.RuneCountInString(Signature) > 25:
		Utilities.SendErrMsg(c, "service::Users::Register", define.SignatureTooLong, "个性签名过长，请重新输入")
		return
	}

	//验证码获取及验证
	//code := define.VerificationDataMap[email].Code
	code, err := DAO.RDB.Get(c, email).Result()
	if err != nil {
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

	//设置用户默认头像
	file, err := os.Open(define.PictureSavePath + "default.jpg")
	defer file.Close()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}
	var avatar = make([]byte, fileInfo.Size())
	_, err = file.Read(avatar)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		return
	}

	newUser := EntitySets.User{
		MyModel:   define.MyModel{},
		UserID:    UUID,
		UserName:  userName,
		Password:  string(hashedPassword),
		Email:     email,
		Signature: Signature,
		Avatar:    avatar,
	}
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
	Account, err := newUser.Create(tx)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = defaultFavorites.Create(tx)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = privateFavorites.Create(tx)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Register", define.CreateUserFailed, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = userLevel.Create(tx)
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
		"code":    http.StatusOK,
		"msg":     "注册成功",
		"Account": Account, //返回创建好的账号
		"token":   token,
	})
}

// SendCode
// @Tags User API
// @summary 发送验证码
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	userEmail := c.Query("email") //从前端获取email信息
	if userEmail == "" {
		Utilities.SendErrMsg(c, "service::Users::SendCode", define.EmptyMail, "邮箱不能为空")
		return
	}
	DAO.RDB.Del(c, userEmail) //之前的验证码可能没有过期，要先删除
	code := logic.CreateVerificationCode()
	err := logic.SendCode(userEmail, code)
	fmt.Println(code)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::SendCode", define.CodeSendFailed, "验证码发送失败:"+err.Error())
		return
	}

	DAO.RDB.Set(c, userEmail, code, time.Duration(define.Expired)*time.Second)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功！",
	})
}

// Login
// @Tags User API
// @summary 用户登录
// @Accept multipart/form-data
// @Produce json,xml
// @Param Account formData string true "用户账号"
// @Param password formData string true "用户密码"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /user/login [post]
func Login(c *gin.Context) {
	var userInfo = new(EntitySets.User)
	var err error
	AccountOrUserName := c.PostForm("Account")
	password := c.PostForm("password")
	fmt.Println("Account:", AccountOrUserName)
	fmt.Println("password:", password)
	if AccountOrUserName == "" || password == "" {
		Utilities.SendErrMsg(c, "service::Users::Login", define.EmptyAccountOrPassword, "账号或密码为空")
		return
	}

	err = DAO.DB.Where("Account=? OR userName=?", AccountOrUserName, AccountOrUserName).First(&userInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound != nil { //未找到用户信息记录
			Utilities.SendErrMsg(c, "service::Users::Login", define.AccountNotFind, "账号不存在，请重新检查输入的账号")
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
		"code": 200,
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
// @Param UserID query string true "用户ID"
// @Success 200 {string}  json "{"code":"200","msg":"注销用户成功"}"
// @Router /user/Logout [delete]
func Logout(c *gin.Context) {
	id := c.Query("UserID")
	favorites, err := logic.GetFavoritesByID(id)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Login", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		fmt.Println("Err in Getting User Favorites:", err.Error())
		return
	}

	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//TODO:删除用户的收藏夹信息
	for _, favorite := range favorites {
		err := favorite.Delete(DAO.DB)
		if err != nil {
			Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
			tx.Rollback()
			fmt.Println("Err in Deleting Favorites:", err.Error())
			return
		}
	}

	//TODO:删除用户的关注列表信息
	err = DAO.DB.Where("UID=?", id).Delete(&RelationshipSets.UserFollows{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting Follows:", err.Error())
		return
	}
	//TODO:删除被关注用户的对应粉丝列表信息
	err = DAO.DB.Where("FID=?", id).Delete(&RelationshipSets.UserFollowed{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting Followed:", err.Error())
		return
	}
	//TODO:删除用户对应等级信息
	err = DAO.DB.Where("UID=?", id).Delete(&EntitySets.Level{}).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting level:", err.Error())
		return
	}
	//TODO:删除用户信息
	var user *EntitySets.User
	err = DAO.DB.Debug().Where("UserID=?", id).Delete(&user).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::Logout", define.LogoutUserFailed, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting User:", err.Error())
		return
	}
	tx.Commit()
	Utilities.SendJsonMsg(c, 200, "注销账户成功")

}

// GetUserDetail
// @Tags User API
// @Summary 获取用户完整、详细的信息
// @Accept json
// @Produce json
// @Param UserID query string true "用户标识"
// @Success 200 {string}  json "{"code":"200","data":userInfo}"
// @Router /user/user-detail [get]
func GetUserDetail(c *gin.Context) {
	userID := c.Query("UserID")
	if userID == "" {
		Utilities.SendErrMsg(c, "service::Users::GetUserDetail", define.ObtainUserInformationFailed, "用户唯一标识不能为空")
		return
	}

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Omit("password").Where("UserID=?", userID).Preload("Videos").Preload("Favorites").
		Preload("Favorites.Videos").Preload("Comments").
		Preload("MessageBox").Preload("Follows").Preload("UserLevel").
		Preload("Followed").Preload("UserWatch").Preload("UserSearch").First(&userInfo).Error
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::GetUserDetail", define.ObtainUserInformationFailed, "Get User Info failed:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": userInfo,
	})
}

// ModifyUserSignature
// @Tags User API
// @summary 用户信息修改-更新用户签名
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param userSignature formData string false "用户签名,为空表示没有签名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户签名成功"}"
// @Router /user/ModifySignature [put]
func ModifyUserSignature(c *gin.Context) {
	id := c.Query("userID")
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

	Utilities.SendJsonMsg(c, 200, "修改用户签名成功")
}

// ModifyUserEmail
// @Tags User API
// @summary 用户信息修改-更新用户邮箱
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param userEmail formData string true "用户新邮箱"
// @Param Code formData string true "验证码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户邮箱成功"}"
// @Router /user/ModifyEmail [put]
func ModifyUserEmail(c *gin.Context) {
	//获取用户id,email,验证码
	id := c.Query("userID")
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
	Utilities.SendJsonMsg(c, 200, "用户邮箱成功")

}

// ForgetPassword
// @Tags User API
// @summary 重置用户密码
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param userEmail formData string true "用户邮箱"
// @Param Code formData string true "验证码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"重置用户密码成功"}"
// @Router /user/ForgetPassword [put]
func ForgetPassword(c *gin.Context) {
	id := c.Query("userID")
	userEmail := c.PostForm("userEmail")
	Code := c.PostForm("Code")
	newPassword := c.PostForm("newPassword")
	repeatPassword := c.PostForm("repeatPassword")

	var userInfo = new(EntitySets.User)
	verify, RDBErr := DAO.RDB.Get(c, userEmail).Result()
	err := DAO.DB.Where("Userid=?", id).First(&userInfo).Error
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
	Utilities.SendJsonMsg(c, 200, "重置密码成功")
	DAO.RDB.Del(c, userEmail)
}

// ModifyPassword
// @Tags User API
// @summary 用户信息修改-修改用户密码
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param password formData string true "用户密码"
// @Param newPassword formData string true "用户新密码"
// @Param repeatPassword formData string true "再次确认密码"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户密码成功"}"
// @Router /user/ModifyPassword [put]
func ModifyPassword(c *gin.Context) {
	id := c.Query("userID")
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
// @Param userID query string true "用户ID"
// @Param userName formData string true "用户名"
// @Success 200 {string}  json "{"code":"200","msg":"修改用户名成功"}"
// @Router /user/ModifyUserName [put]
func ModifyUserName(c *gin.Context) {
	id := c.Query("userID")
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
}

// UploadUserAvatar
// @Tags User API
// @summary 用户信息修改-上传头像
// @Accept multipart/form-data
// @Produce json
// // @Param Authorization header string true "token"
// @Param userID query string true "用户ID"
// @Param file formData file true "头像"
// @Success 200 {string}  json "{"code":"200","msg":"上传头像成功"}"
// @Router /user/face/upload/Avatar [post]
func UploadUserAvatar(c *gin.Context) {
	userID := c.Query("userID")
	FH, _ := c.FormFile("file") //FH=FileHeader
	//TODO:检查文件后缀名是否为 .jpg/.jpeg/.png/.jfif
	fname := FH.Filename
	extension := path.Ext(fname)
	println("ext:", extension)
	if _, ok := define.PicExtCheck[extension]; ok != true {
		Utilities.SendErrMsg(c, "service::Users::UploadUserAvatar", define.ImageFormatError, "图片格式错误或不支持此图片格式")
		return
	}

	//TODO:打开并读取文件
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
// @Param UID query string true "用户ID"
// @Param searchString query string true "搜索记录"
// @Router /user/AddSearchHistory [post]
func AddSearchHistory(c *gin.Context) {
	UID := c.Query("UID")
	searchString := c.Query("searchString")
	SearchHistory := EntitySets.SearchHistory{
		Model:        gorm.Model{},
		UID:          UID,
		SearchString: searchString,
	}
	err := SearchHistory.Create(DAO.DB)
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
// @Param userID query string true "用户ID"
// @Param VID query string true "视频ID"
// @Router /user/AddVideoHistory [post]
func AddVideoHistory(c *gin.Context) {
	UID := c.Query("UID")
	VID := c.Query("VID")
	VideoHistory := EntitySets.VideoHistory{
		MyModel: define.MyModel{},
		UID:     UID,
		VID:     VID,
	}
	err := VideoHistory.Create(DAO.DB)
	if err != nil {
		Utilities.SendErrMsg(c, "service::Users::AddVideoHistory", 5020, "创建视频历史记录失败")
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "创建视频历史记录成功")
}
