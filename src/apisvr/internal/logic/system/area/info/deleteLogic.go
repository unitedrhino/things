package info

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func getAreaIDs(in *sys.AreaInfo, areaIDs []int64) []int64 {
	areaIDs = append(areaIDs, in.AreaID)
	for _, v := range in.Children {
		areaIDs = getAreaIDs(v, areaIDs)
	}
	return areaIDs
}

func (l *DeleteLogic) Delete(req *types.AreaWithID) error {
	treeResp, err := l.svcCtx.AreaM.AreaInfoRead(l.ctx, &sys.AreaInfoReadReq{AreaID: req.AreaID, IsRetTree: true})
	if err != nil {
		return err
	}
	var areaIDs []int64 = getAreaIDs(treeResp, nil)
	dmRep, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{
		AreaIDs: areaIDs})
	if err != nil {
		return err
	}
	if len(dmRep.List) != 0 {
		var devices []*dm.DeviceCore
		for _, v := range dmRep.List {
			devices = append(devices, &dm.DeviceCore{
				DeviceName: v.DeviceName,
				ProductID:  v.ProductID,
			})
		}
		//全部放到未分类中
		_, err := l.svcCtx.DeviceM.DeviceInfoMultiUpdate(l.ctx, &dm.DeviceInfoMultiUpdateReq{
			Devices: devices,
			AreaID:  def.NotClassified,
		})
		if err != nil {
			return err
		}
	}
	_, err = l.svcCtx.AreaM.AreaInfoDelete(l.ctx, &sys.AreaWithID{AreaID: req.AreaID})
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
