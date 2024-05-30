package helper

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

func CheckRedisErrorNil(err error) bool {
	if errors.Is(err, redis.Nil) {
		return true
	}
	return false
}
