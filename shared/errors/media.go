package errors

const MEDIA_ERROR = 3000000

var (
	MediaCreateError   = NewCodeError(MEDIA_ERROR+1, "流服务创建失败")
	MediaUpdateError   = NewCodeError(MEDIA_ERROR+2, "流服务更新失败")
	MediaNotfoundError = NewCodeError(MEDIA_ERROR+3, "流服务不存在")
	MediaActiveError   = NewCodeError(MEDIA_ERROR+4, "流服务激活失败")
)
