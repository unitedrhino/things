package errors

const MediaError = 3000000

var (
	MediaCreateError       = NewCodeError(MediaError+1, "流服务创建失败")
	MediaUpdateError       = NewCodeError(MediaError+2, "流服务更新失败")
	MediaNotfoundError     = NewCodeError(MediaError+3, "流服务不存在")
	MediaActiveError       = NewCodeError(MediaError+4, "流服务激活失败")
	MediaPullCreateError   = NewCodeError(MediaError+5, "拉流创建失败")
	MediaStreamDeleteError = NewCodeError(MediaError+6, "流删除错误")
	MediaRecordNotFound    = NewCodeError(MediaError+7, "未找到录像列表")
	MediaSipUpdateError    = NewCodeError(MediaError+8, "ID或channelID不能都为空")
	MediaSipDevCreateError = NewCodeError(MediaError+9, "设备创建失败")
	MediaSipChnCreateError = NewCodeError(MediaError+10, "通道创建失败")
	MediaSipPlayError      = NewCodeError(MediaError+11, "通道播放失败")
)
