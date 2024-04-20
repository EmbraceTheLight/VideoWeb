package config

import (
	"VideoWeb/logrusLog"
	"bufio"
	"fmt"
	"github.com/siruspen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"time"
)

type MySQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
type DBConfig struct {
	MySQLConf *MySQLConfig `yaml:"mysql"`
	RedisConf *RedisConfig `yaml:"redis"`
}

type LogConfig struct {
	logrus.TextFormatter
	Level      string `yaml:"level"`
	FilePath   string `yaml:"path"`
	TimeFormat string `yaml:"timeFormat"`
	ForceColor bool   `yaml:"forceColor"`
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
	//TODO: parse config file and return Config struct
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
	//fmt.Printf("MySQL %#v:", cfg.DBConf.MySQLConf)
	//fmt.Printf("Redis %#v:", cfg.DBConf.RedisConf)
	//fmt.Printf("Log %#v:", cfg.LogConf)
	return nil
}

func (format *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	tz, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("err in func Format:", err)
		return nil, err
	}
	entry.Time = entry.Time.In(tz)
	entry.Time.Format(format.logConf.TimeFormat)

	//自定义日志输出格式
	formatted := fmt.Sprintf("[%s] %s Error in [%s]: %s\n",
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

	file, err := os.OpenFile(logConf.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	//设置多路输出
	Output := io.MultiWriter(os.Stdout, file)
	logrusLog.Log.SetOutput(Output)

	//设置日志格式
	logrusLog.Log.SetFormatter(&LogFormat{logConf: logConf})
	fmt.Println("init Log Successfully.")
	return nil
}

func InitConfig(configPath string) error {
	err := parseConfig(configPath)
	if err != nil {
		return err
	}
	err = InitLog()
	if err != nil {
		return err
	}
	fmt.Println("init config success.")
	return nil
}
