package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
)

const LoaderName_FileMeta = "file_meta_loader"

type FileMetaLoader struct {
	fids []int64
}

func NewFileMetaLoader(fids []int64) *FileMetaLoader {
	l := &FileMetaLoader{
		fids: fids,
	}
	return l
}

func (l *FileMetaLoader) GetName() string {
	return LoaderName_FileMeta
}

// 1、尝试从 redis 缓存中获取
// 2、如果缓存没有命中，访问 mysql
// 3、异步设置 redis 缓存
func (l *FileMetaLoader) Run() (interface{}, error) {
	metas, missFids, err := redis.FileRedis.MGet(l.fids)
	if err != nil {
		// redis 出错尝试 mysql
		// 因为 redis 缓存获取非核心逻辑，所以允许失败，但是在高并发场景下，可能会将下游 MySQL 打扒，导致整个体统瘫痪
		logrus.Errorf("redis Error: %s", err)
	}

	if len(missFids) == 0 {
		// 全部 cache 命中
		logrus.Debugf("FileInfoMeta redis cache hit Fid: %v", l.fids)
		return metas, nil
	}

	if metas == nil {
		metas = make(map[int64]*model.File, len(l.fids))
	}

	metasFromMySQL, err := mysql.FileMySQL.MGet(missFids)
	if err != nil {
		return nil, err
	}

	// 只缓存 missFids
	go l.saveRedisCache(metasFromMySQL)
	for _, m := range metasFromMySQL {
		metas[m.Fid] = m
	}

	return metas, nil
}

func (l *FileMetaLoader) saveRedisCache(metas []*model.File) {
	if len(metas) == 0 {
		return
	}
	redis.FileRedis.MSet(metas)
}
