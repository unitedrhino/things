package errors

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yl/shared/proto"
)


type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Details []string `json:"details,omitempty"`
}

type rpcError interface {
	GRPCStatus() *status.Status
}

//func TogRPCError(err *Error) error {
//	s, _ := status.New(ToRPCCode(err.Code()), err.Msg()).WithDetails(&pb.Error{Code: int32(err.Code()), Message: err.Msg()})
//	return s.Err()
//}

func (c CodeError)ToRpc() error{
	s, _ := status.New(codes.Unknown, c.Msg).
		WithDetails(&proto.Error{Code: int32(c.Code), Message: c.Msg,Detail: c.Details})
	return s.Err()
}

func (c CodeError)WithMsg(msg string) *CodeError{
	return &CodeError{Code: c.Code, Msg: msg}
}
func (c CodeError)AddDetail(msg string) *CodeError{
	 c.Details = append(c.Details, msg)
	 return &c
}

func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
	return Default.WithMsg(msg)
}

func (e CodeError) Error() string {

	ret,_ :=json.Marshal(e)
	return string(ret)
}

//将普通的error及转换成json的error或error类型的转回自己的error
func Fmt(errs error) *CodeError{
	switch errs.(type) {
	case *CodeError:
		return errs.(*CodeError)
	case rpcError://如果是grpc类型的错误
		s, _ := status.FromError(errs)
		if len(s.Details()) == 0 {
			err := fmt.Sprintf("rpc err detail is nil|err=%#v",s)
			return System.AddDetail(err)
		}
		if er,ok:= s.Details()[0].(*proto.Error);ok{
			return &CodeError{Code: int(er.Code),Msg: er.Message,Details: er.Detail}
		}
		err := fmt.Sprintf("rpc err not suppot|err=%#v",s)
		return System.AddDetail(err)
	default:
		var ce CodeError
		err := json.Unmarshal([]byte(errs.Error()),&ce)
		if err != nil {
			return System.AddDetail(errs.Error())
		}
		return &ce
	}
}


//func FromError(err error) *status.Status {
//	s, _ := status.FromError(err)
//	return s
//}
