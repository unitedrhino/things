package errors

const FILE_ERROR = 1000000

var (
	Upload = NewCodeError(FILE_ERROR+1, "上传失败")
)
