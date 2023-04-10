package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/devices"
	"strings"
	"sync"
)

type (
	Info struct {
		ProductID string
		Script    string
		Lang      int64
	}
	Vm struct {
		sync.Pool
	}
	ConvertFunc  func(data []byte) ([]byte, error)
	TransformRet struct {
		Topic   string `json:"topic"`
		PayLoad []byte `json:"payLoad"`
	}
	TransFormFunc func(topic string, data []byte) (*TransformRet, error)
	Repo          interface {
		// GetProtoFunc 自定义协议转iTHings协议
		GetProtoFunc(ctx context.Context, productID string, dir ConvertType,
			handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
			Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) (ConvertFunc, error)
		// GetTransFormFunc 自定义topic转换函数
		GetTransFormFunc(ctx context.Context, productID string, direction devices.Direction) (TransFormFunc, error)
		// ClearCache 清除缓存
		ClearCache(ctx context.Context, productID string) error
	}
)

const (
	CustomType = "custom" //自定义协议topic type定义
)

type ConvertType int

const (
	ConvertTypeRowToProto ConvertType = iota //自定义协议转iThings协议
	ConvertTypeProtoToRow                    //iThings协议转自定义协议
	ConvertTypeTransform                     //自定义topic协议转换
)

func (i *Info) InitScript() *Vm {
	vmInfo := Vm{Pool: sync.Pool{New: func() any {
		vm := goja.New()
		_, err := vm.RunString(i.Script)
		if err != nil {
			return nil
		}
		return vm
	}}}
	return &vmInfo
}

func (v *Vm) RawDataToProtocol(ctx context.Context,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) ConvertFunc {
	vm := v.Get().(*goja.Runtime)
	funName := fmt.Sprintf("%s%sProtocolToRawData", handle, strings.ToTitle(Type))
	convert, ok := goja.AssertFunction(vm.Get(funName))
	if !ok {
		return nil
	}
	return func(data []byte) ([]byte, error) {
		res, err := convert(goja.Undefined(), vm.ToValue(data))
		if err != nil {
			panic(err)
		}
		return res.ToObject(nil).MarshalJSON()
	}
}
func (v *Vm) ProtocolToRawData(ctx context.Context,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) ConvertFunc {
	vm := v.Get().(*goja.Runtime)
	funName := fmt.Sprintf("%s%sProtocolToRawData", handle, strings.ToTitle(Type))
	convert, ok := goja.AssertFunction(vm.Get(funName))
	if !ok {
		return nil
	}
	return func(data []byte) ([]byte, error) {
		res, err := convert(goja.Undefined(), vm.ToValue(data))
		if err != nil {
			panic(err)
		}
		return res.ToObject(nil).MarshalJSON()
	}
}

func (v *Vm) TransformPayload(ctx context.Context, dir devices.Direction) TransFormFunc {
	vm := v.Get().(*goja.Runtime)
	dirStr := "up"
	if dir == devices.Down {
		dirStr = "down"
	}
	funName := fmt.Sprintf("%sTransformPayload", dirStr)
	convert, ok := goja.AssertFunction(vm.Get(funName))
	if !ok {
		return nil
	}
	return func(topic string, data []byte) (*TransformRet, error) {
		res, err := convert(goja.Undefined(), vm.ToValue(topic), vm.ToValue(data))
		if err != nil {
			panic(err)
		}
		ret, err := res.ToObject(nil).MarshalJSON()
		if err != nil {
			return nil, err
		}
		var retStu TransformRet
		err = json.Unmarshal(ret, &retStu)
		if err != nil {
			return nil, err
		}
		return &retStu, nil
	}
}
