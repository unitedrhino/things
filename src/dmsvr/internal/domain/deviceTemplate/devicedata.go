// Package repo 本文件是提供设备模型数据存储的信息
package deviceTemplate

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"time"
)

type (
	//返回给前端的属性
	PropertyData struct {
		ID        string      `json:"id"`         //属性的id
		Param     interface{} `json:"property" `  //一个属性的参数
		TimeStamp time.Time   `json:"timeStamp" ` //时间戳
	}
	EventData struct {
		ID        string                 `json:"id" `        //事件id
		Type      string                 `json:"type" `      //事件类型: 信息:info  告警alert  故障:fault
		Params    map[string]interface{} `json:"params" `    //事件参数
		TimeStamp time.Time              `json:"timeStamp" ` //时间戳
	}

	DeviceData2Repo interface {
		// √
		InsertEventData(ctx context.Context, t *Template, productID string, deviceName string, event *EventData) error
		InsertPropertyData(ctx context.Context, t *Template, productID string, deviceName string, property *PropertyData) error
		//params key为属性的id,val为属性的值
		InsertPropertiesData(ctx context.Context, t *Template, productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error
		GetEventDataWithID(ctx context.Context, t *Template, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*EventData, error)
		GetPropertyDataWithID(ctx context.Context, t *Template, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*PropertyData, error)
		// InitProduct 初始化产品的物模型相关表及日志记录表 √
		InitProduct(ctx context.Context, t *Template, productID string) error
		// DropProduct 删除产品时需要删除产品下的所有表
		DropProduct(ctx context.Context, t *Template, productID string) error
		// InitDevice 创建设备时为设备创建单独的表 √
		InitDevice(ctx context.Context, t *Template, productID string, deviceName string) error
		// DropDevice 删除设备时需要删除设备的所有表
		DropDevice(ctx context.Context, t *Template, productID string, deviceName string) error
		// ModifyProduct 修改产品物模型
		ModifyProduct(ctx context.Context, oldT *Template, newt *Template, productID string) error
	}

	DeviceDataRepo interface {
		InsertEventData(productID string, deviceName string, event *EventData) error
		InsertPropertyData(productID string, deviceName string, property *PropertyData) error
		//params key为属性的id,val为属性的值
		InsertPropertiesData(productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error
		GetEventDataWithID(productID string, deviceName string, dataID string, page def.PageInfo2) ([]*EventData, error)
		GetPropertyDataWithID(productID string, deviceName string, dataID string, page def.PageInfo2) ([]*PropertyData, error)
		CreatePropertyDB(productID string) error
		CreateEventDB(productID string) error
		CreateLogDB(productID string) error
	}
	GetDeviceDataRepo func(ctx context.Context) DeviceDataRepo
)
