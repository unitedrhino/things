package deviceauthlogic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	devicemanagelogic "github.com/i-Things/things/src/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
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

// 计算签名: 使用 HMAC-sha1 算法对目标串 dest 进行加密，密钥为 secret,将生成的结果进行 Base64 编码
func getSignature(secret string, dest string) string {
	if secret == "" || dest == "" {
		return ""
	}

	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha1(dest, []byte(secret))))
}

// 设备动态注册
func (l *DeviceRegisterLogic) DeviceRegister(in *dm.DeviceRegisterReq) (*dm.DeviceRegisterResp, error) {

	//检查产品动态开关是否开启
	pi, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.Database.AddMsg("产品查询失败")
	}
	if pi.AutoRegister == devices.DeviceRegisterUnable {
		return nil, errors.NotEnable.AddMsg("设备动态注册未启动")
	}

	//检查设备是否已注册
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			//检查设备自动创建是否开启， 开启则自动创建设备，未开启则返回错误
			if pi.AutoRegister == devices.DeviceAutoCreateEnable {
				//检查设备签名是否正确
				sig := getSignature(pi.Secret, fmt.Sprintf("deviceName=%s&nonce=%d&productId=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp))
				if sig != in.Signature {
					return nil, errors.Parameter.AddMsg("无效签名")
				}
				_, err = devicemanagelogic.NewDeviceInfoCreateLogic(l.ctx, l.svcCtx).DeviceInfoCreate(&dm.DeviceInfo{
					ProductID:  in.ProductID,
					DeviceName: in.DeviceName,
				})
				if err != nil {
					return nil, errors.Database.AddMsg(fmt.Sprintf("设备注册失败: %s", err.Error()))
				}
				resp, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
				if err != nil {
					return nil, errors.Database.AddMsg(fmt.Sprintf("设备注册失败: %s", err.Error()))
				}
				return &dm.DeviceRegisterResp{Psk: resp.Secret}, nil
			}
			return nil, errors.NotFind.AddMsg("设备注册失败，无效设备")
		} else {
			return nil, errors.Database.AddMsg(fmt.Sprintf("设备注册失败: %s", err.Error()))
		}
	}

	if di.FirstLogin.Valid == true {
		return nil, errors.NotEmpty.AddMsg("设备已注册")
	}

	//检查设备签名是否正确
	sig := getSignature(pi.Secret, fmt.Sprintf("deviceName=%s&nonce=%d&productId=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp))
	if sig != in.Signature {
		return nil, errors.Parameter.AddMsg("无效签名")
	}

	//给设备返回设备密钥 经过aes cbc base64加密
	psk, err := utils.AesCbcBase64(di.Secret, pi.Secret)
	if err != nil {
		return nil, errors.Default.AddMsg(fmt.Sprintf("设备密钥加密失败: %s", err.Error()))
	}
	return &dm.DeviceRegisterResp{Psk: psk}, nil
}
