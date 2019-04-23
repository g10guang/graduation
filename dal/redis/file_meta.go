package redis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/g10guang/graduation/model"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/g10guang/graduation/constdef"
)

type FileInfoRedis struct {
	conn *redis.Client
}

func NewFileInfoRedis() *FileInfoRedis {
	r := &FileInfoRedis{}
	r.conn = redis.NewClient(&redis.Options{
		Addr:     constdef.RedisAddr,
		Password: "",
		DB:       0,
	})
	return r
}

func (r *FileInfoRedis) genKey(fid int64) string {
	return fmt.Sprintf("f%d", fid)
}

func (r *FileInfoRedis) genKeys(fids []int64) []string {
	keys := make([]string, len(fids))
	for i, fid := range fids {
		keys[i] = r.genKey(fid)
	}
	return keys
}

// 只要获取 redis 失败就认为 cache not found
// 不管是网络超时还是确实 key 不存在
// TODO 区分清楚是 redis key 不存在的情况，免得压垮下游
func (r *FileInfoRedis) Get(fid int64) (meta model.File, err error) {
	s, err := r.conn.Get(r.genKey(fid)).Result()
	if err != nil {
		logrus.Errorf("redis Get Fid: %d not found", fid)
		return meta, err
	}
	if err = json.NewDecoder(strings.NewReader(s)).Decode(&meta); err != nil {
		logrus.Errorf("unmarshal redis cache Error: %s", err)
	}
	return
}

func (r *FileInfoRedis) MGet(fids []int64) (metas map[int64]*model.File, missFids []int64, err error) {
	fidsKey := r.genKeys(fids)
	result, err := r.conn.MGet(fidsKey...).Result()
	if err != nil {
		// redis error
		logrus.Errorf("redis mget Error: %s", err)
		// 如果 redis 发生错误，则所有 fids 都 miss
		return nil, fids, err
	}
	metas = make(map[int64]*model.File, len(fids))
	for i, v := range result {
		if v == nil {
			// cache not found
			missFids = append(missFids, fids[i])
			logrus.Debugf("fid: %d redis cache not found", fids[i])
			continue
		}
		m := new(model.File)
		if err = json.NewDecoder(strings.NewReader(v.(string))).Decode(m); err != nil {
			missFids = append(missFids, fids[i])
			logrus.Errorf("unmarshal FileMetas Error: %s", err)
			continue
		}
		// cache hit
		logrus.Debugf("fid: %d redis cache hit", fids[i])
		metas[fids[i]] = m
	}
	return metas, missFids, nil
}

func (r *FileInfoRedis) Del(fids []int64) error {
	fidsKey := r.genKeys(fids)
	if _, err := r.conn.Del(fidsKey...).Result(); err != nil {
		logrus.Errorf("delete redis cache Error: %s", err)
		return err
	}
	return nil
}

// 因为 MSet 不支持超时
func (r *FileInfoRedis) Set(file *model.File) error {
	b, err := json.Marshal(file)
	logrus.Debugf("set redis cache fid: %d json: %s", file.Fid, string(b))
	if err != nil {
		logrus.Errorf("json Marshal Error: %s", err)
		return err
	}
	if _, err = r.conn.Set(r.genKey(file.Fid), string(b), time.Minute*5).Result(); err != nil {
		logrus.Errorf("redis Set file: %+v", file)
	}
	return err
}

func (r *FileInfoRedis) MSet(files []*model.File) error {
	pipe := r.conn.Pipeline()
	for _, meta := range files {
		b, err := json.Marshal(meta)
		if err != nil {
			logrus.Errorf("json Marshal Error: %s", err)
			continue
		}
		pipe.Set(r.genKey(meta.Fid), string(b), time.Minute*5)
	}
	_, err := pipe.Exec()
	if err != nil {
		logrus.Errorf("Multi Set redis pipeline Error: %s", err)
		return err
	}
	return nil
}

// 注意 Get Set 都是设置的 []byte 而不是 string
// 因为接口返回的是 interface{} 类型，需要注意类型
func (r *FileInfoRedis) GetPageCache(uid, offset, limit int64) (metas []*model.File, err error) {
	key, field := r.genPageKeyField(uid, offset, limit)
	b, err := r.conn.HGet(key, field).Bytes()
	if err != nil {
		logrus.Errorf("redis HGet Error: %s", err)
		return nil, err
	}

	if err := json.Unmarshal(b, &metas); err != nil {
		logrus.Errorf("json Unmarshal Error: %s", err)
		return nil, err
	}

	return metas, nil
}

func (r *FileInfoRedis) SetPageCache(uid, offset, limit int64, metas []*model.File) error {
	key, field := r.genPageKeyField(uid, offset, limit)
	b, err := json.Marshal(&metas)
	if err != nil {
		logrus.Errorf("json Marshal Error: %s", err)
		return err
	}
	if err := r.conn.HSet(key, field, b).Err(); err != nil {
		logrus.Errorf("redis HSet Error: %s", err)
		return err
	}
	if err = r.conn.Expire(key, time.Minute).Err(); err != nil {
		logrus.Errorf("redis expire page cache Error: %s", err)
		return err
	}
	return nil
}

func (r *FileInfoRedis) DelPageCache(uid int64) error {
	key, _ := r.genPageKeyField(uid, 0, 0)
	if err := r.conn.Del(key).Err(); err != nil {
		logrus.Errorf("redis Del uid: %d page cache Error: %s", uid, err)
		return err
	}
	return nil
}

func (r *FileInfoRedis) genPageKeyField(uid, offset, limit int64) (key, field string) {
	key = fmt.Sprintf("p_%d", uid)
	field = fmt.Sprintf("%d_%d", offset, limit)
	return
}
