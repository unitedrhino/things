package userdevicelogic

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMultiDeivcesShareAcceptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMultiDeivcesShareAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMultiDeivcesShareAcceptLogic {
	return &UserMultiDeivcesShareAcceptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 接受批量分享的设备
func (l *UserMultiDeivcesShareAcceptLogic) UserMultiDeivcesShareAccept(in *dm.UserMultiDevicesShareAcceptReq) (*dm.Empty, error) {
	multiDevices, err := l.svcCtx.UserMultiDeviceShare.GetData(l.ctx, in.Keyword)
	if err != nil {
		return &dm.Empty{}, err
	}
	sharedDevices, _ := relationDB.NewUserDeviceShareRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserDeviceShareFilter{SharedUserID: in.SharedUserID}, nil)
	sharedDevicesMap := make(map[string]int64)
	for _, d := range sharedDevices {
		key := fmt.Sprintf("%s_%s", d.ProductID, d.DeviceName)
		sharedDevicesMap[key] = d.ID
	}
	tenantCode := ctxs.GetUserCtxNoNil(l.ctx).TenantCode
	for _, v := range multiDevices.Device {
		key := fmt.Sprintf("%s_%s", v.ProductID, v.DeviceName)
		po := relationDB.DmUserDeviceShare{
			ProjectID:         multiDevices.ProjectID,
			TenantCode:        stores.TenantCode(tenantCode),
			SharedUserID:      in.SharedUserID,
			SharedUserAccount: in.SharedUserAccount,
			ProductID:         v.ProductID,
			AuthType:          multiDevices.AuthType,
			DeviceName:        v.DeviceName,
			AccessPerm:        utils.CopyMap[relationDB.SharePerm](multiDevices.AccessPerm),
			SchemaPerm:        utils.CopyMap[relationDB.SharePerm](multiDevices.SchemaPerm),
			ExpTime:           utils.ToNullTime(multiDevices.ExpTime),
		}
		if po.AccessPerm == nil {
			po.AccessPerm = map[string]*relationDB.SharePerm{}
		}
		if po.SchemaPerm == nil {
			po.SchemaPerm = map[string]*relationDB.SharePerm{}
		}
		if id, ok := sharedDevicesMap[key]; ok {
			po.ID = id
			if err := relationDB.NewUserDeviceShareRepo(l.ctx).Update(l.ctx, &po); err != nil {
				return &dm.Empty{}, err
			}
		} else {
			err = relationDB.NewUserDeviceShareRepo(l.ctx).Insert(l.ctx, &po)
			if err != nil {
				return &dm.Empty{}, err
			}
		}
	}
	return &dm.Empty{}, nil
}
