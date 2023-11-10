package sqlFunc

import (
	"encoding/json"
	"github.com/dop251/goja"
)

type (
	ErrRet struct {
		Err error `json:"err"`
	}
)

func ToJsStu(vm *goja.Runtime, stu any) goja.Value {
	b, err := json.Marshal(stu)
	if err != nil {
		panic(err)
	}
	var ret = map[string]any{}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(err)
	}
	return vm.ToValue(ret)
}
