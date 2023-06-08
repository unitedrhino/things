package deviceinteractlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"
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
func (l *MultiSendPropertyLogic) MultiSendProperty(in *di.MultiSendPropertyReq) (*di.MultiSendPropertyResp, error) {
	list := make([]*di.SendPropertyMsg, 0)
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
			ret, err := sigSend.SendProperty(&di.SendPropertyReq{
				ProductID:  in.ProductID,
				DeviceName: v,
				Data:       in.Data,
				IsAsync:    false,
			})

			if err != nil {
				myErr, _ := err.(*errors.CodeError)
				msg := &di.SendPropertyMsg{
					SysMsg:  err.Error(),
					SysCode: myErr.Code,
				}
				mu.Lock()
				defer mu.Unlock()
				list = append(list, msg)
				return
			}

			msg := &di.SendPropertyMsg{
				Code:        ret.Code,
				Status:      ret.Status,
				ClientToken: ret.ClientToken,
			}

			mu.Lock()
			defer mu.Unlock()
			list = append(list, msg)
		}(v)
	}

	wg.Wait()
	return &di.MultiSendPropertyResp{List: list}, nil
}
