package DAO

import (
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/define"
	"gorm.io/gorm"
)

type User struct {
	define.MyModel

	UserID    string `json:"userID" gorm:"column:UserID;type:char(36);primaryKey"` //用户ID
	UserName  string `json:"userName" gorm:"column:userName;type:varchar(40)"`     //用户名
	Password  string `json:"password" gorm:"column:password;type:char(60)"`        //用户密码
	Email     string `json:"email" gorm:"column:email;type:varchar(100)"`          //用户邮箱
	Signature string `json:"signature" gorm:"column:signature;type:varchar(25)"`   //个性签名，至多25个字
	//以下内容为数据库系统默认生成
	Account       int32                            `json:"account" gorm:"column:account;type:int auto_increment;unique"`  //用户账号
	Shells        uint32                           `json:"shells" gorm:"column:shells;type:int unsigned;default:1000"`    //用户拥有的贝壳数
	IsAdmin       int                              `json:"isAdmin" gorm:"column:isAdmin;type:tinyint;default:-1"`         //用户是否为管理员,-1表示不是管理员，1表示是管理员
	CntMsgNotRead int32                            `json:"count" gorm:"column:CntMsgNotRead;type:int unsigned;default:0"` //用户未读消息数
	UserLevel     Level                            `json:"userLevel" gorm:"foreignKey:UID;references:UserID"`             //用户等级
	Videos        []*Video                         `json:"videos" gorm:"foreignKey:UID;references:UserID;"`               //has-many 一对多
	UserWatch     []*VideoHistory                  `json:"userHistory" gorm:"foreignKey:UID;references:UserID"`           //用户观看记录，has-many
	Favorites     []*Favorites                     `json:"favorites" gorm:"foreignKey:UID;references:UserID"`             //用户收藏夹has-many
	Comments      []*Comments                      `json:"comments" gorm:"foreignKey:UID;references:UserID"`              //用户评论has-many
	MessageBox    []*Message                       `json:"message" gorm:"foreignKey:UID;references:UserID;"`              //用户消息盒子has-many
	UserSearch    []*SearchHistory                 `json:"userSearch" gorm:"foreignKey:UID;references:UserID"`            //用户搜索记录
	Follows       []*RelationshipSets.UserFollows  `json:"follows" gorm:"foreignKey:UID;references:UserID"`               //用户关注列表 many2many 多对多
	Followed      []*RelationshipSets.UserFollowed `json:"followed" gorm:"foreignKey:UID;references:UserID"`              //用户粉丝列表 many2many 多对多
	Avatar        []byte                           `json:"avatar" gorm:"type:MediumBLOB;size:10240000"`                   //用户头像,最大为10MiB
}

func (u *User) TableName() string {
	return "Users"
}

func (u *User) Create(DB *gorm.DB) (Account int32, err error) { //返回用户账号以及错误信息
	result := DB.Create(&u).Scan(&u)
	err = result.Error
	if err != nil {
		return 0, err
	}
	Account = u.Account
	return Account, nil
}
