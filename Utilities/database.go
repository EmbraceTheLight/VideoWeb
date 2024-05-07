package Utilities

// This file is used to delete soft deleted records from the database periodically.

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	RelationshipSets "VideoWeb/DAO/RelationshipSets"
	"VideoWeb/Utilities/logf"
	"VideoWeb/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type errInDelete struct {
	id  int
	err error
}

const (
	barrage = iota + 1
	videos
	users
	comments
	favorite
	videoHistory
	userHistory
	favoriteVideo
	followed
	follows
)

var idToFunc = make(map[int]string)
var errChannel = make(chan *errInDelete, 11)
var db *gorm.DB

func initNecessary() error {
	var err error
	MySQLConf := config.GetConfig().DBConf.MySQLConf
	DBConnection := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s", MySQLConf.User, MySQLConf.Password, MySQLConf.Host, MySQLConf.Port, MySQLConf.Database, MySQLConf.Charset)
	db, err = gorm.Open(mysql.Open(DBConnection), &gorm.Config{})
	if err != nil {
		fmt.Println("Open database failed: ", err)
		return err
	}
	if err := db.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout:", err)
		return err
	}
	idToFunc[barrage] = "HDBarrage"
	idToFunc[videos] = "HDVideos"
	idToFunc[users] = "HDUsers"
	idToFunc[comments] = "HDComments"
	idToFunc[favorite] = "HDFavorites"
	idToFunc[videoHistory] = "HDVideoHistory"
	idToFunc[userHistory] = "HDUserHistory"
	idToFunc[favoriteVideo] = "HDFavoriteVideo"
	idToFunc[followed] = "HDFollowed"
	idToFunc[follows] = "HDFollows"
	return nil
}

// HardDelete 定时删除已经软删除的记录，节省空间
func HardDelete() {
	err := initNecessary()
	if err != nil {
		fmt.Println("[HardDelete] err:", err)
		return
	}
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
				logf.WriteErrLog(idToFunc[err1.id], err1.err.Error())
			} else {
				logf.WriteInfoLog(idToFunc[err1.id], "deleted records successfully.")
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
	go HDVideoHistory()
	go HDUserHistory()
	go HDFavoriteVideo()
	go HDFollowed()
	go HDFollows()
}

func HDBarrage() {
	err := db.Unscoped().Delete(&EntitySets.Barrage{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  barrage,
		err: err,
	}
}

func HDVideos() {
	err := db.Unscoped().Delete(&EntitySets.Video{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  videos,
		err: err,
	}
}

func HDUsers() {
	err := db.Unscoped().Delete(&EntitySets.User{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  users,
		err: err,
	}

}
func HDComments() {
	err := db.Unscoped().Delete(&EntitySets.Comments{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  comments,
		err: err,
	}
}

func HDFavorites() {
	err := db.Unscoped().Delete(&EntitySets.Favorites{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  favorite,
		err: err,
	}
}

func HDVideoHistory() {
	err := db.Unscoped().Delete(&EntitySets.VideoHistory{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  videoHistory,
		err: err,
	}
}

func HDUserHistory() {
	err := db.Unscoped().Delete(&EntitySets.SearchHistory{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  userHistory,
		err: err,
	}
}

func HDFavoriteVideo() {
	err := db.Unscoped().Delete(&RelationshipSets.FavoriteVideo{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  favoriteVideo,
		err: err,
	}
}

func HDFollowed() {
	err := db.Unscoped().Delete(&RelationshipSets.UserFollowed{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  followed,
		err: err,
	}
}

func HDFollows() {
	err := db.Unscoped().Delete(&RelationshipSets.UserFollows{}, "deleted_at IS NOT NULL").Error
	errChannel <- &errInDelete{
		id:  follows,
		err: err,
	}
}
