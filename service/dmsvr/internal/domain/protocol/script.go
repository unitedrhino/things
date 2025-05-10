package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"sort"
	"sync"
	"time"
)

type ScriptInfo struct {
	Name       string
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
	loadFunc               func(context.Context, *ScriptTrans) error
	ScriptSymbol           map[string]map[string]reflect.Value
	ProductUpBeforeCache   map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceUpBeforeCache    map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductUpAfterCache    map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceUpAfterCache     map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
	ProductDownBeforeCache map[string]map[devices.MsgHandle]map[string]ScriptInfos       //第一级是
	DeviceDownBeforeCache  map[devices.Core]map[devices.MsgHandle]map[string]ScriptInfos //第一级是
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
	}
	s.AddSymbol("gjson/gjson", map[string]reflect.Value{
		"Get":          reflect.ValueOf(gjson.Get),
		"GetMany":      reflect.ValueOf(gjson.GetMany),
		"GetManyBytes": reflect.ValueOf(gjson.GetManyBytes),
		"GetBytes":     reflect.ValueOf(gjson.GetBytes),
		"Parse":        reflect.ValueOf(gjson.Parse),
		"ParseBytes":   reflect.ValueOf(gjson.ParseBytes),
		"Set":          reflect.ValueOf(sjson.Set),
		"SetBytes":     reflect.ValueOf(sjson.SetBytes),
		"Delete":       reflect.ValueOf(sjson.Delete),
		"DeleteBytes":  reflect.ValueOf(sjson.DeleteBytes),
		"SetRaw":       reflect.ValueOf(sjson.SetRaw),
		"SetRawBytes":  reflect.ValueOf(sjson.SetRawBytes),
	})
	s.AddSymbol("json/json", map[string]reflect.Value{
		"Marshal":   reflect.ValueOf(json.Marshal),
		"Unmarshal": reflect.ValueOf(json.Unmarshal),
	})
	s.AddSymbol("utils/utils", map[string]reflect.Value{
		"ToInt64":                reflect.ValueOf(utils.ToInt64),
		"ToBool":                 reflect.ValueOf(utils.ToBool),
		"ToTime":                 reflect.ValueOf(utils.ToTime),
		"ToDuration":             reflect.ValueOf(utils.ToDuration),
		"ToFloat64":              reflect.ValueOf(utils.ToFloat64),
		"ToFloat32":              reflect.ValueOf(utils.ToFloat32),
		"ToInt32":                reflect.ValueOf(utils.ToInt32),
		"ToInt16":                reflect.ValueOf(utils.ToInt16),
		"ToInt8":                 reflect.ValueOf(utils.ToInt8),
		"ToInt":                  reflect.ValueOf(utils.ToInt),
		"ToUint":                 reflect.ValueOf(utils.ToUint),
		"ToUint64":               reflect.ValueOf(utils.ToUint64),
		"ToUint32":               reflect.ValueOf(utils.ToUint32),
		"ToUint16":               reflect.ValueOf(utils.ToUint16),
		"ToUint8":                reflect.ValueOf(utils.ToUint8),
		"ToString":               reflect.ValueOf(utils.ToString),
		"BoolToInt":              reflect.ValueOf(utils.BoolToInt),
		"ToStringMapStringSlice": reflect.ValueOf(utils.ToStringMapStringSlice),
		"ToStringMapBool":        reflect.ValueOf(utils.ToStringMapBool),
		"ToStringMapInt":         reflect.ValueOf(utils.ToStringMapInt),
		"ToStringMapInt64":       reflect.ValueOf(utils.ToStringMapInt64),
		"ToSlice":                reflect.ValueOf(utils.ToSlice),
		"ToBoolSlice":            reflect.ValueOf(utils.ToBoolSlice),
		"ToStringSlice":          reflect.ValueOf(utils.ToStringSlice),
		"ToIntSlice":             reflect.ValueOf(utils.ToIntSlice),
		"ToDurationSlice":        reflect.ValueOf(utils.ToDurationSlice),
	})
	ctx := ctxs.WithRoot(context.Background())
	utils.Go(ctx, func() {
		var t = time.NewTicker(10 * time.Minute) //10分钟刷新一次
		for range t.C {
			if s.loadFunc != nil {
				if err := s.loadFunc(ctx, &s); err != nil {
					logx.WithContext(ctx).Error(err.Error())
				}
			}
		}
	})
	return &s
}

func (s *ScriptTrans) SetLoad(f func(context.Context, *ScriptTrans) error) {
	s.loadFunc = f
	ctx := ctxs.WithRoot(context.Background())
	if err := s.loadFunc(ctx, s); err != nil {
		logx.WithContext(ctx).Error(err.Error())
	}
}

func (s *ScriptTrans) AddSymbol(key string, syb map[string]reflect.Value) {
	if s.ScriptSymbol[key] == nil {
		s.ScriptSymbol[key] = syb
		return
	}
	for k, v := range syb {
		s.ScriptSymbol[key][k] = v
	}
}

func (s *ScriptTrans) GetFunc(ctx context.Context, script string, funcName string) (any, *[]string, error) {
	var logs = make([]string, 0, 5)
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(s.ScriptSymbol)
	i.Use(interp.Symbols)
	i.Use(map[string]map[string]reflect.Value{
		"log/log": {
			"PrintLn": reflect.ValueOf(func(v ...any) {
				logs = append(logs, fmt.Sprint(v...))
			}),
			"Print": reflect.ValueOf(func(v ...any) {
				logs = append(logs, fmt.Sprint(v...))
			}),
			"Printf": reflect.ValueOf(func(format string, v ...any) {
				logs = append(logs, fmt.Sprintf(format, v...))
			}),
		},
	})
	_, err := i.EvalWithContext(ctx, script)
	if err != nil {
		return nil, nil, err
	}
	handle, err := i.EvalWithContext(ctx, `Handle`)
	if err != nil {
		return nil, nil, err
	}
	return handle.Interface(), &logs, nil
}
func (s *ScriptTrans) PublishMsgRun(ctx context.Context, msg *deviceMsg.PublishMsg, script string) (msgs *deviceMsg.PublishMsg, retLogs []string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.Parameter.AddMsgf("执行panic:%v", p)
			return
		}
	}()
	var out = *msg
	handle, logs, err := s.GetFunc(ctx, script, "Handle")
	if err != nil {
		return nil, nil, errors.Parameter.AddMsgf("脚本定义错误:%s", err.Error())
	}
	fn, ok := handle.(func(context.Context, *deviceMsg.PublishMsg) *deviceMsg.PublishMsg)
	if !ok {
		return nil, nil, errors.Parameter.AddMsg("结构体中需要定义: func Handle(context.Context, *dm.PublishMsg) *dm.PublishMsg")
	}
	newMsg := fn(ctx, &out)
	if newMsg != nil {
		out = *newMsg
	}
	return newMsg, *logs, nil
}

func (s *ScriptTrans) RespMsgRun(ctx context.Context, req *deviceMsg.PublishMsg, resp *deviceMsg.PublishMsg, script string) (retLogs []string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.Parameter.AddMsgf("执行panic:%v", p)
			return
		}
	}()
	handle, logs, err := s.GetFunc(ctx, script, "Handle")
	if err != nil {
		return nil, errors.Parameter.AddMsg("结构体中需要定义: func Handle(context.Context, *dm.PublishMsg) *dm.PublishMsg")
	}
	fn, ok := handle.(func(context.Context, *deviceMsg.PublishMsg, *deviceMsg.PublishMsg))
	if !ok {
		return nil, errors.Parameter.AddMsg("结构体中需要定义: func Handle(context.Context, *dm.PublishMsg) *dm.PublishMsg")
	}
	fn(ctx, req, resp)
	return *logs, nil
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

func (s *ScriptTrans) UpAfterTrans(ctx context.Context, req *deviceMsg.PublishMsg, resp *deviceMsg.PublishMsg) error {
	//todo 后面需要加上缓存
	var scripts ScriptInfos
	func() {
		s.ProductUpBeforeMutex.RLock()
		defer s.ProductUpBeforeMutex.RUnlock()
		pc, ok := s.ProductUpBeforeCache[req.ProductID]
		if ok {
			script := s.GetScripts(ctx, pc, req)
			scripts = append(scripts, script...)
		}
	}()
	func() {
		s.DeviceUpBeforeMutex.RLock()
		defer s.DeviceUpBeforeMutex.RUnlock()
		dc, ok := s.DeviceUpBeforeCache[devices.Core{ProductID: req.ProductID, DeviceName: req.DeviceName}]
		if ok {
			script := s.GetScripts(ctx, dc, req)
			scripts = append(scripts, script...)
		}
	}()
	if len(scripts) == 0 {
		return nil
	}
	sort.Sort(scripts)
	logs := make([]string, 0)
	for _, script := range scripts {
		log, err := s.RespMsgRun(ctx, req, resp, script.Script)
		if err != nil {
			continue
		}
		for _, l := range log {
			logs = append(logs, fmt.Sprintf("%s:[%s]  ", script.Name, l))
		}
	}
	if len(logs) > 0 {
		logx.WithContext(ctx).Infof("脚本执行日志为:%s", logs)
	}
	return nil
}

func (s *ScriptTrans) UpBeforeTrans(ctx context.Context, msg *deviceMsg.PublishMsg) *deviceMsg.PublishMsg {
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
	if len(scripts) == 0 {
		return &out
	}
	sort.Sort(scripts)
	logs := make([]string, 0)
	for _, script := range scripts {
		newMsg, log, err := s.PublishMsgRun(ctx, &out, script.Script)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			continue
		}
		for _, l := range log {
			logs = append(logs, fmt.Sprintf("%s:[%s]  ", script.Name, l))
		}
		if newMsg != nil {
			out = *newMsg
		}
	}
	if len(logs) > 0 {
		logx.WithContext(ctx).Infof("脚本执行日志为:%s", logs)
	}
	return &out
}

func (s *ScriptTrans) DownBeforeTrans(ctx context.Context, msg *deviceMsg.PublishMsg) *deviceMsg.PublishMsg {
	var out = *msg
	var scripts ScriptInfos
	func() {
		s.ProductDownBeforeMutex.RLock()
		defer s.ProductDownBeforeMutex.RUnlock()
		pc, ok := s.ProductDownBeforeCache[msg.ProductID]
		if ok {
			script := s.GetScripts(ctx, pc, msg)
			scripts = append(scripts, script...)
		}
	}()
	func() {
		s.DeviceDownBeforeMutex.RLock()
		defer s.DeviceDownBeforeMutex.RUnlock()
		dc, ok := s.DeviceDownBeforeCache[devices.Core{ProductID: msg.ProductID, DeviceName: msg.DeviceName}]
		if ok {
			script := s.GetScripts(ctx, dc, msg)
			scripts = append(scripts, script...)
		}
	}()
	if len(scripts) == 0 {
		return &out
	}
	sort.Sort(scripts)
	logs := make([]string, 0)
	for _, script := range scripts {
		newMsg, log, err := s.PublishMsgRun(ctx, &out, script.Script)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			continue
		}
		for _, l := range log {
			logs = append(logs, fmt.Sprintf("%s:[%s]  ", script.Name, l))
		}
		if newMsg != nil {
			out = *newMsg
		}
	}
	if len(logs) > 0 {
		logx.WithContext(ctx).Infof("脚本执行日志为:%s", logs)
	}
	return &out
}
