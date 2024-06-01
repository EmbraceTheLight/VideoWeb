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
	MongoDB *mongo.Database
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

func initMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("[initMongo] Error connecting to MongoDB:", err)
		return nil
	}
	db := client.Database("IM")
	logf.WriteInfoLog("initMongo", "Mongo init success")
	return db
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
	if exist := DB.Migrator().HasTable("barrages"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Barrage{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Barrage{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("favorite_video"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.FavoriteVideo{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.FavoriteVideo{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_level"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Level{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Level{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("tags"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Tags{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Tags{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("comments"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Comments{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Comments{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_followed"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollowed{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserFollowed{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_follows"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollows{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserFollows{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("video"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Video{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Video{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("favorites"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Favorites{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Favorites{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("follow_list"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.FollowList{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.FollowList{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.User{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.User{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("video_history"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.VideoHistory{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.VideoHistory{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("search_history"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.SearchHistory{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.SearchHistory{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_video"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserVideo{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserVideo{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_liked_comments"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserLikedComments{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserComments{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("user_Disliked_comments"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserDislikedComments{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserDislikedComments{}):", err))
			panic(err)
		}
	}
}
func setIndexes() error {
	err := DB.Exec("CREATE FULLTEXT INDEX idx_user_name ON user(user_name) with parser ngram;").Error
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
