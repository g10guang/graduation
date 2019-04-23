package handler

import (
	"context"
	"encoding/json"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// 按照 uid 纬度获取文件信息，支持分页功能
type UserFileHandler struct {
	uid       int64
	offset    int64
	limit     int64
	fileMetas []*model.File
}

func NewUserFileHandler() *UserFileHandler {
	h := new(UserFileHandler)
	return h
}

func (h *UserFileHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	if err := h.parseParams(ctx, r); err != nil {
		h.genResponse(out, http.StatusBadRequest)
		return err
	}

	jobmgr := tools.NewJobMgr(time.Second)
	jobmgr.AddJob(loader.NewUserFileLoader(h.uid, h.offset, h.limit))
	if err := jobmgr.Start(ctx); err != nil {
		logrus.Errorf("jobmgr Error: %s", err)
		h.genResponse(out, http.StatusInternalServerError)
		return err
	}
	if result := jobmgr.GetResult(loader.LoaderName_UserFile); result.Result != nil {
		if v, ok := result.Result.([]*model.File); ok {
			h.fileMetas = v
		}
	}
	h.genResponse(out, http.StatusOK)
	return nil
}

func (h *UserFileHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	h.offset, err = strconv.ParseInt(r.FormValue(constdef.Param_Offset), 10, 64)
	if h.limit, err = strconv.ParseInt(r.FormValue(constdef.Param_Limit), 10, 64); err != nil {
		h.limit = 10
	}
	if h.uid, err = strconv.ParseInt(r.FormValue(constdef.Param_Uid), 10, 64); err != nil {
		logrus.Errorf("parse uid Error: %s", err)
		return err
	}
	return nil
}

func (h *UserFileHandler) genResponse(out http.ResponseWriter, statusCode int) {
	if statusCode == http.StatusOK {
		b, err := json.Marshal(h.fileMetas)
		if err != nil {
			logrus.Errorf("json Marshal Error: %s", err)
			out.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = out.Write(b)
		if err != nil {
			logrus.Errorf("write http body Error: %s", err)
			out.WriteHeader(http.StatusInternalServerError)
			return
		}
		out.WriteHeader(http.StatusOK)
	} else {
		out.WriteHeader(statusCode)
	}
}
