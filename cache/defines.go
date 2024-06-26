// Package cache defines 定义基本变量和常量
package cache

import "time"

var (
	DefaultSleep = 50 * time.Millisecond //重试间隔时间
	DefaultTry   = 5                     //最大重试次数

	NotFoundExpireTime = 5 * time.Minute  //设置的值为空缓存的过期时间。空缓存是由于Mysql数据库中无对应数据时，防止缓存穿透而设定。
	UserExpireTime     = 60 * time.Minute //用户缓存过期时间
	CommentExpireTime  = 24 * time.Hour   //评论缓存过期时间
	VideoExpireTime    = 24 * time.Hour   //视频缓存过期时间
)
