package DAO

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

var userName = "root"
var password = "213103"
var ipPort = "127.0.0.1:3306"
var dataBase = "VideoWeb"
var charset = "utf8mb4"
var redisIpPort = "127.0.0.1:6379"

func newClient() *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     redisIpPort,
			Password: "",
			DB:       0,
		})
}

func InitMySQL() (err error) {
	dbc := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", userName, password, ipPort, dataBase, charset)
	RDB = newClient()

	DB, err = gorm.Open(mysql.Open(dbc), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		return err
	}
	//
	//if DB.Debug().AutoMigrate(&RelationshipSets.FavoriteVideo{}) != nil {
	//	fmt.Println("err in AutoMigrate(&FavoriteVideo{}): ", err)
	//}
	//if DB.Debug().AutoMigrate(&DAO.Barrage{}) != nil {
	//	fmt.Println("err int Barrage")
	//}
	//if DB.Debug().AutoMigrate(&DAO.Video{}) != nil {
	//	fmt.Println("err in AutoMigrate(&Video{}): ", err)
	//}
	//
	//// 设置锁等待超时时间为 10 秒
	if err := DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
		return err
	}

	return err
}
