package deviceauthlogic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/dgsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dgsvr/pb/dg"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type DeviceRegisterPayload struct {
	EncryptionType int    `json:"encryptionType"` //加密类型，1表示证书认证，2表示签名认证aes-128-cbc
	Psk            string `json:"psk"`            //设备密钥，当产品认证类型为签名认证时有此参数。
	ClientCert     string `json:"clientCert"`     //设备证书文件字符串格式，当产品认证类型为证书认证时有此参数
	ClientKey      string `json:"clientKey"`      //设备私钥文件字符串格式，当产品认证类型为证书认证时有此参数
}

// DeviceRegisterPlainPayload 是 retEnc=hex 时返回给轻量 MCU 的明文密钥载荷。
type DeviceRegisterPlainPayload struct {
	EncryptionType int    `json:"encryptionType"` // 加密类型，2 表示签名认证密钥。
	Psk            string `json:"psk"`            // Base64 设备密钥解码后的十六进制字符串。
}

func NewDeviceRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceRegisterLogic {
	return &DeviceRegisterLogic{
		ctx:    ctxs.WithRoot(ctx),
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// getSignatureHex 使用 HMAC-sha1 生成旧协议 Base64 编码前的 hex 签名。
func getSignatureHex(secret string, dest string) string {
	if secret == "" || dest == "" {
		return ""
	}

	return utils.HmacSha1(dest, []byte(secret))
}

// getSignature 保留旧协议签名格式：HMAC-sha1 的 hex 字符串再做 Base64。
func getSignature(secret string, dest string) string {
	hexSign := getSignatureHex(secret, dest)
	if hexSign == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(hexSign))
}

// checkSignature 同时兼容旧 Base64 签名和新 MCU 直接 hex 签名。
func checkSignature(secret string, dest string, signature string) bool {
	if signature == "" {
		return false
	}
	hexSign := getSignatureHex(secret, dest)
	if hexSign == "" {
		return false
	}
	return signature == base64.StdEncoding.EncodeToString([]byte(hexSign)) ||
		strings.EqualFold(signature, hexSign)
}

// getPayload 按 retEnc 生成设备动态注册响应，hex 模式返回明文十六进制 psk。
func getPayload(encryptionType int, retEnc string, psk string, productSecret string) (size int, payload string, err error) {
	if strings.EqualFold(retEnc, "hex") {
		rawPsk, err := base64.StdEncoding.DecodeString(psk)
		if err != nil {
			return 0, "", errors.Parameter.AddMsg("设备密钥Base64解码失败")
		}
		data := DeviceRegisterPlainPayload{
			EncryptionType: devices.EncTypeKey,
			Psk:            hex.EncodeToString(rawPsk),
		}
		pay, err := json.Marshal(data)
		if err != nil {
			return 0, "", err
		}
		return len(pay), string(pay), nil
	}

	var data DeviceRegisterPayload
	data.EncryptionType = encryptionType
	data.Psk = psk
	pay, err := json.Marshal(data)
	if err != nil {
		return 0, "", err
	}
	if retEnc == "aes128ecb" {
		payloadStr, err := utils.AesEcbBase64(string(pay), productSecret)
		if err != nil {
			return 0, "", err
		}
		return len(pay), payloadStr, nil
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
				signSource := fmt.Sprintf("deviceName=%s&nonce=%d&productID=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp)
				if !checkSignature(pi.Secret, signSource, in.Signature) {
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
				length, payload, err := getPayload(devices.EncTypeCert, in.RetEnc, resp.Secret, pi.Secret)
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
	signSource := fmt.Sprintf("deviceName=%s&nonce=%d&productID=%s&timestamp=%d", in.DeviceName, in.Nonce, in.ProductID, in.Timestamp)
	if !checkSignature(pi.Secret, signSource, in.Signature) {
		return nil, errors.Parameter.AddMsg("无效签名")
	}

	//将应答信息封装json 并加密
	length, payload, err := getPayload(devices.EncTypeCert, in.RetEnc, di.Secret, pi.Secret)
	return &dg.DeviceRegisterResp{Len: int64(length), Payload: payload}, nil
}
