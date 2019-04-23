package handler

import (
	"context"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type CommonHandler struct {
	Fid        int64
	FileMeta   model.File
	FileReader io.Reader
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	h.Fid, err = strconv.ParseInt(r.FormValue(constdef.Param_Fid), 10, 64)
	if err != nil {
		logrus.Errorf("parse Fid Error: %s", err)
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
