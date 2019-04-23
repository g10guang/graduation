package handler

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/store"
)

var storage store.Storage

func init() {
	//if tools.IsProductEnv() {
	//	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
	//} else {
	//	storage = store.NewLocalStorage()
	//}

	storage = store.NewHdfsStorage(constdef.WebHdfsAddr, constdef.WebHdfsUser, constdef.WebHdfsDir)
}
