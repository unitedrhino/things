package userdevicelogic

import (
	"context"
	"fmt"
	"time"

	"gitee.com/unitedrhino/core/share/dataType"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	shareerrors "gitee.com/unitedrhino/share/errors"
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

type multiAcceptDeviceProjectLookup func(productID string, deviceName string) (int64, error)

// validateMultiAcceptDevices 校验接收时设备仍归属于生成分享 Token 时的项目。
func validateMultiAcceptDevices(
	tokenProjectID int64,
	tokenDevices []*dm.DeviceShareInfo,
	requestedDevices []*dm.DeviceCore,
	lookup multiAcceptDeviceProjectLookup,
) error {
	if tokenProjectID <= def.NotClassified {
		return shareerrors.DeviceNotBound.WithMsg("设备已解绑，分享已失效")
	}

	tokenDeviceSet := make(map[string]struct{}, len(tokenDevices))
	for _, device := range tokenDevices {
		tokenDeviceSet[fmt.Sprintf("%s_%s", device.ProductID, device.DeviceName)] = struct{}{}
	}
	for _, device := range requestedDevices {
		key := fmt.Sprintf("%s_%s", device.ProductID, device.DeviceName)
		if _, ok := tokenDeviceSet[key]; !ok {
			return shareerrors.Permissions.WithMsg("设备不在当前分享中")
		}
		currentProjectID, err := lookup(device.ProductID, device.DeviceName)
		if err != nil {
			if shareerrors.Cmp(err, shareerrors.NotFind) {
				return shareerrors.DeviceNotBound.WithMsg("设备已解绑，分享已失效")
			}
			return err
		}
		if currentProjectID <= def.NotClassified || currentProjectID != tokenProjectID {
			return shareerrors.DeviceNotBound.WithMsg("设备已解绑或已更换主人，分享已失效")
		}
	}
	return nil
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
	err = validateMultiAcceptDevices(multiDevices.ProjectID, multiDevices.Devices, in.Devices, func(productID string, deviceName string) (int64, error) {
		device, findErr := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(
			ctxs.WithRoot(l.ctx),
			relationDB.DeviceFilter{
				ProductID:   productID,
				DeviceNames: []string{deviceName},
			},
		)
		if findErr != nil {
			return 0, findErr
		}
		return int64(device.ProjectID), nil
	})
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
	acceptedCount := 0
	acceptedDevices := make([]*dm.DeviceShareInfo, 0, len(in.Devices))
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
		acceptedCount++
		acceptedDevices = append(acceptedDevices, v)
	}
	if acceptedCount > 0 {
		event := buildDeviceShareAcceptedEvent(in, multiDevices, tenantCode, acceptedDevices, time.Now().Unix())
		err = publishDeviceShareAcceptedEvent(l.ctx, l.svcCtx.FastEvent, event)
		if err != nil {
			return &dm.Empty{}, err
		}
	}
	if acceptedCount > 0 && shouldConsumeShareTokenAfterAccept(multiDevices.UseBy) {
		err = l.svcCtx.UserMultiDeviceShare.DeleteToken(l.ctx, tenantCode, multiDevices.UserID, in.ShareToken)
		if err != nil {
			return &dm.Empty{}, err
		}
	}
	return &dm.Empty{}, nil
}
