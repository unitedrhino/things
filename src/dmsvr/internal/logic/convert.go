package logic

import (
	"database/sql"
	"fmt"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func DBToRPCFmt(db interface{}) interface{} {
	switch db.(type) {
	case *model.DeviceInfo:
		di := db.(*model.DeviceInfo)
		return &dm.DeviceInfo{
			Version:     &wrappers.StringValue{Value: di.Version},
			LogLevel:    di.LogLevel,
			Cert:        di.Cert,
			ProductID:   di.ProductID,
			DeviceName:  di.DeviceName,
			CreatedTime: di.CreatedTime.Unix(),
			FirstLogin:  GetNullTime(di.FirstLogin),
			LastLogin:   GetNullTime(di.LastLogin),
			Secret:      di.Secret,
		}
	case *model.ProductInfo:
		pi := db.(*model.ProductInfo)
		return &dm.ProductInfo{
			ProductID:    pi.ProductID,                                 //产品id
			ProductName:  pi.ProductName,                               //产品名
			AuthMode:     pi.AuthMode,                                  //认证方式:0:账密认证,1:秘钥认证
			DeviceType:   pi.DeviceType,                                //设备类型:0:设备,1:网关,2:子设备
			CategoryID:   pi.CategoryID,                                //产品品类
			NetType:      pi.NetType,                                   //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
			DataProto:    pi.DataProto,                                 //数据协议:0:自定义,1:数据模板
			AutoRegister: pi.AutoRegister,                              //动态注册:0:关闭,1:打开,2:打开并自动创建设备
			Secret:       pi.Secret,                                    //动态注册产品秘钥 只读
			Description:  &wrappers.StringValue{Value: pi.Description}, //描述
			CreatedTime:  pi.CreatedTime.Unix(),                        //创建时间
		}
	default:
		panic(fmt.Sprintf("ToRPCFmt not suppot:%#v", db))
	}
}
