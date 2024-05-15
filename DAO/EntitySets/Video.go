package DAO

import (
	"VideoWeb/define"
	"gorm.io/gorm"
)

type Video struct {
	define.MyModel

	//用户显式指定
	Title       string `json:"Title" gorm:"column:title;type:varchar(100);index:,class:FULLTEXT"` //视频题目,添加了全文索引
	Description string `json:"Description" gorm:"column:description;type:text"`                   //视频描述
	Class       string `json:"Class" gorm:"column:class;type:varchar(20)"`                        //视频分类
	Path        string `json:"Path" gorm:"column:path;type:varchar(200)"`                         //视频路径
	//系统默认生成
	VideoID      int64  `json:"VID" gorm:"column:video_id;type:bigint;primaryKey"`                    //视频唯一标识
	UID          int64  `json:"UID" gorm:"column:user_id;type:bigint"`                                //视频作者ID
	Likes        uint32 `json:"Likes" gorm:"column:likes;type:int unsigned;default:0"`                //视频被点赞数
	UnLikes      uint32 `json:"UnLikes" gorm:"column:unlikes;type:int unsigned;default:0"`            //视频被点踩数
	Shells       uint32 `json:"Shells" gorm:"column:shells;type:int unsigned;default:0"`              //视频获得的贝壳数
	CntFavorites uint32 `json:"CntFavorites" gorm:"column:cnt_favorites;type:int unsigned;default:0"` //视频被收藏数,似乎可以通过联表计算出来
	CntViews     uint32 `json:"CntViews" gorm:"column:cnt_views;type:int unsigned;default:0"`         //视频播放量
	Duration     string `json:"Duration" gorm:"column:duration;type:varchar(10);"`                    //视频时长,单位秒
	Size         int64  `json:"Size" gorm:"column:size;type:int;"`                                    //视频大小,单位字节
	//以下三项均为一对多
	Comments []*Comments `gorm:"foreignKey:VID;references:VideoID"` //视频评论
	Tags     []*Tags     `gorm:"foreignKey:VID;references:VideoID"` //视频标签
	Barrages []*Barrage  `gorm:"foreignKey:VID;references:VideoID"` //视频弹幕
	//二进制大对象,存放视频封面信息
	VideoCover []byte `json:"VideoCover" gorm:"column:cover;type:MediumBLOB;size:10240000"` //视频封面,最大为10MiB
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
func UpdateVideoField(db *gorm.DB, VID int64, fields string, change int) error {
	return db.Model(&Video{}).Where("video_id=?", VID).Update(fields, gorm.Expr(fields+"+?", change)).Error
}

func UpdateVideoCover(db *gorm.DB, VID int64, coverData []byte) error {
	return db.Model(&Video{}).Where("video_id=?", VID).Update("VideoCover", coverData).Error
}
