//本文件是提供设备模型数据存储的信息
package repo

import (
	"context"
	"time"
)

type (
	//返回给前端的属性
	Property struct {
		ID        string      `json:"id"`         //属性的id
		Param     interface{} `json:"property" `  //一个属性的参数
		TimeStamp time.Time   `json:"timeStamp" ` //时间戳
	}
	Event struct {
		ID        string                 `json:"id" `        //事件id
		Type      string                 `json:"type" `      //事件类型: 信息:info  告警alert  故障:fault
		Params    map[string]interface{} `json:"params" `    //事件参数
		TimeStamp time.Time              `json:"timeStamp" ` //时间戳
	}

	DeviceDataRepo interface {
		InsertEventData(productID string, deviceName string, event *Event) error
		InsertPropertyData(productID string, deviceName string, property *Property) error
		//params key为属性的id,val为属性的值
		InsertPropertiesData(productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error
		GetEventDataWithID(productID string, deviceName string, dataID string, timeStart, timeEnd int64, limit int64) ([]*Event, error)
		GetPropertyDataWithID(productID string, deviceName string, dataID string, timeStart, timeEnd int64, limit int64) ([]*Property, error)
		CreatePropertyDB(productID string) error
		CreateEventDB(productID string) error
		CreateLogDB(productID string) error
	}
	GetDeviceDataRepo func(ctx context.Context) DeviceDataRepo
)
