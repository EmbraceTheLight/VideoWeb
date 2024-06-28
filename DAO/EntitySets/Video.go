package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

// VideoSummary 视频摘要信息
type VideoSummary struct {
	VideoID     int64  `json:"VID"`
	VCName      string `json:"VCName" gorm:"column:user_name"` //Video Creator Name
	CntBarrages uint32 `json:"CntBarrages"`
	Title       string `json:"Title"`
	CntViews    uint32 `json:"CntViews"`
	Duration    string `json:"Duration"`
	CoverPath   string `json:"CoverPath"`
}

type Video struct {
	define.MyModel

	//用户显式指定
	Title       string `json:"Title" gorm:"column:title;type:varchar(100);"`                   //视频题目
	Description string `json:"Description" gorm:"column:description;type:text;"`               //视频描述
	Class       string `json:"Class" gorm:"column:class;type:varchar(20);index:idx_class_hot"` //视频分类
	Path        string `json:"Path" gorm:"column:path;type:varchar(200)"`                      //视频路径
	//系统默认生成
	VideoID      int64  `json:"VID" gorm:"column:video_id;type:bigint;primaryKey"`                     //视频唯一标识
	UID          int64  `json:"UID" gorm:"column:user_id;type:bigint"`                                 //视频作者ID
	UserName     string `json:"UserName" gorm:"column:user_name;type:varchar(40)"`                     //视频作者名称
	Likes        uint32 `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`                 //视频被点赞数
	Shells       uint32 `json:"Shells" gorm:"column:shells;type:int unsigned;default:0"`               //视频获得的贝壳数
	Hot          uint32 `json:"Hot" gorm:"column:hot;type:int unsigned;default:0;index:idx_class_hot"` //视频热度,默认排序的依据
	CntBarrages  uint32 `json:"CntBarrages" gorm:"column:cnt_barrages;type:int unsigned;default:0"`    //视频弹幕数
	CntShares    uint32 `json:"CntShares" gorm:"column:cnt_shares;type:int unsigned;default:0"`        //视频被分享数
	CntFavorites uint32 `json:"CntFavorites" gorm:"column:cnt_favorites;type:int unsigned;default:0"`  //视频被收藏数,似乎可以通过联表计算出来
	CntViews     uint32 `json:"CntViews" gorm:"column:cnt_views;type:int unsigned;default:0"`          //视频播放量
	Duration     string `json:"Duration" gorm:"column:duration;type:varchar(10);"`                     //视频时长,单位秒
	Size         int64  `json:"Size" gorm:"column:size;type:int;"`                                     //视频大小,单位字节
	CoverPath    string `json:"VideoCover" gorm:"type:varchar(200)"`                                   //视频封面路径
	//以下三项均为一对多
	Comments []*Comments `gorm:"foreignKey:VID;references:VideoID"` //视频评论
	Tags     []*Tags     `gorm:"foreignKey:VID;references:VideoID"` //视频标签
	Barrages []*Barrage  `gorm:"foreignKey:VID;references:VideoID"` //视频弹幕
}

func (f *Video) TableName() string {
	return "video"
}

// GetVideoInfoByID 根据视频ID获得视频信息
func GetVideoInfoByID(db *gorm.DB, VID int64) (*Video, error) {
	var info = new(Video)
	err := db.Where("video_id=?", VID).First(info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

// InsertVideoRecord 插入视频信息
func InsertVideoRecord(db *gorm.DB, video *Video) error {
	return db.Model(&Video{}).Create(video).Error
}

// DeleteVideoInfoByVideoID 根据视频ID删除视频信息
func DeleteVideoInfoByVideoID(db *gorm.DB, VID int64) error {
	return db.Model(&Video{}).Where("video_id=?", VID).Delete(&Video{}).Error
}

// DeleteVideoInfoByUserID 根据用户ID删除视频信息
func DeleteVideoInfoByUserID(db *gorm.DB, UID int64) error {
	return db.Debug().Model(&Video{}).Where("user_id=?", UID).Delete(&Video{}).Error
}

// UpdateVideoField 根据视频ID以及字段名更新视频某字段
func UpdateVideoField(db *gorm.DB, VID int64, field string, change int) error {
	return db.Model(&Video{}).Where("video_id=?", VID).Update(field, gorm.Expr(field+"+?", change)).Error
}

func UpdateVideoCover(db *gorm.DB, VID int64, coverData []byte) error {
	return db.Model(&Video{}).Where("video_id=?", VID).Update("VideoCover", coverData).Error
}
