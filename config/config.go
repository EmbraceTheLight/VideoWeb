package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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
type Config struct {
	DBConf *DBConfig `yaml:"database"`
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

func InitConfig() {
	err := parseConfig("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("init config success:%v", cfg.DBConf.MySQLConf)
}

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
	fmt.Printf("MySQL %#v:", cfg.DBConf.MySQLConf)
	fmt.Printf("Redis %#v:", cfg.DBConf.RedisConf)
	return nil
}
