package info

import (
	"context"
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

func (l *DeleteLogic) Delete(req *types.ProjectInfoDeleteReq) error {
	dmRep, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{
		Page: &dm.PageInfo{Page: 1, Size: 2}, //只需要知道是否有设备即可
	})
	if err != nil {
		return err
	}
	if len(dmRep.List) != 0 {
		return errors.Parameter.AddMsg("该项目已绑定了设备，不允许删除")
	}
	_, err = l.svcCtx.ProjectM.ProjectInfoDelete(l.ctx, &sys.ProjectInfoDeleteReq{ProjectID: req.ProjectID})
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
