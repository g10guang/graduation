package handler

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/store"
	"github.com/g10guang/graduation/tools"
)

var storage store.Storage

func init() {
	tools.InitLog()
	//if tools.IsProductEnv() {
	//	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
	//} else {
	//	storage = store.NewLocalStorage()
	//}
	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
}
