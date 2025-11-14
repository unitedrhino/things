package staticEvent

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OneDayHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewOneDayHandle(ctx context.Context, svcCtx *svc.ServiceContext) *OneDayHandle {
	return &OneDayHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *OneDayHandle) Handle() error { //产品品类设备数量统计
	err := l.DeviceStatic()
	if err != nil {
		l.Error(err)
	}
	return nil
}

func (l *OneDayHandle) DeviceStatic() error { //区域下的设备数量统计
	{
		ret, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{})
		if err != nil {
			return err
		}
		var projectIDs []int64
		for _, v := range ret.List {
			projectIDs = append(projectIDs, v.ProjectID)
		}
		err = logic.DirectFillProjectDeviceCount(l.ctx, l.svcCtx, time.Millisecond*50, projectIDs...)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("DirectFillProjectDeviceCount error:%v", err)
		}
		time.Sleep(time.Second * 5) //休息一下减少波峰
	}
	{
		var total int64 = 9999 //如果三次都没有成功自然退出
		var size int64 = 500
		var areas []*sys.AreaInfo
		var errCount int64 = 0
		for page := int64(0); page*size < total; page++ {
			err := func() error {
				ret, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{Page: &sys.PageInfo{
					Page: page + 1,
					Size: size,
				}})
				if err != nil {
					return err
				}
				total = ret.Total
				for _, v := range ret.List {
					areas = append(areas, v)
				}
				return nil
			}()
			if err != nil {
				l.Error(err)
				errCount++
			}
			if errCount > 3 { //只有三次错误的机会
				break
			}
		}
		err := logic.DirectFillAreaDeviceCount(l.ctx, l.svcCtx, time.Millisecond*50, areas...)
		if err != nil {
			l.Error(err)
		}
	}

	return nil
}
