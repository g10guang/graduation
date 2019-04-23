package handler

import (
	"context"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GetHandler struct {
	*CommonHandler
	format constdef.ImageFormat
}

func NewGetHandler() *GetHandler {
	return &GetHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *GetHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, http.StatusBadRequest)
		return err
	}

	// 1、获取图片元信息
	// 2、获取图片二进制内容
	jobmgr := tools.NewJobMgr(time.Second * 3)
	jobmgr.AddJob(loader.NewFileMetaLoader([]int64{h.Fid}))
	jobmgr.AddJob(loader.NewFileContentLoader(h.Fid, storage, h.format))
	if err = jobmgr.Start(ctx); err != nil {
		h.genResponse(out, http.StatusInternalServerError)
		return err
	}
	if result := jobmgr.GetResult(loader.LoaderName_FileMeta); result.Result != nil {
		if v, ok := result.Result.(map[int64]*model.File); ok {
			h.FileMeta = *v[h.Fid]
		}
	}

	if result := jobmgr.GetResult(loader.LoaderName_FileContent); result.Result != nil {
		if v, ok := result.Result.(io.Reader); ok {
			h.FileReader = v
		}
	}

	// 返回正常
	h.genResponse(out, http.StatusOK)
	return nil
}

func (h *GetHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	switch r.FormValue(constdef.Param_Format) {
	case "jpeg":
		h.format = constdef.Jpeg
	case "jpeg_water":
		h.format = constdef.WaterMarkJpeg
	case "png":
		h.format = constdef.Png
	case "png_water":
		h.format = constdef.WaterMarkPng
	}
	return
}

func (h *GetHandler) genResponse(out http.ResponseWriter, statusCode int) {
	if statusCode == http.StatusOK {
		out.Header().Set("fid", strconv.FormatInt(h.FileMeta.Fid, 10))
		out.Header().Set("uid", strconv.FormatInt(h.FileMeta.Uid, 10))
		out.Header().Set("name", h.FileMeta.Name)
		out.Header().Set("size", strconv.FormatInt(h.FileMeta.Size, 10))
		out.Header().Set("md5", h.FileMeta.Md5)
		out.Header().Set("create_time", strconv.FormatInt(h.FileMeta.CreateTime.Unix(), 10))
		out.Header().Set("update_time", strconv.FormatInt(h.FileMeta.UpdateTime.Unix(), 10))
		_, err := io.Copy(out, h.FileReader)
		if err != nil {
			logrus.Errorf("write file to http response Error: %s", err)
			out.WriteHeader(http.StatusInternalServerError)
		} else {
			out.WriteHeader(http.StatusOK)
		}
	} else {
		out.WriteHeader(statusCode)
	}
}
