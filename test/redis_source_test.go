package test

import (
	"context"
	"testing"
)

type cmdable func(ctx context.Context, cmd string) error

func (c cmdable) Test(ctx context.Context, cmd string) {
	println("In test ", cmd)
	_ = c(ctx, cmd)
}

func TestRedis(t *testing.T) {
	s := "Hello world"
	ctx := context.Background()
	var c cmdable
	c.Test(ctx, s)
}
