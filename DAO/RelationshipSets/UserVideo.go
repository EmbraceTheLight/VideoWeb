package DAO

import (
	"VideoWeb/define"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

//	type UserVideo struct {
//		define.MyModel
//		UID int64 `json:"uid" gorm:"primaryKey;column:user_id;type:bigint;index:idx_uid_vid"`
//		VID int64 `json:"vid" gorm:"primaryKey;column:video_id;type:bigint;index:idx_uid_vid"`
//
//		IsLike    bool `json:"is_like" gorm:"column:is_like;type:tinyint(1);default:0"`
//		IsDislike bool `json:"is_unlike" gorm:"column:is_dislike;type:tinyint(1);default:0"`
//		IsFavor   bool `json:"is_favor" gorm:"column:is_favor;type:tinyint(1);default:0"`
//	}

type UserVideo struct {
	UID       int64 `bson:"user_id,omitempty"`
	VID       int64 `bson:"video_id,omitempty"`
	IsLike    bool  `bson:"is_like,omitempty"`
	IsDislike bool  `bson:"is_dislike,omitempty"`
	IsFavor   bool  `bson:"is_favor,omitempty"`
}

// InsertUserVideoRecord 插入用户视频记录
func InsertUserVideoRecord(db *mongo.Client, uv *UserVideo) error {
	ctx, cancel := context.WithTimeout(context.Background(), define.MongoOperationTimeout)
	defer cancel()
	//err := db.Model(&UserVideo{}).Save(uv).Error
	filter := bson.M{
		"user_id":  uv.UID,
		"video_id": uv.VID,
	}
	replaceOpts := options.Replace().SetUpsert(true)
	_, err := db.Database("video_web").Collection("user_video").ReplaceOne(ctx, filter, uv, replaceOpts)
	if err != nil {
		return fmt.Errorf("UserVideo.InsertUserVideo: %w", err)
	}
	return nil
}

// GetUserVideoRecord 根据用户ID和视频ID获取用户视频记录
func GetUserVideoRecord(db *mongo.Client, uid int64, vid int64) (*UserVideo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), define.MongoOperationTimeout)
	defer cancel()
	var uv = new(UserVideo)
	//err := db.Model(&UserVideo{}).Where("user_id =? AND video_id = ?", uid, vid).First(uv).Error
	filter := bson.M{
		"user_id":  uid,
		"video_id": vid,
	}
	err := db.Database("video_web").Collection("user_video").FindOne(ctx, filter).Decode(uv)
	if err != nil {
		return nil, err
	}

	return uv, nil
}

// UpdateUserVideoRecord 更新用户视频记录的三个bool状态
func UpdateUserVideoRecord(db *mongo.Client, userID, videoID int64, field string, change bool) error {
	//return db.Model(&UserVideo{}).Where("userID =? AND videoID = ?", userID, videoID).Update(field, change).Error
	ctx, cancel := context.WithTimeout(context.Background(), define.MongoOperationTimeout)
	defer cancel()
	filter := bson.M{
		"userID":  userID,
		"videoID": videoID,
	}
	update := bson.M{
		"$set": bson.M{
			field: change,
		},
	}
	_, err := db.Database("video_web").Collection("user_video").UpdateOne(ctx, filter, update)
	return err
}

// DeleteUserVideoRecordsByUserID 根据用户ID删除用户所有观看过的视频记录
func DeleteUserVideoRecordsByUserID(db *gorm.DB, UserID int64) error {
	return db.Model(&UserVideo{}).Where("user_id =?", UserID).Delete(&UserVideo{}).Error
}

// DeleteUserVideoRecordsByVideoID 根据视频ID删除所有观看过该视频的用户记录
func DeleteUserVideoRecordsByVideoID(db *gorm.DB, videoID int64) error {
	return db.Model(&UserVideo{}).Where("video_id =?", videoID).Delete(&UserVideo{}).Error
}
