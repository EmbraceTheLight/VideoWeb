package DAO

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func newClient() *redis.Client {
	RedisConf := config.GetConfig().DBConf.RedisConf
	fmt.Println("REDISCONF:", RedisConf)
	return redis.NewClient(
		&redis.Options{
			Addr:     RedisConf.Host + ":" + RedisConf.Port,
			Password: RedisConf.Password,
			DB:       0,
		})
}

func CreateDatabase() {
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	dbConnection := fmt.Sprintf("%s:%s@(%s:%s)/mysql?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Charset)
	db, _ := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	db.Exec("CREATE DATABASE  IF NOT EXISTS VideoWeb")
}

func InitDB() (err error) {
	CreateDatabase()
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	DBConnection := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s",
		MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Database, MySQLConf.Charset)
	RDB = newClient()
	DB, err = gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		return err
	}

	if exist := DB.Migrator().HasTable("Barrages"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Barrage{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Barrage{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("FavoriteVideo"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.FavoriteVideo{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&RelationshipSets.FavoriteVideo{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("Level"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Level{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Level{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("Messages"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Message{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Message{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("Tags"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Tags{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Tags{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("Comments"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Comments{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Comments{}):", err)
		}
	}
	if exist := DB.Migrator().HasTable("UserFollowed"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollowed{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&RelationshipSets.UserFollowed{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("UserFollows"); !exist {
		err = DB.Debug().AutoMigrate(&RelationshipSets.UserFollows{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&RelationshipSets.UserFollows{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("Video"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Video{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Video{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("Favorites"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.Favorites{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.Favorites{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("Users"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.User{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.User{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("VideoHistory"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.VideoHistory{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(EntitySets.VideoHistory{}:", err)
		}
	}
	if exist := DB.Migrator().HasTable("SearchHistory"); !exist {
		err = DB.Debug().AutoMigrate(&EntitySets.SearchHistory{})
		if err != nil {
			fmt.Println("Err in AutoMigrate(&EntitySets.SearchHistory{}:", err)
		}
	}
	//// 设置锁等待超时时间为 10 秒
	if err := DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
		return err
	}
	return err
}
