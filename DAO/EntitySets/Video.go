package DAO

import (
	"VideoWeb/define"
)

type Video struct {
	define.MyModel

	VideoID     string `json:"VID" gorm:"column:videoID;type:char(36);primaryKey"` //视频唯一标识
	UID         string `json:"UID" gorm:"column:UID;type:char(36)"`                //视频作者ID
	Title       string `json:"Title" gorm:"column:title;type:varchar(100)"`        //视频题目
	Description string `json:"Description" gorm:"column:description;type:text"`    //视频描述
	Class       string `json:"Class" gorm:"column:class;type:varchar(20)"`         //视频分类
	Path        string `json:"Path" gorm:"column:path;type:varchar(200)"`          //视频路径

	Likes        uint32      `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`               //视频被点赞数
	UnLikes      uint32      `json:"UnLikes" gorm:"column:unlikes;type:int unsigned;default:0"`           //视频被点踩数
	Shells       uint32      `json:"Shells" gorm:"column:shells;type:int unsigned;default:0"`             //视频获得的贝壳数
	CntFavorites uint32      `json:"CntFavorites" gorm:"column:cntFavorites;type:int unsigned;default:0"` //视频被收藏数,似乎可以通过联表计算出来
	CntViews     uint32      `json:"CntViews" gorm:"column:cntViews;type:int unsigned;default:0"`         //视频播放量
	Comments     []*Comments `gorm:"foreignKey:VID;references:videoID"`
	Tags         []*Tags     `gorm:"foreignKey:VID;references:VideoID"`
}

func (f *Video) TableName() string {
	return "Video"
}
