package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoUpdateLogic {
	return &TenantInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新区域
func (l *TenantInfoUpdateLogic) TenantInfoUpdate(in *sys.TenantInfo) (*sys.Response, error) {
	repo := relationDB.NewTenantInfoRepo(l.ctx)
	old, err := repo.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.AdminUserID != 0 && in.AdminUserID != old.AdminUserID { //只有default的超管才有权限修改管理员
		err := logic.IsSupperAdmin(l.ctx, def.TenantCodeDefault)
		if err != nil {
			return nil, err
		}
		old.AdminUserID = in.AdminUserID
	}
	if in.BaseUrl != "" {
		old.BaseUrl = in.BaseUrl
	}
	if in.LogoUrl != "" {
		old.LogoUrl = in.LogoUrl
	}
	if in.Desc != nil {
		old.Desc = utils.ToEmptyString(in.Desc)
	}
	err = repo.Update(l.ctx, old)

	return &sys.Response{}, err
}
