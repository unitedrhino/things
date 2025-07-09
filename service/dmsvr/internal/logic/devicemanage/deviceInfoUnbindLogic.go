package devicemanagelogic

import (
	"context"
	"encoding/base64"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/product"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	devicemsglogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/devicemsg"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
	"gorm.io/gorm"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoUnbindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoUnbindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoUnbindLogic {
	return &DeviceInfoUnbindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoUnbindLogic) DeviceInfoUnbind(in *dm.DeviceInfoUnbindReq) (*dm.Empty, error) {
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	dev := devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}
	di, err := diDB.FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
		ProductID:   in.ProductID,
		DeviceNames: []string{in.DeviceName},
	})
	if err != nil {
		return nil, err
	}
	if di.ProjectID <= def.NotClassified {
		return nil, errors.DeviceNotBound.WithMsg("设备已解绑")
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(di.ProjectID)})
	if err != nil && !errors.Cmp(err, errors.NotFind) { //解绑的时候家庭已经不存在了也需要能正确解绑
		return nil, err
	}
	adminUserID := di.UserID
	if pi != nil {
		adminUserID = pi.AdminUserID
	}
	pc, err := l.svcCtx.ProductCache.GetData(l.ctx, di.ProductID)
	if err != nil {
		return nil, err
	}
	//如果是超管有全部权限
	if !uc.IsAdmin && (di.TenantCode != di.TenantCode || adminUserID != uc.UserID || int64(di.ProjectID) != uc.ProjectID) {
		switch pc.BindLevel {
		case product.BindLeveWeak3: //弱绑定谁都可以解绑
		case product.BindLeveMiddle2:
			if in.Signature == "" {
				return nil, errors.Permissions
			}
			var secert string
			ret, err := devicemsglogic.NewPropertyLogLatestIndexLogic(ctxs.WithRoot(l.ctx), l.svcCtx).PropertyLogLatestIndex(&dm.PropertyLogLatestIndexReq{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
				DataIDs:    []string{in.SecretType},
			})
			if err != nil {
				return nil, err
			}
			secert = ret.List[0].Value
			sig := getSignature(in.SignType, secert, fmt.Sprintf("deviceName=%s&nonce=%d&productID=%s&timestamp=%d",
				in.DeviceName, in.Nonce, in.ProductID, in.Timestamp))
			if sig != in.Signature {
				return nil, errors.Parameter.AddMsg("无效签名")
			}
		default:
			return nil, errors.Permissions
		}

	}
	//dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	//if err != nil {
	//	return nil, err
	//}
	di.TenantCode = def.TenantCodeDefault
	di.ProjectID = def.NotClassified
	di.UserID = def.RootNode
	di.AreaID = dataType.AreaID(def.NotClassified)
	di.AreaIDPath = def.NotClassifiedPath
	di.DeviceAlias = pc.ProductName
	if di.FirstBind.Valid && di.FirstBind.Time.After(time.Now().AddDate(0, 0, -1)) { //绑定一天内的不算绑定时间
		if pc.TrialTime != nil && di.ExpTime.Valid { //如果设备的有效期大于从当前算起的有效期,那说明充值过,这时候不能清除过期时间
			expTime := time.Now().Add(time.Hour * 24 * time.Duration(pc.TrialTime.GetValue()))
			if expTime.After(di.ExpTime.Time) {
				di.FirstBind.Valid = false
				di.ExpTime.Valid = false
			}
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewDeviceInfoRepo(tx).Update(ctxs.WithRoot(l.ctx), di)
		if err != nil {
			return err
		}
		err = relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx),
			relationDB.DeviceProfileFilter{Device: dev})
		if err != nil {
			return err
		}
		err = relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
			&dev,
		}})
		if err != nil {
			return err
		}
		err = logic.UpdateDevice(l.ctx, l.svcCtx, []*devices.Core{
			{ProductID: di.ProductID, DeviceName: di.DeviceName}}, devices.Affiliation{
			TenantCode: string(di.TenantCode), ProjectID: int64(di.ProjectID),
			AreaID: int64(di.AreaID), AreaIDPath: string(di.AreaIDPath)})
		return err
	})
	if err != nil {
		return nil, err
	}
	l.svcCtx.DeviceCache.SetData(l.ctx, dev, nil)
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, in.ProductID, in.DeviceName, DeleteModeThing)
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUnbind, &dev)
	if err != nil {
		l.Error(err)
	}
	BindChange(l.ctx, l.svcCtx, pc, dev, int64(di.ProjectID))

	if di.DeviceType == def.DeviceTypeGateway { //网关类型的需要解绑子设备
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			subs, err := relationDB.NewGatewayDeviceRepo(ctx).FindByFilter(l.ctx, relationDB.GatewayDeviceFilter{Gateway: &devices.Core{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
			}}, nil)
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return
			}
			for _, sub := range subs {
				_, err := NewDeviceInfoUnbindLogic(ctx, l.svcCtx).DeviceInfoUnbind(&dm.DeviceInfoUnbindReq{
					ProductID:  sub.ProductID,
					DeviceName: sub.DeviceName,
				})
				if err != nil {
					logx.WithContext(ctx).Error(err)
					continue
				}
			}
		})
	}
	return &dm.Empty{}, err
}

// 计算签名: 使用 HMAC-sha1 算法对目标串 dest 进行加密，密钥为 secret,将生成的结果进行 Base64 编码
func getSignature(sigType string, secret string, dest string) string {
	if secret == "" || dest == "" {
		return ""
	}
	str := utils.Hmac(sigType, dest, []byte(secret))
	return base64.StdEncoding.EncodeToString([]byte(str))
}
