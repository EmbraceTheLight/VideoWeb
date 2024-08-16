package DAO

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities/logf"
	"VideoWeb/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
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

/*
以下摘自mongo-driver官方文档：

	创建客户端的最佳方式是使用连接池。连接池可以重用现有的连接，而不是每次都创建新的连接。
	为每个进程创建一个客户端，并在所有操作中重复使用。为每个请求创建一个新客户端是个常见的错误，效率非常低。
	每个Client实例都有一个内置连接池。 连接池按需打开套接字以支持并发 MongoDB 操作或 goroutine
*/
func initMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.GetConfig().DBConf.MongoConf
	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)).
			SetMaxPoolSize(150))
	if err != nil {
		log.Println("[initMongoClient] Error connecting to MongoDB: ", err)
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("[initMongoClient] Error pinging MongoDB: ", err)
		panic(err)
	}

	logf.WriteInfoLog("initMongoClient", "Mongo init success")
	return client
}
func initMongoCollections(client *mongo.Client) {
	db := client.Database(config.GetConfig().DBConf.MongoConf.Database)
	collections := config.GetConfig().DBConf.MongoConf.Collections
	fmt.Println("collections: ", collections)
	for _, collection := range collections {
		c := db.Collection(collection.Name)
		var indexModel mongo.IndexModel

		// handle indexes of a collection
		for _, index := range collection.Indexes {
			keys := bson.D{}
			// handle keys and orders of an index
			for key, order := range index.Fields {
				keys = append(keys, bson.E{Key: key, Value: order})
			}
			// create index model
			indexModel = mongo.IndexModel{
				Keys:    keys,
				Options: options.Index().SetUnique(index.Type == "unique"),
			}
		}

		fmt.Println("keys: ", indexModel.Keys)
		_, err := c.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			logf.WriteErrLog("initMongoCollections", err.Error())
			log.Fatalf("Failed to create index: %v", err)
		}
	}
}
func initMongo() {
	MongoDB = initMongoClient()
	initMongoCollections(MongoDB)
}

func createDatabase() {
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	dbConnection := fmt.Sprintf("%s:%s@(%s:%s)/mysql?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Charset)
	db, _ := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	db.Exec("CREATE DATABASE  IF NOT EXISTS video_web")
}
func connectMysql() *gorm.DB {
	var err error
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	DBConnection := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Database, MySQLConf.Charset)
	db, err := gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		panic(err)
	}
	return db
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
		&RelationshipSets.UserDislikedComments{},
	}
	err = DB.AutoMigrate(tables...)
	if err != nil {
		fmt.Println("Create table failed: ", err)
		panic(err)
	}
}
func setFulltextIndexes() error {
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
	// First, create the database
	// Then, connect to the database and create tables
	createDatabase()
	DB = connectMysql()
	checkAndCreateTable()
	err := setFulltextIndexes()
	if err != nil {
		logf.WriteErrLog("setFulltextIndexes", err.Error())
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
	initMongo()
	logf.WriteInfoLog("InitDBs", "All DBs init success")
}
