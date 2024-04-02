package test

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	//e.From = "zey <10eltzey10@gmail.com>"
	e.Subject = "邮箱验证码发送测试"
	e.HTML = []byte("您的验证码为： <b>625</b>")
	//err := e.Send("smtp.gmail.com:465", smtp.PlainAuth("", "10eltzey10@gmail.com", "nzimqkbvxnozyrpq", "smtp.gmail.com"))

	e.From = "zey <1010642166@qq.com>"
	//e.To = []string{"3293136088@qq.com"}
	e.To = []string{"10eltzey10@gmail.com"}

	//err := e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", "10eltzey10@gmail.com", "nzimqkbvxnozyrpq", "smtp.gmail.com"),
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1010642166@qq.com", "exwxhwxuqwljbfdc", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: false, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateVerificationCode(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	t1 := rand.Int() % 1000000
	ret := fmt.Sprintf("%06d", t1)
	println(ret)

	t1 = 1
	ret = fmt.Sprintf("%06d", t1)
	println(ret)

}
