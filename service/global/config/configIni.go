package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type Info struct {
	SqlConfig     *SqlConfigStruct
	RConfig       *RConfigStruct
	ProjectConfig *ProjectConfigStruct
	LiveConfig    *LiveConfigStruct
	AliyunOss     *AliyunOss
	ProjectUrl    string
}

func init() {
	//Avoiding global duplication of package guides
	ReturnsInstance()
}

var Config = new(Info)
var cfg *ini.File
var err error

type SqlConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Host     int    `ini:"host"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

type RConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
}
type LiveConfigStruct struct {
	IP        string `ini:"ip"`
	Agreement string `ini:"agreement"`
	RTMP      string `ini:"rtmp"`
	FLV       string `ini:"flv"`
	HLS       string `ini:"hls"`
	Api       string `ini:"api"`
}

type ProjectConfigStruct struct {
	ProjectStates bool   `ini:"project_states"`
	Url           string `ini:"url"`
	UrlTest       string `ini:"url_test"`
}

type AliyunOss struct {
	Region                   string `ini:"region"`
	Bucket                   string `ini:"bucket"`
	AccessKeyId              string `ini:"accessKeyId"`
	AccessKeySecret          string `ini:"accessKeySecret"`
	Host                     string `ini:"host"`
	CallbackUrl              string `ini:"callbackUrl"`
	Endpoint                 string `ini:"endpoint"`
	RoleArn                  string `ini:"roleArn"`
	RoleSessionName          string `ini:"roleSessionName"`
	DurationSeconds          int    `ini:"durationSeconds"`
	IsOpenTranscoding bool `ini:"isOpenTranscoding"`
	TranscodingTemplate360p  string `ini:"transcodingTemplate360p"`
	TranscodingTemplate480p  string `ini:"transcodingTemplate480p"`
	TranscodingTemplate720p  string `ini:"transcodingTemplate720p"`
	TranscodingTemplate1080p string `ini:"transcodingTemplate1080p"`
}

func ReturnsInstance() *Info {
	Config.SqlConfig = &SqlConfigStruct{}
	cfg, err = ini.Load("config/config.ini")
	if err != nil {
		fmt.Printf("Configuration file does not exist, please check the environment. %v \n", err)
		os.Exit(1)
	}

	err = cfg.Section("mysql").MapTo(Config.SqlConfig)
	if err != nil {
		fmt.Printf("Mysql Read Configuration File Error. %v \n", err)
		os.Exit(1)
	}
	Config.RConfig = &RConfigStruct{}
	err = cfg.Section("redis").MapTo(Config.RConfig)
	if err != nil {
		fmt.Printf("Redis Read Configuration File Error. %v \n", err)
		os.Exit(1)
	}
	Config.ProjectConfig = &ProjectConfigStruct{}
	err = cfg.Section("project").MapTo(Config.ProjectConfig)
	if err != nil {
		fmt.Printf("Project read configuration file error. %v \n", err)
		os.Exit(1)
	}

	Config.LiveConfig = &LiveConfigStruct{}
	err = cfg.Section("live").MapTo(Config.LiveConfig)
	if err != nil {
		fmt.Printf("Live read configuration file error. %v \n", err)
		os.Exit(1)
	}

	Config.AliyunOss = &AliyunOss{}
	err = cfg.Section("aliyunOss").MapTo(Config.AliyunOss)
	if err != nil {
		fmt.Printf("Live read configuration file error. %v \n", err)
		os.Exit(1)
	}

	//Determining whether it is a formal environment
	if Config.ProjectConfig.ProjectStates {
		Config.ProjectUrl = Config.ProjectConfig.Url
	} else {
		Config.ProjectUrl = Config.ProjectConfig.UrlTest
	}

	return Config
}
