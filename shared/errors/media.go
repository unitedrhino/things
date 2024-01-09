package errors

const MediaError = 3000000

var (
	MediaCreateError   = NewCodeError(MediaError+1, "流服务创建失败")
	MediaUpdateError   = NewCodeError(MediaError+2, "流服务更新失败")
	MediaNotfoundError = NewCodeError(MediaError+3, "流服务不存在")
	MediaActiveError   = NewCodeError(MediaError+4, "流服务激活失败")
)
