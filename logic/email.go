package logic

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

// SendCode 发送验证码
func SendCode(toUser, code string) error {
	e := email.NewEmail()
	e.Subject = "验证码"
	e.HTML = []byte("您的验证码为： <b>" + code + "</b>")

	e.From = "zey <1010642166@qq.com>"
	e.To = []string{toUser}

	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1010642166@qq.com", "exwxhwxuqwljbfdc", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: false, ServerName: "smtp.qq.com"})
	if err != nil {
		log.Printf("Error sending Email to %s:%v", toUser, err)
	}
	return err
}

// CreateVerificationCode 生成6位验证码
func CreateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	t1 := rand.Int() % 1000000
	ret := fmt.Sprintf("%06d", t1)
	return ret

}
