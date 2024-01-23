package info

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.UserInfoReadReq) (resp *types.UserInfo, err error) {
	info, err := l.svcCtx.UserRpc.UserInfoRead(l.ctx, &sys.UserInfoReadReq{UserID: req.UserID})
	if err != nil {
		return nil, err
	}
	var (
		roles  []*sys.RoleInfo
		tenant *sys.TenantInfo
	)
	if req.WithRoles == true {
		ret, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
			UserID: req.UserID,
		})
		if err != nil {
			return nil, err
		}
		roles = ret.List
	}
	if req.WithTenant == true {
		ret, err := l.svcCtx.TenantRpc.TenantInfoRead(l.ctx, &sys.WithIDCode{Code: ctxs.GetUserCtx(l.ctx).TenantCode})
		if err != nil {
			return nil, err
		}
		tenant = ret
	}
	//if req.WithAreas {
	//	ret, err := l.svcCtx.UserRpc.UserAreaIndex(l.ctx, &sys.UserAreaIndexReq{
	//		UserID: req.UserID,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	if len(ret.List) != 0 {
	//		var areaIDs []int64
	//		for _, v := range ret.List {
	//			areaIDs = append(areaIDs, v.AreaID)
	//		}
	//		ret2, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{AreaIDs: areaIDs})
	//		if err != nil {
	//			return nil, err
	//		}
	//		areas = ret2.List
	//	}
	//}

	return user.UserInfoToApi(info, roles, tenant), nil
}
