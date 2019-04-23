package jobs

import "github.com/g10guang/graduation/store"

const JobName_DeleteFile = "delete_file"

type DeleteFileJob struct {
	fid int64
	storage store.Storage
}

func NewDeleteFileJob(fid int64, storage store.Storage) *DeleteFileJob {
	j := &DeleteFileJob{
		fid: fid,
		storage: storage,
	}
	return j
}

func (j *DeleteFileJob) GetName() string {
	return JobName_DeleteFile
}

func (j *DeleteFileJob) Run() (interface{}, error) {
	var err error
	for i := 0; i < retryTime; i++ {
		if err = j.storage.Delete(j.fid); err == nil {
			break
		}
	}
	return nil, err
}
