package cache

import (
	"VideoWeb/DAO"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type HashMap interface {
	GetKey() string
	GetMap() map[string]any
}

// SetNilHash sets a nil value of a hashmap to show that the key is not exist
func SetNilHash(ctx context.Context, key string) (err error) {
	_, err = DAO.RDB.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		var err1 error
		err1 = pipe.HSet(ctx, key, "empty", true).Err()
		if err1 != nil {
			return err1
		}
		err1 = pipe.Expire(ctx, key, NotFoundExpireTime).Err()
		return err1
	})
	if err != nil {
		return fmt.Errorf("SetNilHash::RDB.Pipelined --> %w", err)
	}
	return nil
}

// HSetWithRetry sets a hashmap with retry and interval time
func HSetWithRetry(
	ctx context.Context,
	key string,
	retryCount int, sleep time.Duration,
	expire time.Duration,
	values map[string]any) (err error) {
	err = DAO.RDB.HSet(ctx, key, values).Err()
	for err != nil {
		err = DAO.RDB.HSet(ctx, key, values).Err()
		retryCount--
		if err == nil || retryCount == 0 {
			break
		}
		time.Sleep(sleep)
	}

	if err != nil {
		return err
	}
	err = DAO.RDB.Expire(ctx, key, expire).Err()
	return
}

// HSets sets many hashmap with retry and interval time
func HSets(
	ctx context.Context,
	expire time.Duration,
	values ...HashMap) (err error) {

	pipe := DAO.RDB.Pipeline()

	for _, value := range values {
		pipe.HSet(ctx, value.GetKey(), value.GetMap())
		pipe.Expire(ctx, value.GetKey(), expire)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		return fmt.Errorf("cache.util.HSets: %w", err)
	}

	return
}

// HGetAll gets all values of a hashmap
func HGetAll(ctx context.Context, key string, expire time.Duration) (mp map[string]string, err error) {
	mp, err1 := DAO.RDB.HGetAll(ctx, key).Result()
	if err1 != nil {
		return nil, fmt.Errorf("cache.util.HGetAll: %w", err1)
	}
	if _, ok := mp["empty"]; ok {
		return nil, fmt.Errorf("cache.util.HGetAll: %w", errors.New("hash is not exist"))
	}

	err2 := DAO.RDB.Expire(ctx, key, expire).Err()
	if err2 != nil {
		return nil, fmt.Errorf("cache.util.HGetAll: %w", err2)
	}

	return
}

// SAddWithRetry sets values of a hashmap with retry and interval time
func SAddWithRetry(
	ctx context.Context,
	key string,
	retryCount int,
	sleep time.Duration,
	expire time.Duration,
	values ...any,
) (err error) {
	err = DAO.RDB.SAdd(ctx, key, values...).Err()
	for err != nil {
		err = DAO.RDB.SAdd(ctx, key, values...).Err()
		retryCount--
		if err == nil || retryCount == 0 {
			break
		}
		time.Sleep(sleep)
	}
	if err != nil {
		return fmt.Errorf("cache.SAddWithRetry: %w", err)
	}

	err = DAO.RDB.Expire(ctx, key, expire).Err()
	if err != nil {
		return fmt.Errorf("cache.SAddWithRetry: %w", err)
	}

	return nil
}

// ZAddWithRetry sets values of a ZSet with retry and interval time
func ZAddWithRetry(
	ctx context.Context,
	key string,
	retryCount int,
	sleep time.Duration,
	expire time.Duration,
	values ...redis.Z,
) (err error) {
	err = DAO.RDB.ZAdd(ctx, key, values...).Err()
	for err != nil {
		err = DAO.RDB.ZAdd(ctx, key, values...).Err()
		retryCount--
		if err == nil || retryCount == 0 {
			break
		}
		time.Sleep(sleep)
	}
	if err != nil {
		return fmt.Errorf("cache.ZAddWithRetry: %w", err)
	}

	err = DAO.RDB.Expire(ctx, key, expire).Err()
	if err != nil {
		return fmt.Errorf("cache.ZAddWithRetry: %w", err)
	}

	return nil
}

// AddToZSet adds values to redis zset
func AddToZSet(ctx context.Context, key string, expire time.Duration, values ...ZSetValue) (err error) {
	ZSetValues := make([]redis.Z, len(values))
	for i, v := range values {
		ZSetValues[i] = redis.Z{Score: v.GetScore(), Member: v.GetValue()}
	}
	err = ZAddWithRetry(
		ctx,
		key,
		DefaultTry, DefaultSleep, UserExpireTime,
		ZSetValues...,
	)
	if err != nil {
		return fmt.Errorf("cache.util.AddToZSet::%w", err)
	}

	err = DAO.RDB.Expire(ctx, key, expire).Err()
	if err != nil {
		return fmt.Errorf("cache.util.AddToZSet: %w", err)
	}
	return nil
}

// SMembers gets all members of a set
func SMembers(ctx context.Context, key string) (members []string, err error) {
	members, err = DAO.RDB.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("cache.util.SMembers: %w", err)
	}

	err = DAO.RDB.Expire(ctx, key, 60*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("cache.util.SMembers: %w", err)
	}
	return
}

// SIsMember checks if a member is in a set
func SIsMember(ctx context.Context, key, member string) (err error) {
	exist, err := DAO.RDB.SIsMember(ctx, key, member).Result()
	switch {
	case err != nil:
		return fmt.Errorf("cache.util.SIsMember : %w", err)
	case !exist:
		return fmt.Errorf("cache.util.SIsMember : %w", errors.New("member is not exist"))
	default:
		return nil
	}
}

// SInter gets the intersection of set1 and set2
func SInter(ctx context.Context, key1, key2 string) (res []string, err error) {
	res, err = DAO.RDB.SInter(ctx, key1, key2).Result()
	if err != nil {
		return nil, fmt.Errorf("cache.util.SInter: %w", err)
	}
	return
}

// LPush pushes values to the left of a list
func LPush(ctx context.Context, key string, values ...any) (err error) {
	DAO.RDB.Pipeline()
	err = DAO.RDB.LPush(ctx, key, values...).Err()

	if err != nil {
		return fmt.Errorf("::cache.util.LPush : %w", err)
	}
	return nil
}
