package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/g10guang/graduation/write_api/jobs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DeleteHandler struct {
	*CommonHandler
	Fids []int64
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, http.StatusBadRequest)
		return err
	}
	if err = h.delete_(ctx); err != nil {
		h.genResponse(out, http.StatusInternalServerError)
		return err
	}
	h.genResponse(out, http.StatusOK)
	return nil
}

func (h *DeleteHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	fids := r.PostForm[constdef.Param_Fid]
	h.Fids = make([]int64, len(fids))
	for i, fid := range fids {
		h.Fids[i], err = strconv.ParseInt(fid, 10, 64)
		if err != nil {
			logrus.Errorf("parse fid Error: %s", err)
		}
	}
	if len(h.Fids) == 0 {
		logrus.Errorf("empty delete fids")
		return errors.New("empty delete fids")
	}
	logrus.Infof("fids: %+v", h.Fids)
	return
}

func (h *DeleteHandler) delete_(ctx context.Context) (err error) {
	jobmgr := tools.NewJobMgr(time.Second * 2)
	jobmgr.AddJob(jobs.NewDeleteFileMetaJob(h.Fids, h.UserId))
	if err := jobmgr.Start(ctx); err != nil {
		logrus.Errorf("delete fids job exec Error: %s", err)
		return err
	}
	// 发送消息队列异步化
	go h.PublishDeleteEvent(ctx)
	return
}

func (h *DeleteHandler) PublishDeleteEvent(ctx context.Context) (err error) {
	for _, fid := range h.Fids {
		msg := &model.DeleteFileEvent{
			Uid:       h.UserId,
			Fid:       fid,
			Timestamp: time.Now().Unix(),
		}
		b, _ := json.Marshal(msg)
		logrus.Debugf("topic: %s publish msg: %s", constdef.DeleteFileEventTopic, string(b))
		if err = mq.PublishNsq(constdef.DeleteFileEventTopic, b); err != nil {
			logrus.Errorf("Publish Delete Event Error: %s", err)
		}
	}
	return err
}

func (h *DeleteHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
}
