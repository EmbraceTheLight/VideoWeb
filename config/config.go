package config

import (
	"VideoWeb/Utilities/logf"
	"VideoWeb/define"
	"VideoWeb/logrusLog"
	"bufio"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/shiena/ansicolor"
	"github.com/siruspen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"time"
)

//---------- MySQL Config ----------//

type MySQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

//---------- Redis Config ----------//

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

//---------- MongoDB Config ----------//

type Index struct {
	Fields map[string]int `yaml:"fields"` // key is field name, value is index order:1 is ascending, -1 is descending
	Type   string         `yaml:"type"`   // mongodb index type
}
type Collection struct {
	Name    string  `yaml:"name"`
	Indexes []Index `yaml:"indexes"`
}
type MongoConfig struct {
	Host        string       `yaml:"host"`
	Port        string       `yaml:"port"`
	Database    string       `yaml:"database"`
	Collections []Collection `yaml:"collections"`
}

//--------------------------------------------------------------------------//

type DBConfig struct {
	MySQLConf *MySQLConfig `yaml:"mysql"`
	MongoConf *MongoConfig `yaml:"mongo"`
	RedisConf *RedisConfig `yaml:"redis"`
}
type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"path"`
	TimeFormat string `yaml:"timeFormat"`
}

type LogFormat struct {
	logConf *LogConfig
}

type Config struct {
	DBConf  *DBConfig  `yaml:"database"`
	LogConf *LogConfig `yaml:"log"`
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

// path为空则默认读取config目录下的config.yaml文件
func parseConfig(path string) error {
	if path == "" {
		path = "./config/config.yaml"
	}
	//parse config file and return Config struct
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	return nil
}

func (format *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	//设置时区、时间
	tz, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("err in func Format:", err)
		return nil, err
	}
	entry.Time = entry.Time.In(tz)
	entry.Time.Format(format.logConf.TimeFormat)

	//自定义日志输出格式
	formatted := fmt.Sprintf("[%s] %s In Function [%s]: %s\n",
		entry.Level.String(),
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Data["function"],
		entry.Message,
	)
	return []byte(formatted), nil
}

func InitLog() error {
	logConf := GetConfig().LogConf
	err := os.MkdirAll("./logs", 0755)
	if err != nil {
		return err
	}

	//设置多路输出
	writer1 := ansicolor.NewAnsiColorWriter(os.Stdout)
	writer2, _ := rotatelogs.New(
		logConf.FilePath+".%Y%m%d",
		rotatelogs.WithLinkName(logConf.FilePath),
		rotatelogs.WithRotationCount(5),
		rotatelogs.WithRotationSize(300*define.MiB),
	)

	Output := io.MultiWriter(writer1, writer2)
	logrusLog.Log.SetOutput(Output)
	//设置日志格式
	logrusLog.Log.SetFormatter(&LogFormat{logConf: logConf})
	logf.WriteInfoLog("InitConfig", "Init Config Successfully.")
	return nil
}

func InitConfig(configPath string) {
	err := parseConfig(configPath)
	if err != nil {
		panic(err)
	}
	err = InitLog()
	if err != nil {
		panic(err)
	}
}
