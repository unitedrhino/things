package script

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
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
	ConvertFunc func(data []byte) ([]byte, error)
)

func (i *Info) InitScript() Vm {
	vmInfo := Vm{Pool: sync.Pool{New: func() any {
		vm := goja.New()
		_, err := vm.RunString(i.Script)
		if err != nil {
			return nil
		}
		return vm
	}}}
	return vmInfo
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
