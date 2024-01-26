package apply

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/area/info"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.UserAreaApplyIndexReq) (resp *types.UserAreaApplyIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserAreaApplyIndex(l.ctx, &sys.UserAreaApplyIndexReq{
		Page:      logic.ToSysPageRpc(req.Page),
		AuthTypes: req.AuthTypes,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.UserAreaApplyInfo
	var userInfoMap = map[int64]*sys.UserInfo{}
	var areaInfoMap = map[int64]*sys.AreaInfo{}
	for _, v := range ret.List {
		if req.WithUserInfo {
			userInfoMap[v.UserID] = nil
		}
		if req.WithAreaInfo {
			areaInfoMap[v.AreaID] = nil
		}
		list = append(list, &types.UserAreaApplyInfo{
			ID:          v.Id,
			UserID:      v.UserID,
			AreaID:      v.AreaID,
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime,
		})
	}
	if req.WithUserInfo {
		list, err := l.svcCtx.UserRpc.UserInfoIndex(l.ctx, &sys.UserInfoIndexReq{
			UserIDs: lo.Keys(userInfoMap),
		})
		if err != nil {
			return nil, err
		}
		for _, v := range list.List {
			userInfoMap[v.UserID] = v
		}
	}
	if req.WithAreaInfo {
		list, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{
			AreaIDs: lo.Keys(areaInfoMap),
		})
		if err != nil {
			return nil, err
		}
		for _, v := range list.List {
			areaInfoMap[v.AreaID] = v
		}
	}
	if req.WithAreaInfo || req.WithUserInfo {
		for _, v := range list {
			v.UserInfo = user.UserInfoToApi(userInfoMap[v.UserID], nil, nil)
			v.AreaInfo = info.ToAreaInfoTypes(areaInfoMap[v.AreaID], nil)
		}
	}
	return &types.UserAreaApplyIndexResp{
		Total: ret.Total,
		List:  list,
	}, nil
}
