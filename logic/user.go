package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode/utf8"
)

// GetUserIpInfo 获取并返回用户所在国家和地区
func GetUserIpInfo(c *gin.Context) (Country, City string) {
	UserIP := c.ClientIP()
	fmt.Println(UserIP)
	info, err := Utilities.GetIPInfo(UserIP)
	if err != nil {
		fmt.Println("获取用户IP失败:", err)
		return "", ""
	}
	return info.Country, info.City
}

// GetUserNameByID 通过用户ID获取用户名
func GetUserNameByID(id string) (string, error) {
	var userName string
	err := DAO.DB.Where("UserID=?", id).Pluck("userName", &userName).Limit(1).Error
	if err != nil {
		return "", err
	}
	return userName, nil
}

// GetFavoritesByID 通过用户ID来获取该用户的收藏夹列表
func GetFavoritesByID(id string) ([]*EntitySets.Favorites, error) {
	var favorites []*EntitySets.Favorites
	err := DAO.DB.Where("UID = ?", id).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// ComparePassword  比较用户输入的密码与数据库中的密码
func ComparePassword(userPassword, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword)); err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// ModifyPassword 用户修改辅助函数
func ModifyPassword(id, newPassword, repeatPassword string) (int, error) {
	switch {
	case len(newPassword) < 6: //密码长度小于6位
		return 4002, errors.New("密码长度不能小于6位，请重新输入密码")
	case newPassword != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		return 4003, errors.New("第一次输入的密码与第二次输入的密码不一致，请重新输入")
	}
	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return define.PasswordEncryptionError, errors.New("密码加密错误")
	}

	//更新用户密码
	err = DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("Password", string(hashedPassword)).Error
	if err != nil {
		return define.ModifyPasswordFailed, errors.New("修改密码失败")
	}

	return http.StatusOK, nil
}

// CheckRegisterInfo 检查注册信息是否正确。
func CheckRegisterInfo(checkInfo *EntitySets.User, repeatPassword string) error {
	var countName int64
	err := DAO.DB.Model(&EntitySets.User{}).Where("userName=?", checkInfo.UserName).Count(&countName).Error
	if err != nil {
		return errors.New("查询用户信息失败")
	}
	switch {
	case countName > 0: //已有同名用户
		return errors.New("用户名已存在")
	case len(checkInfo.Password) < 6: //密码长度小于6位
		return errors.New("密码长度小于6位")
	case checkInfo.Password != repeatPassword: //第一次输入的密码与第二次输入的密码不一致
		return errors.New("第一次输入的密码与第二次输入的密码不一致")
	case utf8.RuneCountInString(checkInfo.Signature) > 25:
		return errors.New("个性签名过长")
	}
	return nil
}

// UpdateUserFieldForUpdate 更新用户某个字段(悲观锁)
func UpdateUserFieldForUpdate(c *gin.Context, UserID int64, field string, change int) error {
	tx := DAO.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Set("gorm:query_option", "FOR UPDATE") //添加行级锁(悲观)
	err := EntitySets.UpdateUserField(tx, UserID, field, change)
	funcName, _ := c.Get("funcName")
	if err != nil {
		tx.Rollback()
		Utilities.SendErrMsg(c, funcName.(string), 5000, err.Error())
		return err
	}
	tx.Commit()
	return nil
}
