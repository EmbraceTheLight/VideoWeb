package userCache

// UserBasic cache the user's basic information
type UserBasic struct {
	Userinfo map[string]any
}

// UserComments cache the user comments' IDs
type UserComments struct {
	key          string
	UserComments []int64
}

// UserFollows cache the user's following-users IDs
type UserFollows struct {
	key         string
	UserFollows []int64
}

// UserFollowed cache the user's followed-users IDs
type UserFollowed struct {
	key          string
	UserFollowed []int64
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
	Uv  UserVideo
	Us  UserSearch
	Uw  UserWatch
}
