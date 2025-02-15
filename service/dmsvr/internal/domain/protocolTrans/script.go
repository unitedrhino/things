package protocolTrans

import (
	"context"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"github.com/traefik/yaegi/interp"
	"reflect"
	"sort"
	"sync"
)

type ScriptInfo struct {
	Priority   int64
	ScriptLang int64
	Script     string
}

type ScriptInfos []ScriptInfo

// 实现 sort.Interface 接口的 Len 方法
func (a ScriptInfos) Len() int {
	return len(a)
}

// 实现 sort.Interface 接口的 Less 方法
func (a ScriptInfos) Less(i, j int) bool {
	return a[i].Priority < a[j].Priority
}

// 实现 sort.Interface 接口的 Swap 方法
func (a ScriptInfos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ScriptTrans struct {
	ScriptSymbol           map[string]map[string]reflect.Value
	ProductUpBeforeCache   map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceUpBeforeCache    map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductUpAfterCache    map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceUpAfterCache     map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductDownBeforeCache map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceDownBeforeCache  map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductDownAfterCache  map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceDownAfterCache   map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductUpBeforeMutex   sync.RWMutex
	DeviceUpBeforeMutex    sync.RWMutex
	ProductUpAfterMutex    sync.RWMutex
	DeviceUpAfterMutex     sync.RWMutex
	ProductDownBeforeMutex sync.RWMutex
	DeviceDownBeforeMutex  sync.RWMutex
	ProductDownAfterMutex  sync.RWMutex
	DeviceDownAfterMutex   sync.RWMutex
}

func NewScriptTrans() *ScriptTrans {
	s := ScriptTrans{ScriptSymbol: make(map[string]map[string]reflect.Value),
		ProductUpBeforeCache:   make(map[string]map[devices.MsgHandle]map[string]ScriptInfos),
		DeviceUpBeforeCache:    make(map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos),
		ProductUpAfterCache:    make(map[string]map[devices.MsgHandle]map[string]ScriptInfos),
		DeviceUpAfterCache:     make(map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos),
		ProductDownBeforeCache: make(map[string]map[devices.MsgHandle]map[string]ScriptInfos),
		DeviceDownBeforeCache:  make(map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos),
		ProductDownAfterCache:  make(map[string]map[devices.MsgHandle]map[string]ScriptInfos),
		DeviceDownAfterCache:   make(map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos),
	}
	return &s
}

func (s *ScriptTrans) AddSymbol(key string, syb map[string]reflect.Value) {
	s.ScriptSymbol[key] = syb
}

func (s *ScriptTrans) DownBeforeTrans(ctx context.Context, msg *deviceMsg.PublishMsg) *deviceMsg.PublishMsg {
	return msg
}

func (s *ScriptTrans) UpBeforeTrans(ctx context.Context, msg *deviceMsg.PublishMsg) *deviceMsg.PublishMsg {
	//todo 后面需要加上缓存
	var out = *msg
	var scripts ScriptInfos
	func() {
		s.ProductUpBeforeMutex.RLock()
		defer s.ProductUpBeforeMutex.RUnlock()
		pc, ok := s.ProductUpBeforeCache[msg.ProductID]
		if ok {
			script := s.GetScripts(ctx, pc, msg)
			scripts = append(scripts, script...)
		}
	}()
	func() {
		s.DeviceUpBeforeMutex.RLock()
		defer s.DeviceUpBeforeMutex.RUnlock()
		dc, ok := s.DeviceUpBeforeCache[devices.Core{ProductID: msg.ProductID, DeviceName: msg.DeviceName}]
		if ok {
			script := s.GetScripts(ctx, dc, msg)
			scripts = append(scripts, script...)
		}
	}()

	sort.Sort(scripts)
	for _, script := range scripts {
		i := interp.New(interp.Options{})
		i.Use(s.ScriptSymbol)
		_, err := i.EvalWithContext(ctx, script.Script)
		if err != nil {
			continue
		}
		handle, err := i.EvalWithContext(ctx, `Handle`)
		if err != nil {
			continue
		}
		fn, ok := handle.Interface().(func(context.Context, *deviceMsg.PublishMsg) *deviceMsg.PublishMsg)
		if !ok {
			continue
		}
		newMsg := fn(ctx, &out)
		if newMsg != nil {
			out = *newMsg
		}
	}
	return &out
}
func (s *ScriptTrans) GetScripts(ctx context.Context, script map[devices.MsgHandle]map[string]ScriptInfos, msg *deviceMsg.PublishMsg) (ret ScriptInfos) {
	hs := func(h map[string]ScriptInfos) {
		if h == nil {
			return
		}
		{
			t, ok := h[msg.Type]
			if ok {
				ret = append(ret, t...)
			}
		}
		{
			t, ok := h[All]
			if ok {
				ret = append(ret, t...)
			}
		}
	}
	hs(script[msg.Handle])
	hs(script[All])
	return
}
