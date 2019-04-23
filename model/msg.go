package model

// 消息队列消息定义

// 上传图片 topic who when what
type PostFileEvent struct {
	Fid       int64
	Uid       int64
	Timestamp int64
	Extra     string
}

// 删除图片 topic who when what
type DeleteFileEvent struct {
	Fid       int64
	Uid       int64
	Timestamp int64
	Extra     string
}
