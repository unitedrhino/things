package deviceauthlogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	devicemanagelogic "github.com/i-Things/things/src/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceRegisterLogic {
	return &DeviceRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设备动态注册
func (l *DeviceRegisterLogic) DeviceRegister(in *dm.DeviceRegisterReq) (*dm.DeviceRegisterResp, error) {

	//1.检查产品动态开关是否开启
	pi, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, errors.Database.AddMsg(fmt.Sprintf("设备查询失败: %s", err.Error()))
	}
	if pi.AutoRegister == 1 {
		return nil, errors.NotEnable.AddMsg("设备动态注册未启动")
	}

	params := map[string]string{
		"productId":  in.ProductID,
		"deviceName": in.DeviceName,
		"nonce":      cast.ToString(in.Nonce),
		"timestamp":  cast.ToString(in.Timestamp),
	}
	//2.用产品id获取产品密钥，按照算法计算设备签名，与设备发来的签名对比是否一致
	sig := utils.GetSignature(pi.Secret, params)
	if sig != in.Signature {
		return nil, errors.Permissions.AddMsg("无效签名")
	}
	//3.如果一致则是有效设备，获取设备名，查看是否已经注册过，如果未注册过，则录入设备信息
	_, err = l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err == nil {
		return nil, errors.Default.AddMsg("设备已注册")
	}
	_, err = devicemanagelogic.NewDeviceInfoCreateLogic(l.ctx, l.svcCtx).DeviceInfoCreate(&dm.DeviceInfo{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		LogLevel:   1,
	})
	if err != nil {
		return nil, errors.Database.AddMsg(fmt.Sprintf("设备注册失败: %s", err.Error()))
	}
	resp, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		return nil, errors.Database.AddMsg(fmt.Sprintf("设备注册失败: %s", err.Error()))
	}

	//4.给设备返回设备密钥 经过aes cbc base64加密
	psk, err := utils.AesCbcBase64(resp.Secret, pi.Secret)
	if err != nil {
		return nil, errors.Default.AddMsg(fmt.Sprintf("设备密钥加密失败: %s", err.Error()))
	}
	return &dm.DeviceRegisterResp{Psk: psk}, nil
}
