package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareIndexLogic {
	return &UserDeviceShareIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备分享列表(只有)
func (l *UserDeviceShareIndexLogic) UserDeviceShareIndex(in *dm.UserDeviceShareIndexReq) (*dm.UserDeviceShareIndexResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.Device == nil {
		return nil, errors.Parameter.AddMsg("设备需要填写")
	}
	di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithAllProject(l.ctx), relationDB.DeviceFilter{ProductID: in.Device.ProductID, DeviceNames: []string{in.Device.DeviceName}})
	if err != nil {
		return nil, err
	}
	if di.UserID <= def.RootNode || di.ProjectID <= def.NotClassified {
		return &dm.UserDeviceShareIndexResp{}, nil
	}
	pi, err := l.svcCtx.ProjectCache.GetData(l.ctx, int64(di.ProjectID))
	if err != nil {
		return nil, err
	}
	if !uc.IsAdmin && (pi.AdminUserID != uc.UserID) { //只有所有者和被分享者才有权限操作
		return nil, errors.Permissions
	}
	f := relationDB.UserDeviceShareFilter{ProductID: di.ProductID, DeviceName: di.DeviceName}
	total, err := relationDB.NewUserDeviceShareRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &dm.UserDeviceShareIndexResp{Total: total, List: ToUserDeviceSharePbs(list)}, nil
}
