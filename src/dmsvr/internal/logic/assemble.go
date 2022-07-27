package logic

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func ToProductSchema(pt *schema.SchemaInfo) *dm.ProductSchema {
	return &dm.ProductSchema{
		CreatedTime: pt.CreatedTime.Unix(),
		ProductID:   pt.ProductID,
		Schema:      pt.Template,
	}
}

func ToFirmwareInfo(fi *mysql.ProductFirmware) *dm.FirmwareInfo {
	return &dm.FirmwareInfo{
		ProductID:   fi.ProductID,
		Version:     fi.Version,
		Name:        fi.Name,
		Description: fi.Description,
		Dir:         fi.Dir,
		Size:        fi.Size,
	}
}

func ToDeviceInfo(di *mysql.DeviceInfo) *dm.DeviceInfo {
	var (
		tags map[string]string
	)

	if di.Tags != "" {
		_ = json.Unmarshal([]byte(di.Tags), &tags)
	}
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
		IsOnline:    cast.ToBool(di.IsOnline),
		Tags:        tags,
	}
}

func ToProductInfo(pi *mysql.ProductInfo) *dm.ProductInfo {
	dpi := &dm.ProductInfo{
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
		//Model:     &wrappers.StringValue{Value: pi.Model},    //数据模板
	}
	return dpi
}

func ToDataHubLogIndex(log *device.HubLog) *dm.DataHubLogIndex {
	return &dm.DataHubLogIndex{
		Timestamp:  log.Timestamp.UnixMilli(),
		Action:     log.Action,
		RequestID:  log.RequestID,
		TranceID:   log.TranceID,
		Topic:      log.Topic,
		Content:    log.Content,
		ResultType: log.ResultType,
	}
}

//SDK调试日志
func ToDataSdkLogIndex(log *device.SDKLog) *dm.DataSdkLogIndex {
	return &dm.DataSdkLogIndex{
		Timestamp: log.Timestamp.UnixMilli(),
		Loglevel:  log.LogLevel,
		Content:   log.Content,
	}
}
