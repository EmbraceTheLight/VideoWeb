package service

import (
	"VideoWeb/DAO"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendCode
// @Tags Captcha API
// @summary 发送验证码
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱"
// @Success 200 {string}  json "{"code":"200","data":"data"}"
// @Router /send-code [get]
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

	DAO.RDB.Set(c, userEmail, code, define.Expired)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功！",
	})
}

// GenerateGraphicCaptcha
// @Tags Captcha API
// @Summary 生成图形验证码
// @Description 生成图形验证码并返回给客户端
// @Produce json
// @Success 200 {string} json {"captcha_result": CaptchaResult}
// @Router /GenerateGraphicCaptcha [get]
func GenerateGraphicCaptcha(c *gin.Context) {
	id, b64s, ans, err := logic.GenerateGraphicCaptcha()
	if err != nil {
		Utilities.SendErrMsg(c, "service::Captcha::GenerateGraphicCaptcha", define.CaptchaGenerateFailed, "图形验证码生成失败:"+err.Error())
		return
	}
	captchaResult := define.GraphicCaptcha{
		ID:     id,
		B64lob: b64s,
		Answer: ans,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"msg":     "图形验证码生成成功!",
		"captcha": captchaResult,
	})
}

// CheckGraphicCaptcha
// @Tags Captcha API
// @Summary 验证图形验证码
// @Description 验证图形验证码并返回给客户端
// @Accept multipart/form-data
// @Produce json
// @Param id query string true "验证码ID"
// @Param input formData string true "用户输入的验证码"
// @Success 200 {string} json {"captcha_result": CaptchaResult}
// @Router /CheckGraphicCaptcha [post]
func CheckGraphicCaptcha(c *gin.Context) {
	id := c.Query("id")
	input := c.PostForm("input")
	right := logic.GraphicCaptchaVerify(id, input)
	if !right {
		Utilities.SendJsonMsg(c, http.StatusOK, "图形验证码验证失败！")
		return
	}
	Utilities.SendJsonMsg(c, http.StatusOK, "图形验证码验证成功！")
}
