package redis

import (
	"fmt"
	"time"

	"github.com/g10guang/graduation/constdef"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// 图片内容的缓存
type FileContentRedis struct {
	conn *redis.Client
}

func NewFileContentRedis() *FileContentRedis {
	r := &FileContentRedis{}
	r.conn = redis.NewClient(&redis.Options{
		Addr:     constdef.RedisAddr,
		Password: "",
		DB:       1,
	})
	return r
}

func (r *FileContentRedis) genKey(fid int64, format ...constdef.ImageFormat) string {
	if len(format) == 0 || format[0] == constdef.InvalidImageFormat {
		return fmt.Sprintf("c_%d", fid)
	}
	return fmt.Sprintf("c_%d_%d", fid, format[0])
}

func (r *FileContentRedis) genKeys(fids []int64, format ...constdef.ImageFormat) []string {
	keys := make([]string, len(fids))
	for i, fid := range fids {
		keys[i] = r.genKey(fid, format...)
	}
	return keys
}

func (r *FileContentRedis) Set(fid int64, data []byte, format ...constdef.ImageFormat) error {
	if _, err := r.conn.Set(r.genKey(fid, format...), data, time.Minute).Result(); err != nil {
		logrus.Errorf("redis Set fid %d content Error: %s", fid, err)
		return err
	}
	return nil
}

func (r *FileContentRedis) Get(fid int64, format ...constdef.ImageFormat) ([]byte, error) {
	b, err := r.conn.Get(r.genKey(fid, format...)).Bytes()
	if err != nil {
		logrus.Errorf("redis Get fid %d content Error: %s", fid, err)
		return nil, err
	}
	return b, nil
}

func (r *FileContentRedis) Del(fids []int64, format ...constdef.ImageFormat) error {
	if err := r.conn.Del(r.genKeys(fids, format...)...).Err(); err != nil {
		logrus.Errorf("redis Del fids: %v content Error: %s", fids, err)
		return err
	}
	return nil
}
