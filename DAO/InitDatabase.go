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
	db.Exec("CREATE DATABASE  IF NOT EXISTS VideoWeb")
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
	if exist := DB.Migrator().HasTable("Barrages"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Barrage{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Barrage{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("FavoriteVideo"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.FavoriteVideo{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.FavoriteVideo{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("Level"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Level{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Level{}):", err))
			panic(err)
		}
	}

	if exist := DB.Migrator().HasTable("Tags"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Tags{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Tags{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("Comments"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Comments{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Comments{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("UserFollowed"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollowed{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserFollowed{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("UserFollows"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollows{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&RelationshipSets.UserFollows{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("Video"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Video{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Video{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("Favorites"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Favorites{})
		if err != nil {
			logf.WriteErrLog("initMysql::checkAndCreateTable", fmt.Sprintln("Err in AutoMigrate(&EntitySets.Favorites{}):", err))
			panic(err)
		}
	}
	if exist := DB.Migrator().HasTable("Users"); !exist {
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
}
func initMysql() {
	createDatabase()
	connectMysql()
	checkAndCreateTable()
	//// 设置锁等待超时时间为 10 秒
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
