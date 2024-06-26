// Package commentCache 缓存单条评论的详情，与userCache以及videoCache配合使用
package commentCache

// CommentCache 存储单条评论的详情
type CommentCache struct {
	key      string
	comments map[string]any
}

func (c *CommentCache) GetKey() string {
	return c.key
}

func (c *CommentCache) GetMap() map[string]any {
	return c.comments
}

// NewCommentCache 创建一个新的 CommentCache
func NewCommentCache() *CommentCache {
	return &CommentCache{comments: make(map[string]any)}
}
