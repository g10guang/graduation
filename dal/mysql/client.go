package mysql

var FileMySQL *FileInfoMySql
var UserMySQL *UserInfoMySql

func init() {
	FileMySQL = NewFileInfoMySql()
	UserMySQL = NewUserInfoMySql()
}
