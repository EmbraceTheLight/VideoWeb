package logic

import (
	"VideoWeb/define"
	"github.com/mojocn/base64Captcha"
)

// GenerateGraphicCaptcha 生成图形验证码
func GenerateGraphicCaptcha() (id, b64s, ans string, err error) {
	randRGBA := base64Captcha.RandColor()
	var driverString = base64Captcha.DriverString{
		Height:          75,
		Width:           150,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          4,
		Source:          "abcdefghijklmnopqrstuvwxyz0123456789",
		BgColor:         &randRGBA,
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	var driver base64Captcha.Driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, define.Store)
	id, b64s, ans, err = captcha.Generate()
	return
}

// GraphicCaptchaVerify 验证图形验证码
func GraphicCaptchaVerify(id, capt string) bool {
	if define.Store.Verify(id, capt, true) {
		return true
	}
	return false
}
