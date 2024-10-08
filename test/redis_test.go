package test

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/Utilities"
	"VideoWeb/cache"
	"VideoWeb/cache/commentCache"
	"VideoWeb/cache/userCache"
	"VideoWeb/cache/videoCache"
	"VideoWeb/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

func init() {
	config.InitConfig("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWeb\\VideoWeb\\config\\config.yaml")
	DAO.InitDBs()
}

func makeUserInfo(userID int64) error {
	uc := userCache.MakeUserCache()
	ctx := context.Background()
	return uc.MakeUserinfo(ctx, userID)
}
func TestMakeUserInfo(t *testing.T) {
	var userID1 int64 = 63216299135045 // 该用户没有粉丝表
	require.NoError(t, makeUserInfo(userID1))

	var userID2 int64 = 63354558959685 // 该用户没有关注表
	require.NoError(t, makeUserInfo(userID2))

	var userID3 int64 = 63358285594693 // 该用户没有评论表，关注表和粉丝表
	require.NoError(t, makeUserInfo(userID3))
}

func getUserBasic(ctx context.Context, userID int64) (mp map[string]string, err error) {
	mp, err = userCache.GetUserBasicInfo(ctx, userID)
	return
}
func getUserComments(ctx context.Context, userID int64) (ucIDs []string, err error) {
	ucIDs, err = cache.SMembers(ctx, strconv.FormatInt(userID, 10)+"_comments")
	return
}
func TestGetUserAndCommentsInfo(t *testing.T) {
	var userID int64 = 63216299135045
	ctx := context.Background()
	ub, err := getUserBasic(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "63216299135045", ub["user_id"])

	ucIDs, err := getUserComments(ctx, userID)
	require.NoError(t, err)

	for k, v := range ub {
		if k == "avatar" {
			//a := []byte(v)
			continue
		}
		fmt.Println("k:", k, "	v:", v)
	}
	for _, ucID := range ucIDs {
		uc, err := cache.HGetAll(ctx, strconv.FormatInt(Utilities.String2Int64(ucID), 10), cache.CommentExpireTime)
		require.NoError(t, err)

		fmt.Println("ucID:", ucID)
		for k, v := range uc {
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}

	var notExistUser int64 = 000000
	mp, err := getUserBasic(ctx, notExistUser)
	if mp["empty"] == "1" {
		err = gorm.ErrRecordNotFound
	}
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestGetUserInfos(t *testing.T) {
	ctx := context.Background()
	var userIDs = []int64{
		63216299135045,
		63354558959685,
		63358285594693,
		000000000, // this user is not exist
	}
	ubs, err := userCache.GetUsersBasicInfo(ctx, userIDs)
	println("ubslen:", len(ubs))
	if ubs[len(userIDs)-1]["empty"] == "1" {
		err = gorm.ErrRecordNotFound
	}
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)

	for _, ub := range ubs {
		for k, v := range ub {
			if k == "avatar" {
				//a := []byte(v)
				fmt.Println("k: ", k, "	v: avatar")
				continue
			}
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
		fmt.Println()
	}
}

func TestGetUserComments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID int64 = 63216299135045
	var ucIDs []int64
	err := DAO.DB.Model(&EntitySets.Comments{}).Where("user_id = ?", userID).Select("comment_id").Find(&ucIDs).Error
	require.NoError(t, err)
	//ucIds, err :=

	require.NoError(t, err)
	comments, err := commentCache.GetCommentsInfo(ctx, ucIDs)
	require.NoError(t, err)

	for _, comment := range comments {
		for k, v := range comment {
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
	}
	fmt.Println()

}

func makeVideoInfo(videoID int64) error {
	vc := videoCache.MakeVideoCache()
	ctx := context.Background()
	return vc.MakeVideoInfo(ctx, videoID)
}
func TestMakeVideoInfo(t *testing.T) {
	var videoID int64 = 52826949386309
	err := makeVideoInfo(videoID)
	require.NoError(t, err)

}

func TestGetVideoInfo(t *testing.T) {
	var videoID int64 = 63216465526853
	err := makeVideoInfo(videoID)
	require.NoError(t, err)

	ctx := context.Background()
	//video basic info
	vbs, err := videoCache.GetVideosBasicInfo(ctx, videoID)
	require.NoError(t, err)
	fmt.Println("↓↓↓↓↓video basic info↓↓↓↓↓")
	for _, v := range vbs {
		for k, v := range v {
			fmt.Println("k:", k, "	v:", v)
		}
	}

	//tags
	fmt.Println()
	vt, err := videoCache.GetTagsInfo(ctx, videoID)
	require.NoError(t, err)
	fmt.Println("↓↓↓↓↓video tags↓↓↓↓↓")
	for i, tag := range vt {
		fmt.Printf("tag%d: %s\n", i+1, tag)
	}
	fmt.Println()

	//comments
	vc, err := videoCache.GetAllVideoCommentsInfo(ctx, videoID)
	require.NoError(t, err)
	fmt.Println("↓↓↓↓↓video comments↓↓↓↓↓")
	for _, comment := range vc {
		for k, v := range comment {
			fmt.Println("k:", k, "	v:", v)
		}
		println()
		println()
		println()
	}

	//barrages
	barrages, err := videoCache.GetBarragesInfo(ctx, videoID)
	require.NoError(t, err)
	fmt.Println("↓↓↓↓↓video barrages↓↓↓↓↓")
	for _, barrage := range barrages {
		for k, v := range barrage {
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
		fmt.Println()
	}
}

// 测试通过redis交集获取用户点赞的评论
func TestGetUserLikedComments(t *testing.T) {
	var userID int64 = 52826422661189
	var videoID int64 = 52826949386309

	err := makeVideoInfo(videoID)
	require.NoError(t, err)

	ctx := context.Background()
	err = userCache.MakeUserLikedComments(ctx, userID, videoID)
	require.NoError(t, err)

	comments, err := videoCache.GetUserLikedCommentsInfo(ctx, videoID, userID)
	require.NoError(t, err)

	for _, comment := range comments {
		fmt.Println("comment_id:", comment["comment_id"])
		for k, v := range comment {
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
	}
}

func TestMakeAndGetVideoZSetInfos(t *testing.T) {
	ctx := context.Background()
	tmp, err := videoCache.MakeAllVideosZSet(ctx)
	require.NoError(t, err)

	err = videoCache.SaveVideosZSet(ctx, "all_videos", tmp...)
	require.NoError(t, err)

	videoIDs, err := videoCache.GetVideoZSetInfo(ctx, "all_videos", 0, 4)
	require.NoError(t, err)

	videos, err := videoCache.GetVideosBasicInfo(ctx, videoIDs...)
	require.NoError(t, err)

	for _, video := range videos {
		fmt.Println("video:", video["video_id"])
		for k, v := range video {
			fmt.Println("k:", k, "	v:", v)
		}
		fmt.Println()
	}
}
