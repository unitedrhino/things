package errors

const FileError = 1000000

var (
	Upload = NewCodeError(FileError+1, "上传失败")
)
