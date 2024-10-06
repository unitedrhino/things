package deviceauthlogic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/dgsvr/internal/svc"
	"gitee.com/i-Things/things/service/dgsvr/pb/dg"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type DeviceRegisterPayload struct {
	EncryptionType int    `json:"encryptionType"` //加密类型，1表示证书认证，2表示签名认证
	Psk            string `json:"psk"`            //设备密钥，当产品认证类型为签名认证时有此参数。
	ClientCert     string `json:"clientCert"`     //设备证书文件字符串格式，当产品认证类型为证书认证时有此参数
	ClientKey      string `json:"clientKey"`      //设备私钥文件字符串格式，当产品认证类型为证书认证时有此参数
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

func getPayload(encryptionType int, psk string, productSecret string) (size int, payload string, err error) {
	var data DeviceRegisterPayload
	data.EncryptionType = encryptionType
	data.Psk = psk
	pay, err := json.Marshal(data)
	if err != nil {
		return 0, "", err
	}
	payloadStr, err := utils.AesCbcBase64(string(pay), productSecret)
	if err != nil {
		return 0, "", err
	}

	return len(pay), payloadStr, nil
}

// 设备动态注册
func (l *DeviceRegisterLogic) DeviceRegister(in *dg.DeviceRegisterReq) (*dg.DeviceRegisterResp, error) {
	//检查产品动态开关是否开启
	pi, err := l.svcCtx.ProductM.ProductInfoRead(l.ctx, &dm.ProductInfoReadReq{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}
	if pi.AutoRegister == devices.DeviceRegisterUnable {
		return nil, errors.NotEnable.AddMsg("设备动态注册未启动")
	}
	//检查设备是否已注册
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			//检查设备自动创建是否开启， 开启则自动创建设备，未开启则返回错误
			if pi.AutoRegister == devices.DeviceAutoCreateEnable {
				//检查设备签名是否正确
				sig := getSignature(pi.Secret, fmt.Sprintf("deviceName=%s&nonce=%d&productID=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp))
				if sig != in.Signature {
					return nil, errors.Parameter.AddMsg("无效签名")
				}
				_, err := l.svcCtx.DeviceM.DeviceInfoCreate(l.ctx, &dm.DeviceInfo{
					ProductID:  in.ProductID,
					DeviceName: in.DeviceName,
				})
				if err != nil {
					return nil, errors.Database.AddMsgf("设备注册失败: %s", err.Error())
				}
				resp, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{ProductID: in.ProductID, DeviceName: in.DeviceName})
				if err != nil {
					return nil, errors.Database.AddMsgf("设备注册失败: %s", err.Error())
				}
				//将应答信息封装json 并加密
				length, payload, err := getPayload(devices.EncryptionTypeCert, resp.Secret, pi.Secret)
				return &dg.DeviceRegisterResp{Len: int64(length), Payload: payload}, nil
			}
			return nil, errors.NotFind.AddMsg("设备注册失败，无效设备")
		} else {
			return nil, errors.Database.AddMsgf("设备注册失败: %s", err.Error())
		}
	}

	if di.FirstLogin != 0 {
		return nil, errors.NotEmpty.AddMsg("设备已注册")
	}

	//检查设备签名是否正确
	sig := getSignature(pi.Secret, fmt.Sprintf("deviceName=%s&nonce=%d&productID=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp))
	if sig != in.Signature {
		return nil, errors.Parameter.AddMsg("无效签名")
	}

	//将应答信息封装json 并加密
	length, payload, err := getPayload(devices.EncryptionTypeCert, di.Secret, pi.Secret)
	return &dg.DeviceRegisterResp{Len: int64(length), Payload: payload}, nil
}
