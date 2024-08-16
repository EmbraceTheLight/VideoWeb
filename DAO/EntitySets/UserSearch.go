package DAO

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//type UserSearchHistory struct {
//	define.MyModel
//	UID          int64  `json:"UID" gorm:"column:user_id;type:bigint;primaryKey"`
//	SearchString string `json:"searchString" gorm:"column:search_string;type:varchar(100);"` //用户搜索字符串
//}

type UserSearchHistory struct {
	SearchTime   time.Time `bson:"search_time,omitempty"`
	UID          int64     `bson:"user_id,omitempty"`
	SearchString string    `bson:"search_string,omitempty"`
}

func (s *UserSearchHistory) GetScore() float64 {
	return float64(s.SearchTime.UnixMilli())
}

func (s *UserSearchHistory) GetValue() any {
	return s.SearchString
}

// InsertSearchRecord 插入或更新用户搜索记录
func InsertSearchRecord(mongo *mongo.Client, userID int64, searchString string) error {
	filter := bson.M{
		"user_id":       userID,
		"search_string": searchString,
	}

	updateOpts := options.Update().SetUpsert(true)
	update := bson.M{
		"$set": bson.M{
			"search_time": time.Now(),
		},
		"$setOnInsert": bson.M{
			"user_id":       userID,
			"search_string": searchString,
		},
	}
	//err := db.Where("user_id = ?", UID).Omit("created_at").Save(sh).Error
	_, err := mongo.Database("video_web").
		Collection("user_search_history").UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		return fmt.Errorf("UserSearch.InsertSearchRecord: %w", err)
	}
	return nil
}

// DeleteOneSearchRecord 根据记录ID删除单条搜索记录
func DeleteOneSearchRecord(mongo *mongo.Client, ID int64) error {
	//return db.Model(&UserSearchHistory{}).Where("id = ?", ID).Delete(&UserSearchHistory{}).Error
	filter := bson.D{{"_id", ID}}
	_, err := mongo.Database("video_web").
		Collection("user_search_history").DeleteOne(context.TODO(), filter)
	return err
}

// DeleteAllSearchRecord 根据用户ID删除用户的所有搜索记录
func DeleteAllSearchRecord(mongo *mongo.Database, userID int64) error {
	//return db.Model(&UserSearchHistory{}).Where("user_id = ?", uid).Delete(&UserSearchHistory{}).Error
	filter := bson.D{{"user_id", userID}}
	_, err := mongo.Collection("user_search_history").DeleteOne(context.TODO(), filter)
	return err
}
