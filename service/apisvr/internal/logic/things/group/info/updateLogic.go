package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.GroupInfo) error {
	_, err := l.svcCtx.DeviceG.GroupInfoUpdate(l.ctx, ToGroupInfoPbTypes(req))

	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.GroupInfo.upadte failure err=%+v", utils.FuncName(), er)
		return er
	}
	return nil
}
