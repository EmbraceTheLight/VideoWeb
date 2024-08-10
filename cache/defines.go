// Package cache defines 定义基本变量和常量
package cache

import "time"

var (
	DefaultSleep = 50 * time.Millisecond //重试间隔时间
	DefaultTry   = 5                     //最大重试次数

	NotFoundExpireTime  = 5 * time.Minute  //设置的值为空缓存的过期时间。空缓存是由于Mysql数据库中无对应数据时，防止缓存穿透而设定。
	UserExpireTime      = 60 * time.Minute //用户缓存过期时间
	CommentExpireTime   = 24 * time.Hour   //评论缓存过期时间
	VideoExpireTime     = 24 * time.Hour   //视频缓存过期时间
	OperationExpireTime = 3 * time.Second  //缓存操作过期时间
	DeleteTime          = 1 * time.Second  //指定多少秒后进行延时双删
)

const (
	VideoZSetKey = "video_zset"      //视频排行榜zset key
	ULCSfx       = "_liked_comments" //用户喜欢的评论后缀
	CommentSfx   = "_comments"       //评论后缀
	BarrageSfx   = "_barrages"       //弹幕后缀
	TagSfx       = "_tags"           //标签后缀
	FldSfx       = "_followed"
	FlsSfx       = "_follows"
	SearchSfx    = "_searches"
	WatchSfx     = "_watches"
)

type ZSetValue interface {
	GetScore() float64
	GetValue() interface{}
}
