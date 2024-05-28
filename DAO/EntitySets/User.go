package DAO

import (
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/define"
	"gorm.io/gorm"
)

type UserSummary struct {
	UserID    int64  `json:"userID" `
	UserName  string `json:"userName" `
	Avatar    []byte `json:"avatar"`
	UserLevel uint16 `json:"userLevel" `
}
type User struct {
	define.MyModel
	UserID    int64  `json:"userID" gorm:"column:user_id;type:bigint;primaryKey"` //用户ID
	UserName  string `json:"userName" gorm:"column:user_name;type:varchar(40)"`   //用户名
	Password  string `json:"password" gorm:"column:password;type:varchar(72)"`    //用户密码,最多72位
	Email     string `json:"email" gorm:"column:email;type:varchar(100)"`         //用户邮箱
	Signature string `json:"signature" gorm:"column:signature;type:varchar(25)"`  //个性签名，至多25个字
	//以下内容为数据库系统默认生成
	Shells        uint32                           `json:"shells" gorm:"column:shells;type:int unsigned;default:1000"`       //用户拥有的贝壳数
	IsAdmin       int                              `json:"isAdmin" gorm:"column:is_admin;type:tinyint;default:-1"`           //用户是否为管理员,-1表示不是管理员，1表示是管理员
	CntMsgNotRead int32                            `json:"count" gorm:"column:cnt_msg_not_read;type:int unsigned;default:0"` //用户未读消息数
	UserLevel     Level                            `json:"userLevel" gorm:"foreignKey:UID;references:UserID"`                //用户等级
	Videos        []*Video                         `json:"videos" gorm:"foreignKey:UID;references:UserID;"`                  //has-many 一对多
	UserWatch     []*VideoHistory                  `json:"userHistory" gorm:"foreignKey:UID;references:UserID"`              //用户观看记录,has-many
	Favorites     []*Favorites                     `json:"favorites" gorm:"foreignKey:UID;references:UserID"`                //用户收藏夹has-many
	Comments      []*Comments                      `json:"comments" gorm:"foreignKey:UID;references:UserID"`                 //用户评论has-many
	UserSearch    []*SearchHistory                 `json:"userSearch" gorm:"foreignKey:UID;references:UserID"`               //用户搜索记录,has-many
	Follows       []*FollowList                    `json:"follows" gorm:"foreignKey:UID;references:UserID"`                  //用户关注列表 has-many
	Followed      []*RelationshipSets.UserFollowed `json:"followed" gorm:"foreignKey:UID;references:UserID"`                 //用户粉丝列表 has-many
	Avatar        []byte                           `json:"avatar" gorm:"type:MediumBLOB;size:10240000"`                      //用户头像,最大为10MiB
}

func (u *User) TableName() string {
	return "user"
}

// InsertUserRecord 插入用户记录
func InsertUserRecord(db *gorm.DB, user *User) (err error) {
	result := db.Model(User{}).Create(&user)
	err = result.Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserVideosByID 通过用户ID来获取该用户的发送视频的列表
func GetUserVideosByID(db *gorm.DB, id string) ([]*Video, error) {
	var videos []*Video
	err := db.Model(User{}).Where("user_id = ?", id).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// GetUserInfoByName 通过用户名来获取该用户的详细信息
func GetUserInfoByName(db *gorm.DB, name string) (*User, error) {
	var user *User
	err := db.Model(User{}).Where("user_name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil

}

// GetUserCommentsByID 通过用户ID来获取该用户的评论列表
func GetUserCommentsByID(DB *gorm.DB, id string) ([]*Comments, error) {
	var comments []*Comments
	err := DB.Where("user_id = ?", id).Find(&comments).Error
	if err != nil {
		return comments, nil
	}
	return comments, nil

}

// UpdateUserNumField 根据用户ID以及字段名更新用户数值字段
func UpdateUserNumField(db *gorm.DB, UID int64, field string, change int) error {
	return db.Model(&User{}).Where("user_id=?", UID).Update(field, gorm.Expr(field+"+?", change)).Error
}

// UpdateUserStringField 根据用户ID以及字段名更新用户字符串字段
func UpdateUserStringField(db *gorm.DB, UID int64, field string, change string) error {
	return db.Model(&User{}).Where("user_id=?", UID).Update(field, change).Error
}

// DeleteUserRecordByID 根据用户ID删除用户记录
func DeleteUserRecordByID(db *gorm.DB, id int64) error {
	return db.Model(User{}).Where("user_id = ?", id).Delete(&User{}).Error
}

// GetUserInfoByID 根据用户ID获取用户信息
func GetUserInfoByID(db *gorm.DB, id int64) (*User, error) {
	var user *User
	err := db.Model(&User{}).Where("user_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserAvatar(db *gorm.DB, userID int64, data []byte) error {
	return db.Model(&User{}).Where("user_id=?", userID).Update("avatar", data).Error
}
