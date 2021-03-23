package errors

import (
	"encoding/json"
	"fmt"
)


type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c CodeError)WithMsg(msg string) *CodeError{
	return &CodeError{Code: c.Code, Msg: msg}
}
func (c CodeError)AddMsg(msg string) *CodeError{
	return &CodeError{Code: c.Code, Msg: c.Msg+":"+msg}
}

func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return SysErrorDefault.WithMsg(msg)
}

func (e *CodeError) Error() string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`,e.Code,e.Msg)
}

//将普通的error及转换成json的error或error类型的转回自己的error
func Fmt(errs error) *CodeError{
	switch errs.(type) {
	case *CodeError:
		return errs.(*CodeError)
	default:
		var ce CodeError
		err := json.Unmarshal([]byte(errs.Error()),&ce)
		if err != nil {
			return SysErrorDefault.AddMsg(errs.Error())
		}
		return &ce
	}
}