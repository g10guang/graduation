package handler

import (
	"context"
	"encoding/json"
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/write_api/jobs"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	*CommonHandler
	File       multipart.File
	FileHeader *multipart.FileHeader
	FileMeta   *model.File
}

func NewPostHandler() *PostHandler {
	h := &PostHandler{
		CommonHandler: NewCommonHandler(),
	}
	return h
}

func (h *PostHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	var err error
	if err = h.parseParams(ctx, r); err != nil {
		// 用户参数错误
		h.genResponse(out, http.StatusBadRequest)
		return err
	}
	h.BuildFileMeta()
	if err = h.SaveFile(ctx); err != nil {
		h.genResponse(out, http.StatusInternalServerError)
		return err
	}
	out.Header().Set("fid", strconv.FormatInt(h.FileMeta.Fid, 10))
	h.genResponse(out, http.StatusOK)
	return nil
}

func (h *PostHandler) parseParams(ctx context.Context, r *http.Request) error {
	var err error
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	h.File, h.FileHeader, err = r.FormFile(constdef.Param_File)
	if err != nil {
		logrus.Errorf("Get Upload FileBytes Error: %s", err.Error())
		return err
	}
	return nil
}

// 有可能有网络错误，唯一键冲突等
// 因为整个过程涉及到不少网络操作，所以需要使用事务，免得数据库中插入了无用记录
func (h *PostHandler) SaveFile(ctx context.Context) (err error) {
	logrus.Debugf("SaveFile fid: %d", h.FileMeta.Fid)
	db := mysql.FileMySQL.Begin()
	defer func() {
		if err != nil {
			logrus.Debugf("post mysql rollback")
			db.Rollback()
			go storage.Delete(h.FileMeta.Fid)
		} else {
			logrus.Debugf("post mysql commit")
			db.Commit()
			// 异步发送消息队列
			go h.PublishPostFileEvent()
		}
	}()

	jobmgr := tools.NewJobMgr(time.Second * 3)
	jobmgr.AddJob(jobs.NewSaveFileMetaJob(h.FileMeta, db))
	jobmgr.AddJob(jobs.NewStoreFileJob(h.FileMeta.Fid, h.File, storage))
	if err = jobmgr.Start(ctx); err != nil {
		logrus.Errorf("batch Job process Error: %s", err)
		return err
	}

	return nil
}

// 发送一条消息到消息队列
func (h *PostHandler) PublishPostFileEvent() error {
	msg := &model.PostFileEvent{
		Fid:       h.FileMeta.Fid,
		Uid:       h.FileMeta.Uid,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(msg)
	return mq.PublishNsq(constdef.PostFileEventTopic, b)
}

func (h *PostHandler) BuildFileMeta() {
	now := time.Now()
	h.FileMeta = &model.File{
		Uid:        h.UserId,
		Fid:        tools.GenID().Int64(),
		Name:       h.FileHeader.Filename,
		Size:       h.FileHeader.Size,
		//Md5:        "", 		// Md5 计算放在 consumer 中执行
		CreateTime: now,
		UpdateTime: now,
	}
}

// 生成响应
func (h *PostHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
}
