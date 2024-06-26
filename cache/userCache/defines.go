package userCache

import (
	"context"
)

// UserBasic cache the user's basic information
type UserBasic struct {
	Userinfo map[string]any
}

func (ub *UserBasic) Set(k string, v any) {
	ub.Userinfo[k] = v
}
func (ub *UserBasic) Get(k string) any {
	return ub.Userinfo[k]
}
func (ub *UserBasic) Del(k string) {
	delete(ub.Userinfo, k)
}

// UserComments cache the user comments' IDs
type UserComments struct {
	key          string
	UserComments []int64
}

func (uc *UserComments) Set(v any) {
	set(uc.UserComments, v.(int64))
}
func (uc *UserComments) GetAll() []int64 {
	return uc.UserComments
}
func (uc *UserComments) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.userCache.UserComments->GetOne", uc.key, member)
}
func (uc *UserComments) Del(v any) {
	del(uc.UserComments, v.(int64))
}

// UserFollows cache the user's following-users IDs
type UserFollows struct {
	key         string
	UserFollows []int64
}

func (ufs *UserFollows) Set(v any) {
	set(ufs.UserFollows, v.(int64))
}
func (ufs *UserFollows) GetAll() []int64 {
	return ufs.UserFollows
}
func (ufs *UserFollows) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.userCache.UserFollows->GetOne", ufs.key, member)
}
func (ufs *UserFollows) Del(v any) {
	del(ufs.UserFollows, v.(int64))
}

// UserFollowed cache the user's followed-users IDs
type UserFollowed struct {
	key          string
	UserFollowed []int64
}

func (ufd *UserFollowed) Set(v any) {
	set(ufd.UserFollowed, v.(int64))
}
func (ufd *UserFollowed) GetAll() []int64 {
	return ufd.UserFollowed
}
func (ufd *UserFollowed) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.userCache.UserFollowed->GetOne", ufd.key, member)
}
func (ufd *UserFollowed) Del(v any) {
	del(ufd.UserFollowed, v.(int64))
}

// UserVideo cache the user's upload videos
type UserVideo struct {
	key       string
	UserVideo []int64
}

// UserSearch cache the user's search history
type UserSearch struct {
	key        string
	UserSearch []string
}

// UserWatch cache the user's watched videos
type UserWatch struct {
	key       string
	UserWatch []int64
}

// UserCache user info cache
type UserCache struct {
	Ub  UserBasic
	Uc  UserComments
	Ufs UserFollows
	Ufd UserFollowed
	UV  UserVideo
	US  UserSearch
	UW  UserWatch
}

//func (Uc *UserComments) MarshalBinary() (data []byte, err error) {
//	return json.Marshal(Uc.UserComments)
//}
//
//func (Uc *UserComments) UnmarshalBinary(data []byte) error {
//	return json.Unmarshal(data, &Uc.UserComments)
//}
