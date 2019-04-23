package jobs

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/sirupsen/logrus"
)

const JobName_DeleteFileMeta = "delete_file_meta"
const retryTime = 2

type DeleteFileMetaJob struct {
	fids []int64
	uid  int64
}

func NewDeleteFileMetaJob(fids []int64, uid int64) *DeleteFileMetaJob {
	j := &DeleteFileMetaJob{
		fids: fids,
		uid:  uid,
	}
	return j
}

func (j *DeleteFileMetaJob) GetName() string {
	return JobName_DeleteFileMeta
}

// 采用事务
func (j *DeleteFileMetaJob) Run() (interface{}, error) {
	var err error
	conn := mysql.FileMySQL.Begin()
	defer func() {
		if err != nil {
			// 执行失败回滚 mysql
			conn.Rollback()
		} else {
			conn.Commit()
			// 延时删除
			go func() {
				redis.FileRedis.Del(j.fids)
			}()
		}
	}()
	if err = mysql.FileMySQL.Delete(conn, j.fids, j.uid); err != nil {
		logrus.Errorf("Delete MySQL Error: %s", err)
		return nil, err
	}

	if err = redis.FileRedis.Del(j.fids); err != nil {
		return nil, err
	}
	return nil, nil
}
