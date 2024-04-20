package config

import (
	"VideoWeb/logrusLog"
	"bufio"
	"fmt"
	"github.com/shiena/ansicolor"
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

	//var color string
	//switch entry.Level {
	//case logrus.InfoLevel:
	//	color = define.Blue
	//case logrus.WarnLevel:
	//	color = define.Yellow
	//case logrus.ErrorLevel:
	//	color = define.Red
	//default:
	//	color = define.Reset
	//}

	//自定义日志输出格式
	formatted := fmt.Sprintf("[%s] %s %s in [%s]: %s\n",
		//color,
		entry.Level.String(),
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Data["type"],
		entry.Data["function"],
		entry.Message,
		//define.Reset,
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
	writer1 := ansicolor.NewAnsiColorWriter(os.Stdout)
	writer2 := ansicolor.NewAnsiColorWriter(file)
	Output := io.MultiWriter(writer1, writer2)
	logrusLog.Log.SetOutput(Output)

	logrusLog.Log.SetFormatter(&LogFormat{logConf: logConf})
	fmt.Println("init Log Format,logConf is:", logConf)
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
