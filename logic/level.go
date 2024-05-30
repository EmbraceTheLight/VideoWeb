package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/define"
	"VideoWeb/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// AddExpForLogin 每日登录添加经验
func AddExpForLogin(c *gin.Context, uid int64, db *gorm.DB) error {
	var err error
	defer Utilities.DeferFunc(c, err, "AddExpForLogin")

	us := strconv.FormatInt(uid, 10) + "Login"
	//查询用户Level
	level, err := EntitySets.GetLevelRecordByUserID(db, uid)
	if err != nil {
		return err
	}

	//检查redis中是否有今天的登录记录，如果有，则不再添加经验
	_, err = DAO.RDB.Get(c, us).Result()
	if err == nil {
		return nil
	} else {
		if helper.CheckRedisErrorNil(err) {
			level.AddExp(define.ExpLoginOneDay)
			DAO.RDB.Set(c, us, 1, Utilities.GetTomorrowTime()) //设置redis记录,过期时间为明天零时
			err = EntitySets.SaveLevelRecords(db, level)
			return err
		} else { //查询redis出现错误
			return err
		}
	}
}

// AddExpForThrowShells 每日投掷shells添加经验
func AddExpForThrowShells(c *gin.Context, uid int64, shells int, db *gorm.DB) error {
	var err error
	defer Utilities.DeferFunc(c, err, "AddExpForThrowShells")

	us := strconv.FormatInt(uid, 10) + "shells"
	//查询用户Level
	level, err := EntitySets.GetLevelRecordByUserID(db, uid)
	if err != nil {
		return err
	}

	ret, err := DAO.RDB.Get(c, us).Result()
	if err == nil { //查询到了相关记录
		shellsHasThrown, _ := strconv.Atoi(ret)
		if shellsHasThrown >= define.LimitShellsPerDay {
			return nil
		}
		if shellsHasThrown+shells < define.LimitShellsPerDay {
			level.AddExp(shells * define.ExpEachShellThrow)
			DAO.RDB.IncrBy(c, us, int64(shells))
		} else {
			level.AddExp((define.LimitShellsPerDay - shellsHasThrown) * define.ExpEachShellThrow)
			DAO.RDB.IncrBy(c, us, int64(shells))
		}
		err = EntitySets.SaveLevelRecords(db, level)
		return err
	} else {
		if helper.CheckRedisErrorNil(err) { //没有相关记录
			level.AddExp(shells * define.ExpEachShellThrow)
			DAO.RDB.Set(c, us, shells, Utilities.GetTomorrowTime()) //设置redis记录,过期时间为明天零时
			err = EntitySets.SaveLevelRecords(db, level)
			return err
		} else { //查询redis出现错误
			return err
		}
	}
}

// AddExpForGainShells 获得shells添加经验
func AddExpForGainShells(c *gin.Context, uid int64, shells int, db *gorm.DB) error {
	var err error
	defer Utilities.DeferFunc(c, err, "AddExpForGAINShells")

	//查询用户Level
	level, err := EntitySets.GetLevelRecordByUserID(db, uid)
	if err != nil {
		return err
	}
	level.AddExp(shells * define.ExpEachShellGain)
	err = EntitySets.SaveLevelRecords(db, level)
	return err
}

// AddExpForUploadVideo 上传视频添加经验
func AddExpForUploadVideo(c *gin.Context, uid int64, db *gorm.DB) error {
	var err error
	defer Utilities.DeferFunc(c, err, "AddExpForShareVideo")
	level, err := EntitySets.GetLevelRecordByUserID(db, uid)
	if err != nil {
		return err
	}
	level.AddExp(define.ExpEachUploadVideo)
	err = EntitySets.SaveLevelRecords(db, level)
	return err
}
