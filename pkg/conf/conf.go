package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	App          AppConf
	Auth         AuthConf
	Log          LogConf
	Openobserver OpenobserverConf
	Db           DbConf
	Cache        CacheConf
	File         FileConf
	Devops       DevopsConf
)

type Conf struct {
	AppConf          AppConf          `yaml:"app"`
	AuthConf         AuthConf         `yaml:"auth"`
	DbConf           DbConf           `yaml:"db"`
	CacheConf        CacheConf        `yaml:"cache"`
	FileConf         FileConf         `yaml:"file"`
	LogConf          LogConf          `yaml:"log"`
	OpenobserverConf OpenobserverConf `yaml:"openobserver"`
	DevopsConf       DevopsConf       `yaml:"devops"`
}

// DevopsConf DevOps 流程执行等工作区配置
type DevopsConf struct {
	// WorkspaceRoot 本地工作区根（clone/产物/日志）；为空则运行接口拒绝执行
	WorkspaceRoot string `yaml:"workspaceRoot"`
}

type AppConf struct {
	Id         uint16 `yaml:"id"`
	Name       string `yaml:"name"`
	Doc        string `yaml:"doc"`
	Version    string `yaml:"version"`
	BaseUrl    string `yaml:"baseUrl"`
	Addr       string `yaml:"addr"`
	Port       int    `yaml:"port"`
	Env        string `yaml:"env"`
	EncryptKey string `yaml:"encryptKey"` // 非空时启用 API 加解密中间件（与前端 VITE_API_ENCRYPT_KEY 一致）
}

type AuthConf struct {
	RetryCount int64    `yaml:"retryCount"`
	RetryTime  int64    `yaml:"retryTime"`
	Use        string   `yaml:"use"`
	TokenName  string   `yaml:"tokenName"`
	Timeout    int64    `yaml:"timeout"`
	KeyPrefix  string   `yaml:"keyPrefix"`
	PubUrl     []string `yaml:"pubUrl"`
}

// LogConf 与配置中 log 段扁平字段一致（见 admin/cmd/setting.yaml）
type LogConf struct {
	Level      string `yaml:"level"`
	StdOut     bool   `yaml:"stdOut"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
	LocalTime  bool   `yaml:"LocalTime"`
}

type OpenobserverConf struct {
	BaseUrl   string `yaml:"baseUrl"`
	Streams   string `yaml:"streams"`
	Index     string `yaml:"index"`
	User      string `yaml:"user"`
	Token     string `yaml:"token"`
	Interval  int    `yaml:"interval"`
	Threshold int    `yaml:"threshold"`
}

type DbConf struct {
	MaxIdleConn     int    `yaml:"maxIdleConn"`
	MaxOpenConn     int    `yaml:"maxOpenConn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"` // 连接最大存活时间(秒)，默认60
	ShowSQL         bool   `yaml:"showSQL"`
	LogLevel        int    `yaml:"logLevel"`
	Backup         struct {
		Enable       bool   `yaml:"enable"`
		Cron         string `yaml:"cron"`         // 如 "0 0 0 * * *" 每天0点
		Path         string `yaml:"path"`        // 备份目录
		KeepDays     int    `yaml:"keepDays"`    // 保留天数
		MaxSize      int    `yaml:"maxSize"`     // 最大总大小(MB)
		MaxCount     int    `yaml:"maxCount"`    // 最大备份文件数
		Compress     bool   `yaml:"compress"`   // 是否 gzip 压缩
		Encrypt      bool   `yaml:"encrypt"`    // 是否加密
		EncryptKey   string `yaml:"encryptKey"`
		EncryptMethod string `yaml:"encryptMethod"` // aes
		EncryptLevel string `yaml:"encryptLevel"`   // low/medium/high
	} `yaml:"backup"`
	UseDb string `yaml:"useDb"`
	Sqlite struct {
		Db string `yaml:"db"`
	} `yaml:"sqlite"`
	Mysql struct {
		Db   string `yaml:"db"`
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pwd  string `yaml:"pwd"`
	} `yaml:"mysql"`
	Pgsql struct {
		Db   string `yaml:"db"`
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pwd  string `yaml:"pwd"`
	} `yaml:"pgsql"`
}

type CacheConf struct {
	Local struct {
		Expire int `yaml:"expire"`
		Purges int `yaml:"purges"`
	} `yaml:"local"`
	Redis struct {
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
		Pwd  string `yaml:"pwd"`
		Db   int    `yaml:"db"`
	} `yaml:"redis"`
}

type FileConf struct {
	BaseUrl string `yaml:"baseUrl"`
	Local   string `yaml:"local"`
	Oss     struct {
		Endpoint   string `yaml:"endpoint"`
		AccessKey  string `yaml:"accessKey"`
		SecretKey  string `yaml:"secretKey"`
		BucketName string `yaml:"bucketName"`
	} `yaml:"oss"`
}

func New() {
	// 解析配置文件
	fmt.Println("Optional startup parameters, such as:\n - go run main.go -c setting.yaml\ndefault use setting.yaml")
	confPath := flag.String("c", "setting.yaml", "配置文件路径")
	flag.Parse()
	fmt.Println("The configuration file is " + *confPath)
	stRootDir, err := os.Getwd()
	if err != nil {
		panic("config init error: failed to get working directory: " + err.Error())
	}
	fPath := filepath.Join(stRootDir, *confPath)
	file, err := os.ReadFile(fPath)
	if err != nil {
		panic("config init error:" + err.Error())
	}

	var conf Conf
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		panic("config init error:" + err.Error())
	}
	App = conf.AppConf
	Auth = conf.AuthConf
	Db = conf.DbConf
	Cache = conf.CacheConf
	File = conf.FileConf
	Log = conf.LogConf
	Openobserver = conf.OpenobserverConf
	Devops = conf.DevopsConf
	log.Println("✅ config ok!")
}
