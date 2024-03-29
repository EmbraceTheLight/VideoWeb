package define

import (
	"gorm.io/gorm"
	"time"
)

var (
	DefaultPage       = "1"
	DefaultSize       = "20"
	Expired     int64 = 600 //过期时间。单位：秒。用于验证验证码是否过期
	Level1      int   = 2   //从0级升到1级所需经验值,下同
	Level2      int   = 4
	Level3      int   = 8
	Level4      int   = 16
	Level5      int   = 32
	Level6      int   = 64
	Level7      int   = 128
	Level8      int   = 256
	Level9      int   = 512
	Level10     int   = 1024
)

type VerificationData struct {
	Code      string
	Timestamp int64
}

// VerificationDataMap 与认证有关的Map
var VerificationDataMap = make(map[string]VerificationData)

type MyModel struct {
	// 不包含 gorm.Model 中的默认字段
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type IPInfo struct {
	Query         string  `json:"query"`     //待查询的IP
	Status        string  `json:"status"`    //查询状态，如success
	Continent     string  `json:"continent"` //大洲
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"` //纬度
	Lon           float64 `json:"lon"` //经度
	Timezone      string  `json:"timezone"`
	Offset        string  `json:"offset"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"ISP"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	AsName        string  `json:"asName"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
}

// PicExtCheck 验证图片后缀名
var PicExtCheck = make(map[string]struct{})
var PictureSavePath = ".\\resources\\pictures\\"

// VideoExtCheck 验证视频后缀名
var VideoExtCheck = make(map[string]struct{})
var VideoSavePath = "D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWeb\\resources\\Video\\"

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
