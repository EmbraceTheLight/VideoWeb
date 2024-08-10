package DAO

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities/logf"
	"VideoWeb/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

var (
	DB      *gorm.DB
	MongoDB *mongo.Client
	RDB     *redis.Client
)

func initRedisClient() *redis.Client {
	RedisConf := config.GetConfig().DBConf.RedisConf
	logf.WriteInfoLog("initRedisClient", "init redis client success")
	return redis.NewClient(
		&redis.Options{
			Addr:     RedisConf.Host + ":" + RedisConf.Port,
			Password: RedisConf.Password,
			DB:       0,
		})
}

func initMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI("mongodb://localhost:27017").
			SetMaxPoolSize(150))
	if err != nil {
		log.Println("[initMongo] Error connecting to MongoDB: ", err)
		panic(err)
	}
	//db := client.Database("video_web")

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("[initMongo] Error pinging MongoDB: ", err)
		panic(err)
	}

	logf.WriteInfoLog("initMongo", "Mongo init success")
	return client
}

func createDatabase() {
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	dbConnection := fmt.Sprintf("%s:%s@(%s:%s)/mysql?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Charset)
	db, _ := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	db.Exec("CREATE DATABASE  IF NOT EXISTS videoweb")
}
func connectMysql() {
	var err error
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	DBConnection := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Database, MySQLConf.Charset)
	DB, err = gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		panic(err)
	}
}
func checkAndCreateTable() {
	var err error
	var tables = []interface{}{
		&EntitySets.User{},
		&EntitySets.Level{},
		&EntitySets.Favorites{},
		&EntitySets.Video{},
		&EntitySets.Barrage{},
		&EntitySets.Tags{},
		&RelationshipSets.FavoriteVideo{},
		&EntitySets.Comments{},
		&RelationshipSets.UserFollowed{},
		&RelationshipSets.UserFollows{},
		&EntitySets.FollowList{},
		&EntitySets.UserWatch{},
		&EntitySets.UserSearchHistory{},
		&RelationshipSets.UserVideo{},
		&RelationshipSets.UserLikedComments{},
		&RelationshipSets.UserDislikedComments{},
	}
	err = DB.AutoMigrate(tables...)
	if err != nil {
		fmt.Println("Create table failed: ", err)
		panic(err)
	}
}
func setIndexes() error {
	type IndexInfo struct {
		KeyName    string
		SeqInIndex int
		ColumnName string
	}

	var err error
	var results []IndexInfo
	//检查是否存在索引，不存在则创建FULLTEXT索引
	check := fmt.Sprintf("SHOW INDEX FROM `video` WHERE COLUMN_NAME IN ('title','description') AND Index_type = 'FULLTEXT'")
	DB.Raw(check).Scan(&results)
	if len(results) == 0 {
		err = DB.Exec("CREATE FULLTEXT INDEX idx_user_name ON user(user_name) with parser ngram;").Error
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate key name") {
				return nil
			}
			return err
		}
		err = DB.Exec("CREATE FULLTEXT INDEX idx_fulltext_title_description ON video(title,description) with parser ngram;").Error
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate key name") {
				return nil
			}
			return err
		}
	}

	return err
}
func initMysql() {
	createDatabase()
	connectMysql()
	checkAndCreateTable()
	err := setIndexes()
	if err != nil {
		logf.WriteErrLog("setIndexes", err.Error())
		panic(err)
	}
	// 设置锁等待超时时间为 10 秒
	if err := DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
		panic(err)
	}
	logf.WriteInfoLog("initMysql", "Mysql init success")
}

func InitDBs() {
	initMysql()
	RDB = initRedisClient()
	MongoDB = initMongo()
	logf.WriteInfoLog("InitDBs", "All DBs init success")
}
