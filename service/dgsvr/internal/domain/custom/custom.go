package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/dop251/goja"
	"sync"
)

type (
	Info struct {
		ProductID       string
		TransformScript *string //可能为空
		ScriptLang      int64
		CustomTopic     []string
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
	ConvertTypeUp   ConvertType = iota //自定义协议转iThings协议
	ConvertTypeDown                    //iThings协议转自定义协议
)

func (i *Info) InitScript() *Vm {
	if i.TransformScript == nil || *i.TransformScript == "" {
		return nil
	}
	vmInfo := Vm{Pool: sync.Pool{New: func() any {
		vm := goja.New()
		_, err := vm.RunString(*i.TransformScript)
		if err != nil {
			return nil
		}
		return vm
	}}}
	return &vmInfo
}

// 上行数据
func (v *Vm) DataUp(ctx context.Context,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) ConvertFunc {
	vm := v.Get().(*goja.Runtime)
	funName := fmt.Sprintf("%s%sUp", handle, utils.FirstUpper(Type))
	convert, ok := goja.AssertFunction(vm.Get(funName))
	if !ok {
		return nil
	}
	return func(data []byte) ([]byte, error) {
		res, err := convert(goja.Undefined(), vm.ToValue(data))
		if err != nil {
			return nil, errors.System.AddMsgf("调用js函数失败:%v", err)
		}
		str, err := res.ToObject(nil).MarshalJSON()
		return str, err
	}
}

// 下行数据
func (v *Vm) DataDown(ctx context.Context,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) ConvertFunc {
	vm := v.Get().(*goja.Runtime)
	funName := fmt.Sprintf("%s%sDown", handle, utils.FirstUpper(Type))
	convert, ok := goja.AssertFunction(vm.Get(funName))
	if !ok {
		return nil
	}
	return func(data []byte) ([]byte, error) {
		var dataStr = map[string]any{}
		json.Unmarshal(data, &dataStr)
		res, err := convert(goja.Undefined(), vm.ToValue(dataStr))
		if err != nil {
			panic(err)
		}
		ret, err := res.ToObject(nil).MarshalJSON()
		if err != nil {
			return nil, err
		}
		var retData []byte
		err = json.Unmarshal(ret, &retData)
		if err != nil {
			return ret, nil
		}
		return retData, nil
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
