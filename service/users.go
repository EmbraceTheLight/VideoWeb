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
		Utilities.SendJsonMsg(c, 5001, "error while searching database:"+err.Error())
		return
	case countName > 0: //已有同名用户
		Utilities.SendJsonMsg(c, 4001, "用户名已存在，请重新输入待注册用户名")
		return
	case len(password) < 6: //密码长度小于6位
		Utilities.SendJsonMsg(c, 4002, "密码长度不能小于6位，请重新输入密码")
		return
	case password != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		Utilities.SendJsonMsg(c, 4003, "第一次输入的密码与第二次输入的密码不一致，请重新输入")
		return
	case utf8.RuneCountInString(Signature) > 25:
		Utilities.SendJsonMsg(c, 4010, "个性签名过长，请重新输入")
		return
	}

	//验证码获取及验证
	code := define.VerificationDataMap[email].Code
	ts := define.VerificationDataMap[email].Timestamp
	if code != verify {
		Utilities.SendJsonMsg(c, 4004, "验证码输入错误，请重新输入")
		return
	}
	if time.Now().Unix() > ts+define.Expired {
		Utilities.SendJsonMsg(c, 4013, "验证码已过期，请重新获取验证码")
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Utilities.SendJsonMsg(c, 5002, "密码加密错误")
		return
	}

	//设置用户默认头像
	file, err := os.Open("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWeb\\resources\\pictures\\default.jpg")
	defer file.Close()
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		return
	}
	var avatar = make([]byte, fileInfo.Size())
	_, err = file.Read(avatar)
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
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
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = defaultFavorites.Create(tx)
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = privateFavorites.Create(tx)
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	err = userLevel.Create(tx)
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
		tx.Rollback()
		return
	}
	token, err := logic.CreateToken(newUser.UserID, newUser.UserName, newUser.IsAdmin)
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "创建用户失败:"+err.Error())
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
		Utilities.SendJsonMsg(c, 4005, "邮箱不能为空")
		return
	}
	code, ts := logic.CreateVerificationCode()
	err := logic.SendCode(userEmail, code)
	fmt.Println(code)
	if err != nil {
		Utilities.SendJsonMsg(c, 4006, "验证码发送失败:"+err.Error())
		return
	}

	define.VerificationDataMap[userEmail] = define.VerificationData{
		Code:      code,
		Timestamp: ts,
	}

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
		Utilities.SendJsonMsg(c, 4007, "账号或密码为空")
		return
	}

	err = DAO.DB.Where("Account=? OR userName=?", AccountOrUserName, AccountOrUserName).First(&userInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound != nil { //未找到用户信息记录
			Utilities.SendJsonMsg(c, 4008, "账号不存在，请重新检查输入的账号")
			return
		}
		//其他错误
		Utilities.SendJsonMsg(c, 5003, "Get UserInfo failed:"+err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
		Utilities.SendJsonMsg(c, 4009, "密码错误")
		return
	}
	token, err := logic.CreateToken(userInfo.UserID, userInfo.UserName, userInfo.IsAdmin)
	if err != nil {
		Utilities.SendJsonMsg(c, 5004, "CreateToken error:"+err.Error())
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
		Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
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
			Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
			tx.Rollback()
			fmt.Println("Err in Deleting Favorites:", err.Error())
			return
		}
	}

	//TODO:删除用户的关注列表信息
	err = DAO.DB.Where("UID=?", id).Delete(&RelationshipSets.UserFollows{}).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting Follows:", err.Error())
		return
	}
	//TODO:删除被关注用户的对应粉丝列表信息
	err = DAO.DB.Where("FID=?", id).Delete(&RelationshipSets.UserFollowed{}).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting Followed:", err.Error())
		return
	}
	//TODO:删除用户对应等级信息
	err = DAO.DB.Where("UID=?", id).Delete(&EntitySets.Level{}).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting level:", err.Error())
		return
	}
	//TODO:删除用户信息
	var user *EntitySets.User
	err = DAO.DB.Debug().Where("UserID=?", id).Delete(&user).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5005, "注销用户失败:"+err.Error())
		tx.Rollback()
		fmt.Println("Err in Deleting User:", err.Error())
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注销账户成功",
	})

}

// GetUserDetail
// @Tags User API
// @Accept json
// @Produce json
// @Param UserID query string true "用户标识"
// @Success 200 {string}  json "{"code":"200","data":userInfo}"
// @Router /user/user-detail [get]
func GetUserDetail(c *gin.Context) {
	userID := c.Query("UserID")
	if userID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}

	var userInfo = new(EntitySets.User)
	err := DAO.DB.Debug().Omit("password").Where("UserID=?", userID).Preload("Videos").Preload("Favorites").
		Preload("Favorites.Videos").Preload("Comments").
		Preload("MessageBox").Preload("Follows").Preload("UserLevel").
		Preload("Followed").Preload("UserWatch").Preload("UserSearch").First(&userInfo).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5003,
			"msg":  "Get User Info failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": userInfo,
	})
}

// GetUserIpInfo
// @Tags private API
// @Accept json
// @Produce json
// @Success 200
// @Router /userInfo/user-IPInfo [get]
func GetUserIpInfo(c *gin.Context) {
	UserIP := c.ClientIP()
	fmt.Println(UserIP)
	info, err := Utilities.GetIPInfo(UserIP)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5005,
			"msg":  "获取用户IP信息失败",
		})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": info,
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
		c.JSON(http.StatusOK, gin.H{
			"code": 4010,
			"msg":  "个性签名过长，请重新输入",
		})
		return
	}

	err := DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("signature", signature).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5006,
			"mag":  "修改用户签名失败:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"mag":  "修改用户签名成功",
	})
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
	code := define.VerificationDataMap[userEmail].Code
	ts := define.VerificationDataMap[userEmail].Timestamp

	if code != verify {
		c.JSON(http.StatusOK, gin.H{
			"code": 4004,
			"msg":  "验证码输入错误，请重新输入",
		})
		return
	}
	if time.Now().Unix() > ts+define.Expired {
		Utilities.SendJsonMsg(c, -1, "验证码已过期，请重新获取验证码")
		return
	}
	//修改后存入数据库中
	err := DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("Email", userEmail).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4006,
			"msg":  "验证码发送失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改邮箱成功",
	})
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
	verify := define.VerificationDataMap[userEmail].Code
	ts := define.VerificationDataMap[userEmail].Timestamp
	err := DAO.DB.Where("Userid=?", id).First(&userInfo).Error
	switch {
	case err != nil:
		Utilities.SendJsonMsg(c, 5003, "获取用户信息失败")
		return
	case userEmail != userInfo.Email:
		Utilities.SendJsonMsg(c, 4011, "邮箱不匹配")
		return
	case Code != verify:
		Utilities.SendJsonMsg(c, 4004, "验证码输入错误，请重新输入")
		return
	case time.Now().Unix() > ts+define.Expired:
		Utilities.SendJsonMsg(c, -1, "验证码已过期，请重新获取验证码")
		return
	}

	status, err := logic.ModifyPassword(id, newPassword, repeatPassword)
	if err != nil {
		Utilities.SendJsonMsg(c, status, err.Error())
		return
	}
	Utilities.SendJsonMsg(c, 200, "重置密码成功")

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
		Utilities.SendJsonMsg(c, 5003, "获取用户信息失败")
		return
	}

	err = logic.ComparePassword(userInfo.Password, password)
	if err != nil {
		Utilities.SendJsonMsg(c, 4009, err.Error())
		return
	}

	code, err := logic.ModifyPassword(id, newPassword, repeatPassword)
	if err != nil {
		Utilities.SendJsonMsg(c, code, err.Error())
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
		Utilities.SendJsonMsg(c, 5003, "获取用户信息失败")
		return
	}

	var countName int64
	err = DAO.DB.Model(&EntitySets.User{}).Where("userName=?", userName).Count(&countName).Error
	switch {
	case err != nil: //数据库查询错误
		Utilities.SendJsonMsg(c, 5001, "error while searching database:"+err.Error())
		return
	case countName > 0: //已有同名用户
		Utilities.SendJsonMsg(c, 4001, "用户名已存在，请重新输入用户名")
		return
	}

	err = DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("userName", userName).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5009, "修改用户名失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改用户名成功",
	})

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
		Utilities.SendJsonMsg(c, 4012, "图片格式错误或不支持此图片格式")
		return
	}

	//TODO:打开并读取文件
	file, err := FH.Open()
	if err != nil {
		Utilities.SendJsonMsg(c, 5010, "打开文件失败"+err.Error())
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		Utilities.SendJsonMsg(c, 5011, "读取文件内容失败")
		return
	}

	err = DAO.DB.Model(&EntitySets.User{}).Where("userID=?", userID).Update("Avatar", data).Error
	if err != nil {
		Utilities.SendJsonMsg(c, 5012, "上传用户头像失败")
		fmt.Println("err:", err.Error())
		return
	}

	Utilities.SendJsonMsg(c, 200, "上传用户头像成功")
}
