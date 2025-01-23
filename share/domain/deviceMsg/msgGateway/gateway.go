package msgGateway

import (
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/schema"
)

type (
	//Msg 请求和回复结构体
	Msg struct {
		deviceMsg.CommonMsg
		Payload *GatewayPayload `json:"payload,omitempty"`
	}
	Devices []*Device
	Device  struct {
		ProductID    string `json:"productID"`              //产品id
		DeviceName   string `json:"deviceName"`             //设备id
		DeviceAlias  string `json:"deviceAlias"`            //设备名称
		DeviceSecret string `json:"deviceSecret,omitempty"` //设备秘钥
		Register
		Code int64  `json:"code,omitempty"` //子设备绑定结果
		Msg  string `json:"msg,omitempty"`  //错误原因
	}
	Register struct {
		/*
			子设备绑定签名串。 签名算法：
			1. 签名原串，将产品 GroupIDs 设备名称，随机数，时间戳拼接：text=${product_id};${device_name};${random};${expiration_time}
			2. 使用设备 Psk 密钥，或者证书的 Sha1 摘要，进行签名：hmac_sha1(device_secret, text)
		*/
		Signature  string `json:"signature,omitempty"`
		Random     int64  `json:"random,omitempty"`     //随机数。
		Timestamp  int64  `json:"timestamp,omitempty"`  //时间戳，单位：秒。
		SignMethod string `json:"signMethod,omitempty"` //签名算法。支持 hmacsha1、hmacsha256
	}
	GatewayPayload struct {
		Status      def.GatewayStatus   `json:"status,omitempty"`
		Devices     Devices             `json:"devices,omitempty"`
		Identifiers []string            `json:"identifiers,omitempty"` //内为希望设备上报的属性列表,不填为获取全部
		ProductID   string              `json:"productID,omitempty"`   //产品ID
		DeviceName  string              `json:"deviceName,omitempty"`
		Schema      *schema.ModelSimple `json:"schema,omitempty"` //物模型
	}
)

const (
	TypeTopo   = "topo"   //拓扑关系管理
	TypeStatus = "status" //代理子设备上下线
	TypeThing  = "thing"  //物模型操作
)

// 获取产品id列表(不重复的)
func (d Devices) GetProductIDs() []string {
	var (
		set = map[string]struct{}{}
	)
	for _, v := range d {
		set[v.ProductID] = struct{}{}
	}
	return utils.SetToSlice(set)
}
func (d Devices) GetCore() Devices {
	if d == nil {
		return nil
	}
	var ret Devices
	for _, v := range d {
		ret = append(ret, &Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return ret
}
func (d Devices) GetDevCore() []*devices.Core {
	if d == nil {
		return nil
	}
	var ret []*devices.Core
	for _, v := range d {
		ret = append(ret, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return ret
}
