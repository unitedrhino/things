package common

import "errors"

// 不指定x-oss-forbid-overwrite时，默认覆盖同名Object。
// 指定x-oss-forbid-overwrite为false时，表示允许覆盖同名Object。
// 指定x-oss-forbid-overwrite为true时，表示禁止覆盖同名Object，如果同名Object已存在，程序将报错。
var (
	ForbidWriteErr = errors.New(" Object Forbid Over Write")
)
