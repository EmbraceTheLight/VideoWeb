package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func GetFavoritesByID(id string) ([]*EntitySets.Favorites, error) {
	var favorites []*EntitySets.Favorites
	err := DAO.DB.Where("UID = ?", id).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// GetUserVideosByID 方法名中的User是为了与收藏夹Favorites的视频区分开来
func GetUserVideosByID(id string) ([]*EntitySets.Video, error) {
	var videos []*EntitySets.Video
	err := DAO.DB.Where("UID = ?", id).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil

}

func GetFollowsByUserID(id string) ([]*RelationshipSets.UserFollows, error) {
	var follows []*RelationshipSets.UserFollows
	err := DAO.DB.Where("UID = ?", id).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func GetFollowedByFollowedID(id string) ([]*RelationshipSets.UserFollowed, error) {
	var followed []*RelationshipSets.UserFollowed
	err := DAO.DB.Where("UID = ?", id).Find(&followed).Error
	if err != nil {
		return nil, err
	}
	return followed, nil
}

// GetUserCommentsByID 方法名中的User是为了与Video的评论区分开来
func GetUserCommentsByID(id string) ([]*EntitySets.Comments, error) {
	var comments []*EntitySets.Comments
	err := DAO.DB.Where("UID = ?", id).Find(&comments).Error
	if err != nil {
		return comments, nil
	}
	return comments, nil

}

func GetMsgBoxByID(id string) ([]*EntitySets.Message, error) {
	var message []*EntitySets.Message
	err := DAO.DB.Where("UID = ?", id).Find(&message).Error
	if err != nil {
		return message, nil
	}
	return message, nil
}

func ComparePassword(userPassword, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword)); err != nil {
		return errors.New("密码错误")
	}
	return nil
}

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
		return 5002, errors.New("密码加密错误")
	}

	//更新用户密码
	err = DAO.DB.Model(&EntitySets.User{}).Debug().Where("UserID=?", id).Update("Password", string(hashedPassword)).Error
	if err != nil {
		return 5009, errors.New("修改密码失败")
	}

	return 200, nil
}
