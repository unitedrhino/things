package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"golang.org/x/sync/errgroup"
	"sync"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiSendPropertyControlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMultiSendPropertyControlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiSendPropertyControlLogic {
	return &MultiSendPropertyControlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量调用设备属性
func (l *MultiSendPropertyControlLogic) MultiSendPropertyControl(in *dm.MultiSendPropertyControlReq) (*dm.MultiSendPropertyControlResp, error) {
	var list []*dm.SendPropertyControlMsg
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
	return &dm.MultiSendPropertyControlResp{List: list}, nil
}

func (l *MultiSendPropertyControlLogic) MultiSendOneProductProperty(in *dm.MultiSendPropertyControlReq) ([]*dm.SendPropertyControlMsg, error) {
	list := make([]*dm.SendPropertyControlMsg, 0)
	sigSend := NewSendPropertyControlLogic(l.ctx, l.svcCtx)
	err := sigSend.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, v := range in.DeviceNames {
		wg.Add(1)
		go func(v string) {
			defer utils.Recover(l.ctx)
			defer wg.Done()
			ret, err := sigSend.SendPropertyControl(&dm.SendPropertyControlReq{
				ProductID:  in.ProductID,
				DeviceName: v,
				Data:       in.Data,
				IsAsync:    false,
			})

			if err != nil {
				myErr, _ := err.(*errors.CodeError)
				msg := &dm.SendPropertyControlMsg{
					DeviceName: v,
					SysMsg:     myErr.GetMsg(),
					SysCode:    myErr.Code,
				}
				mu.Lock()
				defer mu.Unlock()
				list = append(list, msg)
				return
			}

			msg := &dm.SendPropertyControlMsg{
				ProductID:  in.ProductID,
				DeviceName: v,
				SysCode:    errors.OK.Code,
				SysMsg:     errors.OK.GetMsg(),
				Code:       ret.Code,
				Msg:        ret.Msg,
				MsgToken:   ret.MsgToken,
			}
			mu.Lock()
			defer mu.Unlock()
			list = append(list, msg)
		}(v)
	}

	wg.Wait()
	return list, err
}

func (l *MultiSendPropertyControlLogic) MultiSendMultiProductProperty(in *dm.MultiSendPropertyControlReq) ([]*dm.SendPropertyControlMsg, error) {
	var productMap = map[string][]string{} //key是产品id,value是产品下的设备列表
	for _, v := range in.Devices {
		productMap[v.ProductID] = append(productMap[v.ProductID], v.DeviceName)
	}
	var group errgroup.Group
	var newIn = dm.MultiSendPropertyControlReq{
		ShadowControl: in.ShadowControl,
		Data:          in.Data,
	}
	var mu sync.Mutex
	var list = []*dm.SendPropertyControlMsg{}
	for k, v := range productMap {
		in := newIn
		in.ProductID = k
		in.DeviceNames = v
		group.Go(func() error {
			li, err := l.MultiSendOneProductProperty(&in)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			list = append(list, li...)
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return nil, err
	}
	return list, err
}
