package errors

const MEDIA_ERROR = 3000000

var (
	MediaCreateError = NewCodeError(MEDIA_ERROR+1, "流服务创建失败")
	//MediaRepeatCreateError = NewCodeError(MEDIA_ERROR+1, "流服务端口和IP重复,创建失败")
	MediaUpdateError = NewCodeError(MEDIA_ERROR+2, "流服务更新失败")
)
