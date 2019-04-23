package redis

var FileRedis *FileInfoRedis
var UserRedis *UserInfoRedis
var ContentRedis *FileContentRedis

func init() {
	FileRedis = NewFileInfoRedis()
	UserRedis = NewUserInfoRedis()
	ContentRedis = NewFileContentRedis()
}
