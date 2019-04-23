package jobs

import (
	"github.com/g10guang/graduation/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const JobName_SaveFileMeta = "save_file_meta"

type SaveFileMetaJob struct {
	file *model.File
	db   *gorm.DB
}

func NewSaveFileMetaJob(file *model.File, conn *gorm.DB) *SaveFileMetaJob {
	j := &SaveFileMetaJob{
		file: file,
		db:   conn,
	}
	return j
}

func (j *SaveFileMetaJob) GetName() string {
	return JobName_SaveFileMeta
}

func (j *SaveFileMetaJob) Run() (interface{}, error) {
	defer func() {
		logrus.Debugf("job: %s exit", j.GetName())
	}()
	if err := j.db.Debug().Save(j.file).Error; err != nil {
		logrus.Errorf("Save FileMetas: %+v to mysql Error: %s", j.file, err)
		return nil, err
	}
	logrus.Infof("Save FileMetas: %+v success", j.file)
	return nil, nil
}
