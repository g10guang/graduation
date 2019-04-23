package handler

import (
	"context"
	"errors"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler interface {
	Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error
}

type CommonHandler struct {
	UserId int64
}

func NewCommonHandler() *CommonHandler {
	h := &CommonHandler{}
	return h
}

func (h *CommonHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *CommonHandler) CheckParams(ctx context.Context) error {
	if h.UserId == 0 {
		logrus.Errorf("uid == 0")
		return errors.New("uid == 0")
	}
	return nil
}

func (h *CommonHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if h.UserId, err = strconv.ParseInt(r.FormValue(constdef.Param_Uid), 10, 64); err != nil {
		logrus.Errorf("CommonHandler parseParams Error: %s", err)
	}
	return err
}

var storage store.Storage

func init() {
	//if tools.IsProductEnv() {
	//	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
	//} else {
	//	storage = store.NewLocalStorage()
	//}

	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
}
