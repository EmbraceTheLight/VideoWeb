package logic

import (
	"VideoWeb/define"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

func GenerateGraphicCaptcha() (id, b64s, ans string, err error) {
	var driverString = base64Captcha.DriverString{
		Height:          30,
		Width:           60,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          4,
		Source:          "abcdefghijklmnopqrstuvwxyz0123456789",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 254,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}
	var driver base64Captcha.Driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, define.Store)
	id, b64s, ans, err = captcha.Generate()
	return
}

func GraphicCaptchaVerify(id, capt string) bool {
	if define.Store.Verify(id, capt, true) {
		return true
	}
	return false
}
