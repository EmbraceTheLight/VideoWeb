package test

import (
	"VideoWeb/logic"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

type userClaims struct {
	UserName string `json:"userName"`
	UserId   string `json:"userId"`
	jwt.RegisteredClaims
}

// 生成Token
func TestCreateToken(t *testing.T) {
	userClaims := userClaims{
		UserName:         "zey",
		UserId:           "zey",
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString(logic.JwtSecret)
	if err != nil {
		t.Fatal(err)
	}
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InpleSIsInVzZXJJZCI6InpleSJ9.pjzy4qoZTrpF1_9hE5Vp9dGN7EGs63DhqByFUWDQvCQ
	fmt.Println("token is:", tokenString)
}

// 解析token
func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InpleSIsInVzZXJJZCI6InpleSJ9.pjzy4qoZTrpF1_9hE5Vp9dGN7EGs63DhqByFUWDQvCQ"
	userClaim := new(userClaims)

	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return logic.JwtSecret, nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if token.Valid == false {
		fmt.Println("Error parsing token:token is invalid!")
		return
	}
	fmt.Println(userClaim)
}
