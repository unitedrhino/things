package accessmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoCreateLogic {
	return &ApiInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoCreateLogic) ApiInfoCreate(in *sys.ApiInfo) (*sys.WithID, error) {
	_, err := relationDB.NewAccessRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AccessFilter{Codes: []string{in.AccessCode}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("接口区域不存在")
		}
	}
	po := ToApiInfoPo(in)
	po.ID = 0
	err = relationDB.NewApiInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
