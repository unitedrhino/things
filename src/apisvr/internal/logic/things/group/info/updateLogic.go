package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *UpdateLogic) Update(req *types.GroupInfoUpdateReq) error {
	_, err := l.svcCtx.DeviceG.GroupInfoUpdate(l.ctx, &dm.GroupInfoUpdateReq{
		GroupID:   req.GroupID,
		GroupName: *req.GroupName,
		Desc:      *req.Desc,
		Tags:      toTagsMap(req.Tags),
	})

	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.GroupInfo.upadte failure err=%+v", utils.FuncName(), er)
		return er
	}
	return nil
}
