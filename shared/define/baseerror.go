package define

import "fmt"

const defaultCode = 1001

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c CodeError)WithMsg(msg string) error{
	return &CodeError{Code: c.Code, Msg: msg}
}

func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

func (e *CodeError) Error() string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`,e.Code,e.Msg)
}