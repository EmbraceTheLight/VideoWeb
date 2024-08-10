package cache

import (
	"VideoWeb/DAO"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// MakeMap converts []any to map[string]any.
func MakeMap(values ...any) map[string]any {
	if len(values)%2 != 0 {
		return nil
	}
	mp := make(map[string]any, len(values))
	for i := 0; i < len(values); i += 2 {
		field, ok := values[i].(string)
		if field == "" || !ok {
			fmt.Println(222)
			return nil
		}
		fmt.Println("field:", field, "	value:", values[i+1])
		mp[field] = values[i+1]
	}
	return mp
}

// GetInfos gets detail infos from redis by keys.
func GetInfos(ctx context.Context, keys ...string) (ret []map[string]string, err error) {
	cmds := make([]*redis.MapStringStringCmd, len(keys))

	pipe := DAO.RDB.Pipeline()
	for i, key := range keys {
		cmds[i] = pipe.HGetAll(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("cache.common.GetInfos: %w", err)
	}

	for _, cmd := range cmds {
		ret = append(ret, cmd.Val())
	}
	return
}
