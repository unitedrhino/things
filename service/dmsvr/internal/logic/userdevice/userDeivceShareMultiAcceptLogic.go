package userdevicelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/share/dataType"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeivceShareMultiAcceptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeivceShareMultiAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeivceShareMultiAcceptLogic {
	return &UserDeivceShareMultiAcceptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 接受批量分享的设备
func (l *UserDeivceShareMultiAcceptLogic) UserDeivceShareMultiAccept(in *dm.UserDeviceShareMultiAcceptReq) (*dm.Empty, error) {
	multiDevices, err := l.svcCtx.UserMultiDeviceShare.GetData(l.ctx, in.ShareToken)
	if err != nil {
		return &dm.Empty{}, err
	}
	sharedDevices, _ := relationDB.NewUserDeviceShareRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserDeviceShareFilter{SharedUserID: in.SharedUserID}, nil)
	sharedDevicesMap := make(map[string]int64)
	for _, d := range sharedDevices {
		key := fmt.Sprintf("%s_%s", d.ProductID, d.DeviceName)
		sharedDevicesMap[key] = d.ID
	}
	acceptDevicesMap := make(map[string]bool)
	for _, v := range in.Devices {
		acceptDevicesMap[fmt.Sprintf("%s_%s", v.ProductID, v.DeviceName)] = true
	}
	tenantCode := ctxs.GetUserCtxNoNil(l.ctx).TenantCode
	for _, v := range multiDevices.Devices {
		key := fmt.Sprintf("%s_%s", v.ProductID, v.DeviceName)
		if !acceptDevicesMap[key] {
			continue
		}
		po := relationDB.DmUserDeviceShare{
			ProjectID:         multiDevices.ProjectID,
			TenantCode:        dataType.TenantCode(tenantCode),
			SharedUserID:      in.SharedUserID,
			SharedUserAccount: in.SharedUserAccount,
			ProductID:         v.ProductID,
			AuthType:          multiDevices.AuthType,
			DeviceName:        v.DeviceName,
			UseBy:             multiDevices.UseBy,
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
		l.svcCtx.UserDeviceShare.SetData(l.ctx, userShared.UserShareKey{
			ProductID:    po.ProductID,
			DeviceName:   po.DeviceName,
			SharedUserID: po.SharedUserID,
		}, nil)
	}
	return &dm.Empty{}, nil
}
