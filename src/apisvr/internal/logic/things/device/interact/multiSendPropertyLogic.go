package interact

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"golang.org/x/sync/errgroup"
	"sync"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiSendPropertyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	retMsg []*types.DeviceInteractMultiSendPropertyMsg
	err    error
	mutex  sync.Mutex
}

func NewMultiSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiSendPropertyLogic {
	return &MultiSendPropertyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiSendPropertyLogic) MultiSendProperty(req *types.DeviceInteractMultiSendPropertyReq) (resp *types.DeviceInteractMultiSendPropertyResp, err error) {
	if req.ProductID != "" && len(req.DeviceNames) != 0 {
		err := l.SendProperty(req.ProductID, req.DeviceNames, req.Data, req.ShadowControl)
		return &types.DeviceInteractMultiSendPropertyResp{List: l.retMsg}, err
	}
	if req.GroupID != 0 || req.AreaID != 0 {
		var ds []*dm.DeviceInfo
		if req.GroupID != 0 {
			dgRet, err := l.svcCtx.DeviceG.GroupDeviceIndex(l.ctx, &dm.GroupDeviceIndexReq{
				GroupID: req.GroupID,
			})
			if err != nil {
				return nil, err
			}
			ds = dgRet.List
		}
		if req.AreaID != 0 {
			ret, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{
				AreaIDs: []int64{req.AreaID},
			})
			if err != nil {
				return nil, err
			}
			ds = ret.List
		}
		var devices = map[string][]string{} //key 是产品id value是设备名列表
		for _, v := range ds {
			if p := devices[v.ProductID]; p != nil {
				devices[v.ProductID] = append(p, v.DeviceName)
				continue
			}
			devices[v.ProductID] = []string{v.DeviceName}
		}
		for p, d := range devices {
			var eg errgroup.Group
			productID := p
			deviceNames := d
			eg.Go(func() error {
				err := l.SendProperty(productID, deviceNames, req.Data, req.ShadowControl)
				if err != nil {
					return err
				}
				return nil
			})
			err := eg.Wait()
			if err != nil {
				return nil, err
			}
		}
		return &types.DeviceInteractMultiSendPropertyResp{List: l.retMsg}, nil
	}
	return nil, errors.Parameter.AddMsg("产品id设备名或分组id或区域id必须填一个")
}
func (l *MultiSendPropertyLogic) SendProperty(productID string, deviceNames []string, data string, shadowControl int64) error {
	list := make([]*types.DeviceInteractMultiSendPropertyMsg, 0)
	dmReq := &dm.MultiSendPropertyReq{
		ProductID:     productID,
		DeviceNames:   deviceNames,
		Data:          data,
		ShadowControl: shadowControl,
	}
	dmResp, err := l.svcCtx.DeviceInteract.MultiSendProperty(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.MultiSendProperty productID=%v deviceNames=%v data=%v err=%+v", utils.FuncName(), productID, deviceNames, data, er)
		return er
	}
	if len(dmResp.List) > 0 {
		for _, v := range dmResp.List {
			list = append(list, &types.DeviceInteractMultiSendPropertyMsg{
				ProductID:  productID,
				DeviceName: v.DeviceName,
				Code:       v.Code,
				Msg:        v.Msg,
				MsgToken:   v.MsgToken,
				SysMsg:     v.SysMsg,
				SysCode:    v.SysCode,
			})
		}
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.retMsg = append(l.retMsg, list...)
	return nil
}
