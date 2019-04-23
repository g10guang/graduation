package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"image"
	"io"
	"strings"
	"sync"
)

type CompressHandler struct {
	msg *model.PostFileEvent
}

func NewCompressHandler(msg *model.PostFileEvent) *CompressHandler {
	h := &CompressHandler{
		msg: msg,
	}
	return h
}

// jpeg/png压缩 + 水印
func (h *CompressHandler) Handle(ctx context.Context) error {
	reader, err := storage.Read(h.msg.Fid)
	if err != nil {
		return err
	}
	meta, err := mysql.FileMySQL.Get(h.msg.Fid)
	if err != nil {
		return err
	}

	imageFormat := h.JudgeFormat(ctx, meta.Name)
	if imageFormat == constdef.InvalidImageFormat {
		logrus.Errorf("invalid image format. image name: %s", meta.Name)
		return nil
	}

	im, err := h.DecodeImage(ctx, reader, imageFormat)
	if err != nil {
		logrus.Errorf("DecodeImage Error: %s", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		h.AddWaterMark(ctx, im)
	}()

	go func() {
		defer wg.Done()
		h.ImageStore(im, constdef.Jpeg)
	}()

	go func() {
		defer wg.Done()
		h.ImageStore(im, constdef.Png)
	}()

	wg.Wait()
	return nil
}

func (h *CompressHandler) JudgeFormat(ctx context.Context, filename string) constdef.ImageFormat {
	imageFormat := constdef.InvalidImageFormat
	switch {
	case strings.HasSuffix(filename, "jpeg") || strings.HasSuffix(filename, "jpg"):
		imageFormat = constdef.Jpeg
	case strings.HasSuffix(filename, "png"):
		imageFormat = constdef.Png
	default:
		imageFormat = constdef.InvalidImageFormat
	}
	return imageFormat
}

func (h *CompressHandler) DecodeImage(ctx context.Context, r io.Reader, format constdef.ImageFormat) (image.Image, error) {
	switch format {
	case constdef.Jpeg:
		return tools.JpegDecode(r)
	case constdef.Png:
		return tools.PngDecode(r)
	default:
		panic(fmt.Errorf("unknown image format: %v", format))
	}
}

func (h *CompressHandler) AddWaterMark(ctx context.Context, im image.Image) error {
	if im == nil {
		panic(nil)
	}
	dim := tools.ConvertImage2Draw(im)
	if dim == nil {
		logrus.Errorf("cannot convert image.Image to draw.Image")
		return fmt.Errorf("cannot convert image.Image to draw.Image")
	}
	logrus.Debugf("Add Water Mark Fid: %d", h.msg.Fid)
	tools.WaterMark(dim, fmt.Sprintf("uid=%d", h.msg.Uid))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		h.ImageStore(im, constdef.WaterMarkPng)
	}()
	go func() {
		defer wg.Done()
		h.ImageStore(im, constdef.WaterMarkJpeg)
	}()
	wg.Wait()
	return nil
}

// 将 image 输出到 storage 存储
func (h *CompressHandler) ImageStore(im image.Image, format constdef.ImageFormat) error {
	buf := &bytes.Buffer{}
	writer := bufio.NewWriter(buf)
	var f func(image.Image, io.Writer) error
	switch format {
	case constdef.Jpeg, constdef.WaterMarkJpeg:
		f = tools.JpegCompress
	case constdef.Png, constdef.WaterMarkPng:
		f = tools.PngCompress
	default:
		logrus.Errorf("Unknown image format: %v", format)
		return fmt.Errorf("Unknown image format: %v", format)
	}
	if err := f(im, writer); err != nil {
		logrus.Errorf("Image format: %v Encode: %s", format, err)
		return err
	}
	if err := storage.Write(h.msg.Fid, buf, format); err != nil {
		logrus.Errorf("storage write format: %v Error: %s", format, err)
		return err
	}
	return nil
}
