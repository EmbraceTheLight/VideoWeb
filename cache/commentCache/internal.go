package commentCache

import (
	"VideoWeb/cache"
	"context"
)

func addToCommentCache(ctx context.Context, key string, values ...interface{}) error {
	var mp = cache.MakeMap(values...)
	return cache.HSetWithRetry(
		ctx, key,
		cache.DefaultTry, cache.DefaultSleep, cache.CommentExpireTime,
		mp,
	)
}
