package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
	"io"
)

type ChecksumHandler struct {
	msg *model.PostFileEvent
}

func NewChecksumHandler(msg *model.PostFileEvent) *ChecksumHandler {
	h := &ChecksumHandler{
		msg: msg,
	}
	return h
}

func (h *ChecksumHandler) Handle(ctx context.Context) error {
	reader, err := storage.Read(h.msg.Fid)
	if err != nil {
		return err
	}
	checksum := md5.New()
	_, err = io.Copy(checksum, reader)
	if err != nil {
		logrus.Errorf("Write md5 checksum Error: %s", err)
	}

	err = mysql.FileMySQL.UpdateMd5(h.msg.Fid, hex.EncodeToString(checksum.Sum(nil)))
	if err != nil {
		logrus.Errorf("Update Fid: %d md5 checksum Error: %s", h.msg.Fid, err)
	}

	// 删除用户 uid 分页文件信息缓存
	if err = redis.FileRedis.DelPageCache(h.msg.Uid); err != nil {
		logrus.Errorf("delete uid: %d file page cache Error: %s", h.msg.Uid, err)
	}

	return nil
}
