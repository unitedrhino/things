package deviceinteractlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
	return &dm.MultiSendPropertyResp{List: list}, nil
}
