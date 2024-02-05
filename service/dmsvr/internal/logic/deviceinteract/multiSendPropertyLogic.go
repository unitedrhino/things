package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"golang.org/x/sync/errgroup"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiSendPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMultiSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiSendPropertyLogic {
	return &MultiSendPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量调用设备属性
func (l *MultiSendPropertyLogic) MultiSendProperty(in *dm.MultiSendPropertyReq) (*dm.MultiSendPropertyResp, error) {
	var list []*dm.SendPropertyMsg
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
	return &dm.MultiSendPropertyResp{List: list}, nil
}

func (l *MultiSendPropertyLogic) MultiSendOneProductProperty(in *dm.MultiSendPropertyReq) ([]*dm.SendPropertyMsg, error) {
	list := make([]*dm.SendPropertyMsg, 0)
	sigSend := NewSendPropertyLogic(l.ctx, l.svcCtx)
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
			ret, err := sigSend.SendProperty(&dm.SendPropertyReq{
				ProductID:  in.ProductID,
				DeviceName: v,
				Data:       in.Data,
				IsAsync:    false,
			})

			if err != nil {
				myErr, _ := err.(*errors.CodeError)
				msg := &dm.SendPropertyMsg{
					DeviceName: v,
					SysMsg:     myErr.Msg,
					SysCode:    myErr.Code,
				}
				mu.Lock()
				defer mu.Unlock()
				list = append(list, msg)
				return
			}

			msg := &dm.SendPropertyMsg{
				ProductID:  in.ProductID,
				DeviceName: v,
				SysCode:    errors.OK.Code,
				SysMsg:     errors.OK.Msg,
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

func (l *MultiSendPropertyLogic) MultiSendMultiProductProperty(in *dm.MultiSendPropertyReq) ([]*dm.SendPropertyMsg, error) {
	var productMap = map[string][]string{} //key是产品id,value是产品下的设备列表
	for _, v := range in.Devices {
		productMap[v.ProductID] = append(productMap[v.ProductID], v.DeviceName)
	}
	var group errgroup.Group
	var newIn = dm.MultiSendPropertyReq{
		ShadowControl: in.ShadowControl,
		Data:          in.Data,
	}
	var mu sync.Mutex
	var list = []*dm.SendPropertyMsg{}
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
