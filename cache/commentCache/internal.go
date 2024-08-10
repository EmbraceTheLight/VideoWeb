package commentCache

import (
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/cache"
	"context"
	"fmt"
	"strconv"
)

func addCommentInfo(ctx context.Context, prefix int64, comment *EntitySets.Comments) (err error) {
	err = MakeCommentInfos(ctx, comment)
	if err != nil {
		return err
	}

	err = cache.SAddWithRetry(
		ctx,
		strconv.FormatInt(prefix, 10)+"_comments",
		cache.DefaultTry,
		cache.DefaultSleep,
		cache.CommentExpireTime,
		comment.CommentID,
	)
	if err != nil {
		return fmt.Errorf("commentCache.AddCommentInfo::%w", err)
	}
	
	return nil
}
