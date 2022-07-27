package types

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/usersvr/user"
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
		Role:        core.Role,
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

//func UserInfoFullToApi(ui *user.UserInfo) *types.UserIndexResp {
//	return &types.UserIndexResp{
//		Uid:        ui.Uid,
//		UserName:   ui.UserName,
//		NickName:   ui.NickName,
//		InviterUid: ui.InviterUid,
//		InviterId:  ui.InviterId,
//		Sex:        ui.Sex,
//		City:       ui.City,
//		Country:    ui.Country,
//		Province:   ui.Province,
//		Language:   ui.Language,
//		HeadImgUrl: ui.HeadImgUrl,
//		CreateTime: ui.CreateTime,
//	}
//}
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
	}
}

func ProductSchemaToApi(v *dm.ProductSchema) ProductSchema {
	return ProductSchema{
		CreatedTime: v.CreatedTime, //创建时间 只读
		ProductID:   v.ProductID,   //产品id 只读
		Schema:      v.Schema,      //数据模板
	}
}
