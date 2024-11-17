package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
	"sync/atomic"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 导入批量更新设备
func NewMultiUpdateImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateImportLogic {
	return &MultiUpdateImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type updateImport struct {
	ProductID   string
	DeviceName  string
	DeviceAlias string
	SchemaAlias map[string]string
	Gateway     devices.Core
	Group       map[string]*AddGroup
}
type AddGroup struct {
	Purpose   string //用途
	GroupName map[string]struct{}
}

func (u updateImport) Validate() error {
	if u.ProductID == "" || u.DeviceName == "" {
		return errors.Parameter.AddMsg("产品ID和设备ID必填")
	}
	if u.Gateway.ProductID == "" && u.Gateway.DeviceName != "" {
		return errors.Parameter.AddMsg("填写了绑定网关的设备ID后也需要填写绑定网关的产品ID")
	}
	return nil
}

func getHandle(tt string) func(r *updateImport, value string) error {
	switch tt {
	case "产品ID":
		return func(r *updateImport, value string) error {
			if value == "" {
				return errors.Parameter.AddMsg("产品ID必填")
			}
			r.ProductID = value
			return nil
		}
	case "设备ID":
		return func(r *updateImport, value string) error {
			if value == "" {
				return errors.Parameter.AddMsg("设备ID必填")
			}
			r.DeviceName = value
			return nil
		}
	case "设备名称":
		return func(r *updateImport, value string) error {
			r.DeviceAlias = value
			return nil
		}
	case "绑定网关产品ID":
		return func(r *updateImport, value string) error {
			r.Gateway.ProductID = value
			return nil
		}
	case "绑定网关设备ID":
		return func(r *updateImport, value string) error {
			r.Gateway.DeviceName = value
			return nil
		}
	case "物模型昵称", "schemaAlias":
		return func(r *updateImport, value string) error {
			if r.SchemaAlias == nil {
				r.SchemaAlias = make(map[string]string)
			}
			kvs := strings.Split(value, ",")
			for _, kv := range kvs {
				k, v, ok := strings.Cut(kv, ":")
				if !ok {
					return errors.Parameter.AddMsgf("格式不对:%v,格式为: switch:213,switch2:xfwef", value)
				}
				r.SchemaAlias[k] = v
			}
			return nil
		}
	default: //自定义字段
		g, h, ok := strings.Cut(tt, ".")
		if !ok {
			return nil
		}
		switch g {
		case "deviceGroup":
			return func(r *updateImport, value string) error {
				if r.Group == nil {
					r.Group = map[string]*AddGroup{}
				}
				if r.Group[h] == nil {
					r.Group[h] = &AddGroup{Purpose: h, GroupName: make(map[string]struct{})}
				}
				gns := strings.Split(value, ",")
				for _, gn := range gns {
					r.Group[h].GroupName[gn] = struct{}{}
				}
				return nil
			}
		case "device":
			return getHandle(h)
		}
	}
	return nil
}

func (l *MultiUpdateImportLogic) MultiUpdateImport(rows [][]string) (*types.DeviceMultiUpdateImportResp, error) {
	ds, err := l.svcCtx.DictM.DictDetailIndex(l.ctx, &sys.DictDetailIndexReq{
		DictCode: "deviceMultiUpdateImportRow",
		Status:   def.True,
	})
	if err != nil {
		return nil, err
	}
	var dictMap = map[string]string{}
	for _, v := range ds.List {
		dictMap[v.Label] = v.Value
	}
	if len(rows) < 2 {
		return nil, errors.Parameter.AddMsg("请至少输入一行")
	}
	var titles = rows[0]
	var handle = make([]func(r *updateImport, value string) error, len(titles))
	for i, t := range titles {
		tt, _, _ := strings.Cut(t, "(")
		tt, _, _ = strings.Cut(tt, "（")
		tt = strings.TrimSpace(tt)
		f := getHandle(tt)
		if f == nil {
			d := dictMap[tt]
			if d == "" {
				return nil, errors.Parameter.AddMsgf("不支持的列名:%v", t)
			}
			f = getHandle(d)
			if f == nil {
				return nil, errors.Parameter.AddMsgf("不支持的列名:%v", t)
			}
		}
		handle[i] = f

	}
	var cols []updateImport
	var needAddGroup = map[string]map[string]struct{}{}
	for _, v := range rows[1:] {
		var col updateImport
		for i, value := range v {
			err := handle[i](&col, value)
			if err != nil {
				return nil, err
			}
		}
		err = col.Validate()
		if err != nil {
			return nil, err
		}
		if col.Group != nil {
			for _, v := range col.Group {
				if needAddGroup[v.Purpose] == nil {
					needAddGroup[v.Purpose] = make(map[string]struct{})
				}
				for n := range v.GroupName {
					needAddGroup[v.Purpose][n] = struct{}{}
				}

			}
		}
		cols = append(cols, col)
	}
	var egg errgroup.Group
	egg.SetLimit(100)
	var succ atomic.Int64
	//var groupAdd = map[int64]map[devices.Core]struct{}{}
	var errNum atomic.Int64
	//if len(needAddGroup)>0{
	//	var gs []*dm.GroupInfo
	//	for purpose,v:=range needAddGroup {
	//		for name:=range v{
	//			gs=append(gs, &dm.GroupInfo{
	//				Purpose:     purpose,
	//				Name:        name,
	//			})
	//		}
	//	}
	//	_,err:=l.svcCtx.DeviceG.GroupInfoMultiCreate(l.ctx,&dm.GroupInfoMultiCreateReq{Groups: gs})
	//	if err!=nil{
	//		return nil, err
	//	}
	//	for purpose,v:=range needAddGroup {
	//		gis,err:=l.svcCtx.DeviceG.GroupInfoIndex(l.ctx, &dm.GroupInfoIndexReq{
	//			Names:     utils.SetToSlice(v),
	//			Purpose:  purpose,
	//		})
	//		if err!=nil{
	//			return nil, err
	//		}
	//		for _,gi:=range gis.List{
	//			if groupAdd[gi.Id]==nil{
	//				groupAdd[gi.Id] = make(map[devices.Core]struct{})
	//			}
	//			groupAdd[gi.Id][devices.Core{
	//				ProductID:  col.,
	//				DeviceName: "",
	//			}]
	//		}
	//	}
	//
	//}
	var errDetail = []types.DeviceMultiUpdateImportError{}
	var mutex sync.Mutex
	var addErr = func(col updateImport, err error) {
		errNum.Add(1)
		mutex.Lock()
		defer mutex.Unlock()
		er := errors.Fmt(err)
		errDetail = append(errDetail, types.DeviceMultiUpdateImportError{
			Device: types.DeviceCore{ProductID: col.ProductID, DeviceName: col.DeviceName},
			Code:   er.Code,
			Msg:    er.GetMsg(),
		})
	}
	for _, c := range cols {
		col := c
		egg.Go(func() error {
			defer utils.Recover(l.ctx)
			di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
				ProductID:  col.ProductID,
				DeviceName: col.DeviceName,
			})
			if err != nil {
				l.Errorf("col:%v,err:%v", col, err)
				addErr(col, err)
				return nil
			}
			var updateOne = &dm.DeviceInfo{ProductID: di.ProductID, DeviceName: di.DeviceName}
			if col.SchemaAlias != nil {
				if di.SchemaAlias == nil {
					updateOne.SchemaAlias = col.SchemaAlias
				} else {
					updateOne.SchemaAlias = di.SchemaAlias
					for k, v := range col.SchemaAlias {
						updateOne.SchemaAlias[k] = v
					}
				}
			}
			updateOne.DeviceAlias = utils.ToRpcNullString(col.DeviceAlias)
			_, err = l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, updateOne)
			if err != nil {
				l.Errorf("col:%v,err:%v", col, err)
				addErr(col, err)
				return nil
			}
			if col.Gateway.DeviceName != "" && di.DeviceType == def.DeviceTypeSubset {
				_, err := l.svcCtx.DeviceM.DeviceGatewayMultiCreate(l.ctx, &dm.DeviceGatewayMultiCreateReq{
					Gateway: &dm.DeviceCore{
						ProductID:  col.Gateway.ProductID,
						DeviceName: col.Gateway.DeviceName,
					},
					IsNotNotify: false,
					IsAuthSign:  false,
					List:        []*dm.DeviceGatewayBindDevice{{ProductID: col.ProductID, DeviceName: col.DeviceName}},
				})
				if err != nil {
					l.Errorf("col:%v,err:%v", col, err)
					addErr(col, err)
					return nil
				}
			}
			succ.Add(1)
			return nil
			//if len(col.Group) > 0 {
			//}
		})
	}
	err = egg.Wait()
	return &types.DeviceMultiUpdateImportResp{
		Total:       int64(len(rows) - 1),
		ErrCount:    errNum.Load(),
		IgnoreCount: 0,
		ErrDetail:   errDetail,
		SuccCount:   succ.Load(),
	}, err
}
