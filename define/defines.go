package define

import (
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
	"time"
)

var (
	Store = base64Captcha.NewMemoryStore(base64Captcha.GCLimitNumber, Expired)

	DefaultPage = 1
	DefaultSize = 20
	Expired     = time.Minute * 10 //过期时间。单位：秒。用于验证验证码是否过期
)

const (
	KiB = int64(1024)
	MiB = 1024 * KiB
	GiB = 1024 * MiB
)

// BaseTime 时间戳基准点，2024-01-01 00:00:00, UTC 时间
const BaseTime int64 = 1704067200000000000 / 1e6

// Hot 视频热度相关
const (
	// AddHotEachView 每访问一次增加的热度
	AddHotEachView = 1

	// AddHotEachComment 每评论一次增加的热度
	AddHotEachComment = 3 * AddHotEachView

	// AddHotEachBarrage 每发表一次弹幕增加的热度
	AddHotEachBarrage = 3 * AddHotEachComment

	// AddHotEachLike 每点赞一次增加的热度
	AddHotEachLike = 10 * AddHotEachView

	// AddHotEachShell 每投一个贝壳增加的热度
	AddHotEachShell = 3 * AddHotEachView

	// AddHotEachFavorite 每收藏一次增加的热度
	AddHotEachFavorite = 75 * AddHotEachView

	// AddHotEachShare 每分享一次增加的热度
	AddHotEachShare = 50 * AddHotEachView
)

// Level 等级相关
const (
	ToLevel2 = 200
	ToLevel3 = 1500
	ToLevel4 = 4500
	ToLevel5 = 10800
	ToLevel6 = 28800

	ExpLoginOneDay     = 5
	ExpEachUploadVideo = 10
	ExpEachShellGain   = 1
	ExpEachShellThrow  = 5 //投递一次贝壳获得经验值。上限：一天5次
	ExpEachShare       = 5

	LimitShellsPerDay = 5 //每天投递贝壳获得经验的上限
	LimitSharesPerDay = 10
)

// MyModel 不包含 gorm.Model 中的默认id字段
type MyModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type GraphicCaptcha struct {
	ID     string `json:"ID"`
	B64lob string `json:"b64lob"`
	Answer string `json:"answer"`
}

type IPInfo struct {
	Query     string `json:"query"`     //待查询的IP
	Status    string `json:"status"`    //查询状态，如success
	Continent string `json:"continent"` //大洲	//Mobile  bool `json:"mobile"`
	Country   string `json:"country"`
	City      string `json:"city"`
	//Proxy   bool `json:"proxy"`
	//Hosting bool `json:"hosting"`
	//Lat           float64 `json:"lat"`       //纬度
	//Lon           float64 `json:"lon"`       //经度
	//ContinentCode string  `json:"continentCode"`
	//CountryCode   string  `json:"countryCode"`
	//Region        string  `json:"region"`
	//RegionName    string  `json:"regionName"`
	//District      string  `json:"district"`
	//Zip           string  `json:"zip"`
	//Timezone      string  `json:"timezone"`
	//Offset        string  `json:"offset"`
	//Currency      string  `json:"currency"`
	//ISP           string  `json:"ISP"`
	//Org           string  `json:"org"`
	//As            string  `json:"as"`
	//AsName        string  `json:"asName"`
}

// PicExtCheck 验证图片后缀名
var PicExtCheck = make(map[string]struct{})
var PictureSavePath = "./resources/pictures/"

// VideoExtCheck 验证视频后缀名
var VideoExtCheck = make(map[string]struct{})
var VideoSavePath = "./resources/videos/"
var FFProbe = "ffprobe"

func init() {
	//支持的图片的格式
	PicExtCheck[".jpg"] = struct{}{}
	PicExtCheck[".jpeg"] = struct{}{}
	PicExtCheck[".png"] = struct{}{}
	PicExtCheck[".jfif"] = struct{}{}

	//支持的视频的格式
	VideoExtCheck[".mp4"] = struct{}{}
	VideoExtCheck[".mov"] = struct{}{}
	VideoExtCheck[".avi"] = struct{}{}
	VideoExtCheck[".mkv"] = struct{}{}
	VideoExtCheck[".m4v"] = struct{}{}
	VideoExtCheck[".3gp"] = struct{}{}
	VideoExtCheck[".3g2"] = struct{}{}
}
