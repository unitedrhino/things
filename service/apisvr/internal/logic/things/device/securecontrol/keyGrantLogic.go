package securecontrol

import (
	"context"
	stderrors "errors"
	"strings"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	shareerr "gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

// KeyGrantLogic 负责校验用户控制权限并派生 App 安全控制 key。
type KeyGrantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewKeyGrantLogic 创建 App 安全控制 key-grant 逻辑对象。
func NewKeyGrantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeyGrantLogic {
	return &KeyGrantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// KeyGrant 返回当前 App 会话可使用的 App 安全控制 key 或 wrappedKey。
func (l *KeyGrantLogic) KeyGrant(req *types.DeviceSecureControlKeyGrantReq) (*types.DeviceSecureControlKeyGrantResp, error) {
	if strings.TrimSpace(req.ProductID) == "" || strings.TrimSpace(req.DeviceName) == "" {
		return nil, shareerr.Parameter.WithMsg("productID和deviceName必填")
	}
	purpose := strings.TrimSpace(req.Purpose)
	if purpose != "" && purpose != VersionAppSecControlV1 {
		return nil, shareerr.Parameter.WithMsg("purpose不支持")
	}
	dev, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		return nil, shareerr.Fmt(err)
	}
	if err := l.checkGrantPermission(dev); err != nil {
		return nil, err
	}
	if dev.Secret == "" {
		return nil, shareerr.Parameter.WithMsg("设备密钥为空")
	}
	grant, err := BuildKeyGrant(GrantInput{
		ProductID:     dev.ProductID,
		DeviceName:    dev.DeviceName,
		DeviceSecret:  dev.Secret,
		KeyVersion:    req.KeyVersion,
		AuthAlg:       req.AuthAlg,
		AppKeyWrapPub: strings.TrimSpace(req.AppKeyWrapPub),
	})
	if err != nil {
		return nil, grantError(err)
	}
	return &types.DeviceSecureControlKeyGrantResp{
		Version:       grant.Version,
		ProductID:     grant.ProductID,
		DeviceName:    grant.DeviceName,
		AuthAlg:       grant.AuthAlg,
		KeyVersion:    grant.KeyVersion,
		KeyEncoding:   grant.KeyEncoding,
		AppControlKey: grant.AppControlKey,
		WrappedKey:    grant.WrappedKey,
		WrapNonce:     grant.WrapNonce,
		ServerWrapPub: grant.ServerWrapPub,
		ScopeBits:     grant.ScopeBits,
		TtlSec:        grant.TTLSec,
	}, nil
}

// checkGrantPermission 收紧 key-grant 权限到项目读写及以上或设备分享全权限。
func (l *KeyGrantLogic) checkGrantPermission(dev *dm.DeviceInfo) error {
	uc := ctxs.GetUserCtx(l.ctx)
	if uc == nil || (!uc.IsAdmin && !uc.IsSuperAdmin && uc.UserID == 0) {
		return shareerr.Permissions
	}
	if uc.AuthProject(dev.ProjectID, def.AuthReadWrite) {
		return nil
	}
	share, err := l.svcCtx.UserDevice.UserDeviceShareRead(l.ctx, &dm.UserDeviceShareReadReq{
		Device: &dm.DeviceCore{
			ProductID:  dev.ProductID,
			DeviceName: dev.DeviceName,
		},
	})
	if err != nil {
		return shareerr.Permissions
	}
	if share.AuthType == def.AuthAdmin {
		return nil
	}
	return shareerr.Permissions.AddMsg("没有设备安全控制授权")
}

// grantError 把内部协议错误映射为对外参数错误。
func grantError(err error) error {
	switch {
	case stderrors.Is(err, ErrUnsupportedAuthAlg):
		return shareerr.Parameter.WithMsg("authAlg不支持")
	case stderrors.Is(err, ErrInvalidDeviceSecret):
		return shareerr.Parameter.WithMsg("设备密钥格式不正确")
	case stderrors.Is(err, ErrInvalidKeyVersion):
		return shareerr.Parameter.WithMsg("keyVersion不正确")
	case stderrors.Is(err, ErrInvalidWrapPublicKey):
		return shareerr.Parameter.WithMsg("appKeyWrapPub格式不正确")
	case stderrors.Is(err, ErrInvalidWrapNonce):
		return shareerr.Parameter.WithMsg("wrapNonce不正确")
	default:
		return err
	}
}
