package deviceinteractlogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"golang.org/x/sync/errgroup"
	"sync"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyGetReportMultiSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyGetReportMultiSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyGetReportMultiSendLogic {
	return &PropertyGetReportMultiSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 请求设备获取设备最新属性
func (l *PropertyGetReportMultiSendLogic) PropertyGetReportMultiSend(in *dm.PropertyGetReportMultiSendReq) (*dm.PropertyGetReportMultiSendResp, error) {
	var list []*dm.PropertyGetReportSendMsg
	var err error
	if len(in.DeviceNames) != 0 {
		list, err = l.MultiSendOneProductProperty(in)
		if err != nil {
			return nil, err
		}
	} else {
		list, err = l.MultiSendMultiProductProperty(in)
		if err != nil {
			return nil, err
		}
	}
	return &dm.PropertyGetReportMultiSendResp{List: list}, nil
}

func (l *PropertyGetReportMultiSendLogic) MultiSendOneProductProperty(in *dm.PropertyGetReportMultiSendReq) ([]*dm.PropertyGetReportSendMsg, error) {
	list := make([]*dm.PropertyGetReportSendMsg, 0)
	sigSend := NewPropertyGetReportSendLogic(l.ctx, l.svcCtx)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, dev := range in.DeviceNames {
		v := dev
		wg.Add(1)
		utils.Go(l.ctx, func() {
			defer wg.Done()
			ret, err := sigSend.PropertyGetReportSend(&dm.PropertyGetReportSendReq{
				ProductID:  in.ProductID,
				DataIDs:    in.DataIDs,
				DeviceName: v,
			})
			if err != nil {
				myErr, _ := err.(*errors.CodeError)
				msg := &dm.PropertyGetReportSendMsg{
					ProductID:  in.ProductID,
					DeviceName: v,
					SysMsg:     myErr.GetMsg(),
					SysCode:    myErr.Code,
				}
				if ret != nil {
					msg.Code = ret.Code
					msg.Msg = ret.Msg
					msg.MsgToken = ret.MsgToken
				}
				mu.Lock()
				defer mu.Unlock()
				list = append(list, msg)
				return
			}
			msg := &dm.PropertyGetReportSendMsg{
				ProductID:  in.ProductID,
				DeviceName: v,
				SysCode:    errors.OK.Code,
				SysMsg:     errors.OK.GetMsg(),
				Code:       ret.Code,
				Msg:        ret.Msg,
				MsgToken:   ret.MsgToken,
				Timestamp:  ret.Timestamp,
				Params:     ret.Params,
			}
			mu.Lock()
			defer mu.Unlock()
			list = append(list, msg)
		})
	}
	wg.Wait()
	return list, nil
}

func (l *PropertyGetReportMultiSendLogic) MultiSendMultiProductProperty(in *dm.PropertyGetReportMultiSendReq) ([]*dm.PropertyGetReportSendMsg, error) {
	var productMap = map[string]map[string]struct{}{} //key是产品id,value是产品下的设备列表
	for _, v := range in.Devices {
		if productMap[v.ProductID] == nil {
			productMap[v.ProductID] = map[string]struct{}{v.DeviceName: {}}
		} else {
			productMap[v.ProductID][v.DeviceName] = struct{}{}
		}
	}
	if in.AreaID != 0 {
		dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{AreaIDs: []int64{in.AreaID}}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range dis {
			if productMap[v.ProductID] == nil {
				productMap[v.ProductID] = map[string]struct{}{v.DeviceName: {}}
			} else {
				productMap[v.ProductID][v.DeviceName] = struct{}{}
			}
		}
	}
	if in.AreaIDPath != "" {
		dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{AreaIDPath: in.AreaIDPath}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range dis {
			if productMap[v.ProductID] == nil {
				productMap[v.ProductID] = map[string]struct{}{v.DeviceName: {}}
			} else {
				productMap[v.ProductID][v.DeviceName] = struct{}{}
			}
		}
	}
	if in.GroupID != 0 {
		dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{GroupIDs: []int64{in.GroupID}}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range dis {
			if productMap[v.ProductID] == nil {
				productMap[v.ProductID] = map[string]struct{}{v.DeviceName: {}}
			} else {
				productMap[v.ProductID][v.DeviceName] = struct{}{}
			}
		}
	}
	var group errgroup.Group
	var newIn = dm.PropertyGetReportMultiSendReq{
		DataIDs: in.DataIDs,
	}
	var mu sync.Mutex
	var list = []*dm.PropertyGetReportSendMsg{}
	for k, v := range productMap {
		in2 := newIn
		in2.ProductID = k
		in2.DeviceNames = utils.SetToSlice(v)
		group.Go(func() error {
			logx.Errorf("开始")
			li, err := l.MultiSendOneProductProperty(&in2)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			list = append(list, li...)
			logx.Errorf("完成")
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return nil, err
	}
	return list, err
}
