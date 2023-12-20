package apimanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.ApiInfoRepo
}

func NewApiInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoCreateLogic {
	return &ApiInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiInfoCreateLogic) ApiInfoCreate(in *sys.ApiInfo) (*sys.WithID, error) {
	po := relationDB.SysApiInfo{
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		Group:        in.Group,
		AppCode:      in.AppCode,
	}
	err := l.AiDB.Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
