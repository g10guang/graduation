package redis

import (
	"fmt"
	"github.com/g10guang/graduation/model"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestFileSet(t *testing.T) {
	logrus.SetReportCaller(true)
	err := FileRedis.Set(&model.File{
		Fid:        1,
		Uid:        1,
		Name:       "test",
		Size:       10,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Md5:        "5d41402abc4b2a76b9719d911017c592",
		Extra:      "hello",
		Status:     0,
	})
	if err != nil {
		logrus.Errorf("Redis Set Error: %s", err)
	} else {
		logrus.Infof("Redis Set Success")
	}
}

func TestFileGet(t *testing.T) {
	file, err := FileRedis.Get(1)
	if err != nil {
		logrus.Errorf("Redis Get Error: %s", err)
	} else {
		logrus.Infof("file: %+v Error: %s", file, err)
	}
}

func TestUserSet(t *testing.T) {
	err := UserRedis.Set(&model.User{
		Uid:    1,
		Status: 1,
		Extra:  "hello world",
	})
	if err != nil {
		logrus.Errorf("Redis Set Error: %s", err)
	} else {
		logrus.Infof("Redis Set Success")
	}
}

func TestUserGet(t *testing.T) {
	user, err := UserRedis.Get(1)
	if err != nil {
		logrus.Errorf("Redis Get Error: %s", err)
	} else {
		logrus.Infof("Redis Get User: %+v", user)
	}
}

func TestPipeline(t *testing.T) {
	pipe := FileRedis.conn.Pipeline()
	pipe.Get("hello")
	pipe.Get("world")
	pipe.SetNX("hello", "100", 0)
	cmd, err := pipe.Exec()
	if err != nil {
		panic(err)
	}
	for _, c := range cmd {
		logrus.Infof("cmd: %+v", c)
		s, ok := c.(*redis.StringCmd)
		if !ok {
			panic(errors.New("type conversion error"))
		}
		r, err := s.Result()
		if err != nil {
			panic(err)
		}
		logrus.Infof("r: %+v", r)
	}
}

func TestMGet(t *testing.T) {
	r, err := FileRedis.conn.MGet("hello", "world", "no_exist").Result()
	if err != nil {
		logrus.Panicf("access redis Error: %s", err)
		panic(err)
	}
	for _, v := range r {
		logrus.Info(v == nil)
		if v != nil {
			logrus.Info(v.(string))
		}
	}
}

func TestBytes(t *testing.T) {
	r, err := ContentRedis.conn.Set("hello", []byte{10, 0, 0, 10, 0}, time.Second*100).Result()
	if err != nil {
		panic(err)
	}
	logrus.Infof("redis Content: %v", r)
	p, err := ContentRedis.conn.Get("hello").Bytes()
	if err != nil {
		panic(err)
	}
	logrus.Infof("redis response: %v", p)
}

func TestBytes2(t *testing.T) {
	FileRedis.conn.Set("hello", []byte{0, 10, 0, 0, 1, 0}, time.Hour)
	result, err := FileRedis.conn.MGet("hello", "world").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n%+v", result[0], result[1])
	fmt.Printf("%T\t%T", result[0], result[1])
}
