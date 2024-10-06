package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	devicemanagelogic "gitee.com/i-Things/things/service/dmsvr/internal/logic/devicemanage"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"sync"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInitLogic {
	return &ProductInitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductInitLogic) ProductInit(in *dm.ProductInitReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}

	f := relationDB.ProductFilter{ProductIDs: in.ProductIDs}
	pis, err := relationDB.NewProductInfoRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	var wait sync.WaitGroup
	for _, v := range pis {
		pi := v
		wait.Add(1)
		utils.Go(l.ctx, func() {
			defer func() {
				wait.Done()
			}()
			err := l.initOne(pi)
			if err != nil {
				logx.Error(pi, err)
			}
		})
	}
	wait.Wait()
	return &dm.Empty{}, nil
}
func (l *ProductInitLogic) initOne(in *relationDB.DmProductInfo) error {
	err := NewProductInfoCreateLogic(l.ctx, l.svcCtx).InitProduct(in)
	if err != nil {
		logx.Error(in, err)
		return err
	}
	{ //物模型初始化
		t, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
		if err != nil {
			l.Errorf("%s.SchemaManaRepo.GetSchemaModel failure,err:%v", utils.FuncName(), err)
			return err
		}
		if err := l.svcCtx.SchemaManaRepo.InitProduct(l.ctx, t, in.ProductID); err != nil {
			l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
			return err
		}
	}

	dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID}, nil)
	if err != nil {
		logx.Error(in, err)
		return err
	}
	dic := devicemanagelogic.NewDeviceInfoCreateLogic(l.ctx, l.svcCtx)
	for _, di := range dis {
		err := dic.InitDevice(devices.Info{
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
			TenantCode: string(di.TenantCode),
			ProjectID:  int64(di.ProjectID),
			AreaID:     int64(di.AreaID),
		})
		if err != nil {
			logx.Error(in, di, err)
		}
	}
	return nil
}
