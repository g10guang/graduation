package main

import (
	"net/http"

	"github.com/g10guang/graduation/read_api/handler"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
)

func main() {
	tools.InitLog()
	var err error
	initHttpHandler()
	if err = http.ListenAndServe("0.0.0.0:10002", nil); err != nil {
		logrus.Panicf("http.ListenAndServe Error: %s", err)
	}
	logrus.Infof("Main goroutine exit")
}

func initHttpHandler() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/meta", meta)
	http.HandleFunc("/userFile", userFile)
}

func get(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewGetHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("get Error: %s", err)
	}
}

func meta(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewMetaHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("head Error: %s", err)
	}
}

func userFile(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewUserFileHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("userFile Error: %s", err)
	}
}
