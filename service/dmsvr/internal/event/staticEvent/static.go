package staticEvent

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type StaticHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewStaticHandle(ctx context.Context, svcCtx *svc.ServiceContext) *StaticHandle {
	return &StaticHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *StaticHandle) Handle() error { //产品品类设备数量统计
	w := sync.WaitGroup{}
	w.Add(3)
	utils.Go(l.ctx, func() {
		err := l.ProductCategoryStatic()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		err := l.AreaDeviceStatic()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		err := l.DeviceExp()
		if err != nil {
			l.Error(err)
		}
	})
	w.Wait()
	return nil
}
func (l *StaticHandle) AreaDeviceStatic() error { //区域下的设备数量统计
	ret, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{})
	if err != nil {
		return err
	}
	var areaPaths []string
	for _, v := range ret.List {
		areaPaths = append(areaPaths, v.AreaIDPath)
	}
	err = logic.FillAreaDeviceCount(l.ctx, l.svcCtx, areaPaths...)
	return err
}

func (l *StaticHandle) DeviceExp() error { //设备过期处理
	{ //有效期到了之后不启用
		err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx,
			relationDB.DeviceFilter{ExpTime: stores.CmpAnd(stores.CmpLte(time.Now()), stores.CmpIsNull(false))},
			map[string]any{"is_enable": def.False})
		if err != nil {
			l.Error(err)
		}
	}
	{ //清除设置了过期时间且过期了的分享
		err := relationDB.NewUserDeviceShareRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ExpTime: stores.CmpAnd(stores.CmpLte(time.Now()), stores.CmpIsNull(false)),
		})
		if err != nil {
			l.Error(err)
		}
	}
	return nil
}

func (l *StaticHandle) ProductCategoryStatic() error { //产品品类设备数量统计
	pcDB := relationDB.NewProductCategoryRepo(l.ctx)
	pcs, err := pcDB.FindByFilter(l.ctx, relationDB.ProductCategoryFilter{}, nil)
	if err != nil {
		return err
	}
	for _, pc := range pcs {
		ids := utils.GetIDPath(pc.IDPath)
		total, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{ProductCategoryIDs: ids})
		if err != nil {
			l.Error(err)
			continue
		}
		pc.DeviceCount = total
		err = pcDB.Update(l.ctx, pc)
		if err != nil {
			l.Error(err)
			continue
		}
	}
	return nil
}
