package handler

import (
	"context"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
)

type DeleteStorageHandler struct {
	msg *model.DeleteFileEvent
}

func NewDeleteStorageHandler(msg *model.DeleteFileEvent) *DeleteStorageHandler {
	h := &DeleteStorageHandler{
		msg: msg,
	}
	return h
}

func (h *DeleteStorageHandler) Handle(ctx context.Context) error {
	if err := redis.ContentRedis.Del([]int64{h.msg.Fid}); err != nil {
		logrus.Errorf("Delete Fid: %d cache Error: %s", h.msg.Fid, err)
	}
	if err := storage.Delete(h.msg.Fid); err != nil {
		logrus.Errorf("Delete Fid: %d Error: %s", h.msg.Fid, err)
	}
	// 用于缓存用户分页的缓存信息
	if err := redis.FileRedis.DelPageCache(h.msg.Uid); err != nil {
		logrus.Errorf("Delete Uid: %d File Page cache Error: %s", h.msg.Uid, err)
	}
	return nil
}