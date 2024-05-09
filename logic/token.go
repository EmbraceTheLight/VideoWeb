package logic

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte("secretString")

type UserClaims struct {
	UserName string `json:"userName"`
	UserId   int64  `json:"userId"`
	IsAdmin  int    `json:"isAdmin"`
	jwt.RegisteredClaims
}

// CreateToken 生成Token
func CreateToken(identity int64, name string, isAdmin int) (string, error) {
	userClaims := UserClaims{
		UserName:         name,
		UserId:           identity,
		IsAdmin:          isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	fmt.Println("length of tokenString:", len(tokenString))
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)

	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if token.Valid == false {
		fmt.Println("Error parsing token:token is invalid!")
		return nil, errors.New("token is invalid")
	}
	return userClaim, nil
}
