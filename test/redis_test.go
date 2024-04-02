package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

func TestRedisSet(t *testing.T) {
	err := rdb.Set(ctx, "name", "mmmcskmc", time.Second*10).Err()
	if err != nil {
		println(err.Error())
	}
}

func TestRedisGet(t *testing.T) {
	v, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		println(err.Error())
	}
	println("value:", v)

}
