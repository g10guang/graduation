package constdef

import (
	"fmt"
	"os"
)

// request 请求中的字段常量定义
const (
	Param_Uid      = "uid"
	Param_Fid      = "fid"
	Param_File     = "file"
	Param_Filename = "filename"
	Param_Offset   = "offset"
	Param_Limit    = "limit"
	Param_Format   = "format"
)

// OS environment variable key
const (
	ENV_TestEnv    = "test_env"
	ENV_ProductEnv = "product_env"
)

// NSQ 消息队列相关常量
const (
	// Picture Post Event Topic
	PostFileEventTopic = "post_file"
	// Delete Picture Event Topic
	DeleteFileEventTopic = "delete_file"
)

type ImageFormat int16

// 支持的图片格式
// 由于使用了 iota 所以新增格式只能够在最后追加
const (
	InvalidImageFormat ImageFormat = iota
	Jpeg
	Png
	WaterMarkJpeg
	WaterMarkPng
)

var ImageFormatList = []ImageFormat{Jpeg, Png, WaterMarkJpeg, WaterMarkPng}

// NSQ 配置信息
const (
	NsqLookupdAddr = "nsqlookupd:4161"
	NsqdAddr       = "nsqd:4150"
)

// HDFS 配置
const (
	WebHdfsAddr = "hadoop:50070"
	WebHdfsUser = "root"
	WebHdfsDir  = "/oss/image"
)

// MySQL 配置
var (
	MySqlDbName = "oss_meta"
	MySqlUrl    = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), "mysql", 3306, MySqlDbName)
)

// redis 配置
const (
	RedisAddr = "redis:6379"
)
