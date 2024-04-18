package Utilities

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type errInDelete struct {
	id  int
	err error
}

var idToFunc = make(map[int]string)
var errChannel = make(chan *errInDelete, 11)
var db *gorm.DB

func init() {
	DBConnection := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", DAO.UserName, DAO.Password, DAO.IpPort, DAO.DataBase, DAO.Charset)
	var err error
	db, err = gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
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
}

// HardDelete 定时删除已经软删除的记录，节省空间
func HardDelete() {
	ticker := time.NewTicker(5 * time.Hour)
	defer func() {
		ticker.Stop()
		close(errChannel)
	}()
	for {
		select {
		case <-ticker.C:
			HDHelper()
		case err1 := <-errChannel:
			if err1.err != nil {
				log.Printf("error in Function [%s]: %v\n", idToFunc[err1.id], err1.err)
			} else {
				log.Printf("Function [%s] deleted records successfully\n", idToFunc[err1.id])
			}
		}
	}
}
func HDHelper() {
	go HDBarrage()
	go HDVideos()
	go HDUsers()
	go HDComments()
	go HDFavorites()
	go HDMessages()
	go HDVideoHistory()
	go HDUserHistory()
	go HDFavoriteVideo()
	go HDFollowed()
	go HDFollows()
}

func HDBarrage() {
	err := db.Unscoped().Delete(&EntitySets.Barrage{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  1,
		err: err,
	}
}

func HDVideos() {
	err := db.Unscoped().Delete(&EntitySets.Video{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  2,
		err: err,
	}
}

func HDUsers() {
	err := db.Unscoped().Delete(&EntitySets.User{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  3,
		err: err,
	}
}

func HDComments() {
	err := db.Unscoped().Delete(&EntitySets.Comments{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  4,
		err: err,
	}
}

func HDFavorites() {
	err := db.Unscoped().Delete(&EntitySets.Favorites{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  5,
		err: err,
	}
}

func HDMessages() {
	err := db.Unscoped().Delete(&EntitySets.Message{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  6,
		err: err,
	}
}

func HDVideoHistory() {
	err := db.Unscoped().Delete(&EntitySets.VideoHistory{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  7,
		err: err,
	}
}

func HDUserHistory() {
	err := db.Unscoped().Delete(&EntitySets.SearchHistory{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  8,
		err: err,
	}
}

func HDFavoriteVideo() {
	err := db.Unscoped().Delete(&RelationshipSets.FavoriteVideo{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  9,
		err: err,
	}
}

func HDFollowed() {
	err := db.Unscoped().Delete(&RelationshipSets.UserFollowed{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  10,
		err: err,
	}
}

func HDFollows() {
	err := db.Unscoped().Delete(&RelationshipSets.UserFollows{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  11,
		err: err,
	}
}
