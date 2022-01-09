package types

import (
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/usersvr/user"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

func UserCoreToApi(core *user.UserCore) *UserCore {
	return &UserCore{
		Uid:         core.Uid,
		UserName:    core.UserName,
		Email:       core.Email,
		Phone:       core.Phone,
		Wechat:      core.Wechat,
		LastIP:      core.LastIP,
		RegIP:       core.RegIP,
		CreatedTime: core.CreatedTime,
		Status:      core.Status,
	}
}

func UserInfoToApi(ui *user.UserInfo) *UserInfo {
	return &UserInfo{
		Uid:        ui.Uid,
		UserName:   ui.UserName,
		NickName:   ui.NickName,
		InviterUid: ui.InviterUid,
		InviterId:  ui.InviterId,
		Sex:        ui.Sex,
		City:       ui.City,
		Country:    ui.Country,
		Province:   ui.Province,
		Language:   ui.Language,
		HeadImgUrl: ui.HeadImgUrl,
		CreateTime: ui.CreateTime,
	}
}

func DeviceInfoToApi(v *dm.DeviceInfo) *DeviceInfo {
	return &DeviceInfo{
		ProductID:   v.ProductID,           //产品id 只读
		DeviceName:  v.DeviceName,          //设备名称 读写
		CreatedTime: v.CreatedTime,         //创建时间 只读
		Secret:      v.Secret,              //设备秘钥 只读
		FirstLogin:  v.FirstLogin,          //激活时间 只读
		LastLogin:   v.LastLogin,           //最后上线时间 只读
		Version:     GetNullVal(v.Version), // 固件版本  读写
		LogLevel:    v.LogLevel,            // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		Cert:        v.Cert,                // 设备证书  只读
	}
}

func ProductInfoToApi(v *dm.ProductInfo) *ProductInfo {
	return &ProductInfo{
		CreatedTime:  v.CreatedTime,             //创建时间 只读
		ProductID:    v.ProductID,               //产品id 只读
		ProductName:  v.ProductName,             //产品名称
		AuthMode:     v.AuthMode,                //认证方式:0:账密认证,1:秘钥认证
		DeviceType:   v.DeviceType,              //设备类型:0:设备,1:网关,2:子设备
		CategoryID:   v.CategoryID,              //产品品类
		NetType:      v.NetType,                 //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
		DataProto:    v.DataProto,               //数据协议:0:自定义,1:数据模板
		AutoRegister: v.AutoRegister,            //动态注册:0:关闭,1:打开,2:打开并自动创建设备
		Secret:       v.Secret,                  //动态注册产品秘钥 只读
		Description:  GetNullVal(v.Description), //描述
		DevStatus:    GetNullVal(v.DevStatus),   // 产品状态
	}
}

func ProductTemplateToApi(v *dm.ProductTemplate) *ProductTemplate {
	return &ProductTemplate{
		CreatedTime: v.CreatedTime,          //创建时间 只读
		ProductID:   v.ProductID,            //产品id 只读
		Template:    GetNullVal(v.Template), //数据模板
	}
}
