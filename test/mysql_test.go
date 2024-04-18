package test

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

type errInDelete struct {
	id  int
	err error
}

var idToFunc = make(map[int]string)
var errChannel = make(chan *errInDelete, 11)
var db *gorm.DB

//	func TestMain(m *testing.M) {
//		DBConnection := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", DAO.UserName, DAO.Password, DAO.IpPort, DAO.DataBase, DAO.Charset)
//		db, err := gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
//		fmt.Println("db in TestMain:", db)
//		if err != nil {
//			fmt.Println("Open database failed: ", err)
//		}
//		if err := db.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
//			fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
//		}
//		idToFunc[1] = "HDBarrage"
//		idToFunc[2] = "HDVideos"
//		idToFunc[3] = "HDUsers"
//		idToFunc[4] = "HDComments"
//		idToFunc[5] = "HDFavorites"
//		idToFunc[6] = "HDMessages"
//		idToFunc[7] = "HDVideoHistory"
//		idToFunc[8] = "HDUserHistory"
//		idToFunc[9] = "HDFavoriteVideo"
//		idToFunc[10] = "HDFollowed"
//		idToFunc[11] = "HDFollows"
//
//		t := new(testing.T)
//		TestHardDelete(t)
//	}
func TestHardDelete(t *testing.T) {
	DBConnection := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", DAO.UserName, DAO.Password, DAO.IpPort, DAO.DataBase, DAO.Charset)
	db, err := gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
	}
	if err := db.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
	}
	idToFunc[1] = "HDBarrage"
	idToFunc[2] = "HDVideos"
	idToFunc[3] = "HDUsers"
	idToFunc[4] = "HDComments"
	idToFunc[5] = "HDFavorites"
	idToFunc[6] = "HDMessages"
	idToFunc[7] = "HDVideoHistory"
	idToFunc[8] = "HDUserHistory"
	idToFunc[9] = "HDFavoriteVideo"
	idToFunc[10] = "HDFollowed"
	idToFunc[11] = "HDFollows"
	fmt.Println("db:", db)
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		close(errChannel)
	}()
	fmt.Println("db in HardDelete:", db)
	var cnt = 0
	for cnt < 5 {
		select {
		case <-ticker.C:
			err := db.Debug().Unscoped().Delete(&EntitySets.Favorites{}, "deleted_at IS NOT NULL").Error
			errChannel <- &errInDelete{
				id:  5,
				err: err,
			}
			cnt++
		case err1 := <-errChannel:
			if err1.err != nil {
				log.Printf("error in Function [%s]: %v\n", idToFunc[err1.id], err1.err)
			} else {
				log.Printf("Function [%s] deleted records successfully\n", idToFunc[err1.id])
			}
			cnt++
		}
	}
}

// HardDelete 定时删除已经软删除的记录，节省空间
func HardDelete() {

}
