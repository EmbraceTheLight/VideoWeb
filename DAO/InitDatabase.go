package DAO

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

var userName string = "root"
var password string = "213103"
var ipPort string = "127.0.0.1:3306"
var dataBase string = "VideoWeb"
var charset string = "utf8mb4"

func InitMySQL() (err error) {
	dbc := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", userName, password, ipPort, dataBase, charset)
	//dbc := "root:213103@(127.0.0.1:3306)/VideoWeb?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s"
	DB, err = gorm.Open(mysql.Open(dbc), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		return err
	}
	//
	//if DB.Debug().AutoMigrate(&RelationshipSets.FavoriteVideo{}) != nil {
	//	fmt.Println("err in AutoMigrate(&FavoriteVideo{}): ", err)
	//}

	// 设置锁等待超时时间为 10 秒
	if err := DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
		return err
	}

	return err
}
