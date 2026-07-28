package userdevicelogic

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/domain/usershare"
	"github.com/spf13/cast"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareCreateLogic {
	return &UserDeviceShareCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// invalidateCreatedUserDeviceShare 清除新建分享对应的旧权限缓存。
func invalidateCreatedUserDeviceShare(
	ctx context.Context,
	setData func(context.Context, userShared.UserShareKey, *dm.UserDeviceShareInfo) error,
	share *relationDB.DmUserDeviceShare,
) {
	_ = setData(ctx, userShared.UserShareKey{
		ProductID:    share.ProductID,
		DeviceName:   share.DeviceName,
		SharedUserID: share.SharedUserID,
	}, nil)
}

// 分享设备
func (l *UserDeviceShareCreateLogic) UserDeviceShareCreate(in *dm.UserDeviceShareInfo) (*dm.WithID, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID != 0 && (uc.IsAdmin || uc.ProjectAuth[in.ProjectID] != nil) {
		l.ctx = ctxs.WithProjectID(l.ctx, in.ProjectID)
	}
	di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.Device.ProductID, DeviceNames: []string{in.Device.DeviceName}})
	if err != nil {
		return nil, err
	}
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(di.ProjectID)})
	if err != nil {
		return nil, err
	}
	if pi.AdminUserID != uc.UserID && !uc.IsAdmin {
		pa := uc.ProjectAuth[pi.ProjectID]
		if pa.AuthType != def.AuthAdmin {
			if pa.Area == nil || pa.Area[int64(di.AreaID)] != def.AuthAdmin {
				return nil, errors.Permissions.AddMsg("只有管理员才能分享设备")
			}
		}
	}
	if in.SharedUserID == uc.UserID {
		return nil, errors.Parameter.AddMsg("不能分享给自己")
	}
	ui, err := l.svcCtx.UserM.UserInfoRead(l.ctx, &sys.UserInfoReadReq{UserID: in.SharedUserID})
	if err != nil {
		return nil, err
	}
	var account = ui.UserName
	if account == "" {
		account = ui.Phone.GetValue()
	}
	if account == "" {
		account = ui.Email.GetValue()
	}
	if account == "" {
		account = cast.ToString(ui.UserID)
	}
	po := relationDB.DmUserDeviceShare{
		ProjectID:         pi.ProjectID,
		SharedUserID:      in.SharedUserID,
		SharedUserAccount: account,
		ProductID:         in.Device.ProductID,
		UseBy:             in.UseBy,
		Desc:              in.Desc,
		AuthType:          in.AuthType,
		DeviceName:        in.Device.DeviceName,
		AccessPerm:        utils.CopyMap[relationDB.SharePerm](in.AccessPerm),
		SchemaPerm:        utils.CopyMap[relationDB.SharePerm](in.SchemaPerm),
		ExpTime:           utils.ToNullTime2(in.ExpTime),
	}
	if po.AccessPerm == nil {
		po.AccessPerm = map[string]*relationDB.SharePerm{}
	}
	if po.SchemaPerm == nil {
		po.SchemaPerm = map[string]*relationDB.SharePerm{}
	}
	err = relationDB.NewUserDeviceShareRepo(l.ctx).Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	invalidateCreatedUserDeviceShare(l.ctx, l.svcCtx.UserDeviceShare.SetData, &po)
	event := buildDeviceShareGrantedEvent(
		usershare.DeviceShareGrantSourceAccountDirect,
		"",
		uc.UserID,
		po.SharedUserID,
		po.SharedUserAccount,
		po.ProjectID,
		uc.TenantCode,
		po.UseBy,
		[]usershare.DeviceShareGrantedDevice{
			{
				ShareID:    po.ID,
				ProductID:  po.ProductID,
				DeviceName: po.DeviceName,
			},
		},
		time.Now().Unix(),
	)
	err = publishDeviceShareGrantedEvent(l.ctx, l.svcCtx.FastEvent, event)
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: po.ID}, nil
}
