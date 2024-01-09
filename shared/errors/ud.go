package errors

const UdError = 4000000

var (
	SceneTriggerType = NewCodeError(UdError+1, "触发类型不支持")
)
