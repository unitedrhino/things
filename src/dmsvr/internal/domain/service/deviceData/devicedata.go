// Package repo 本文件是提供设备模型数据存储的信息
package deviceData

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"time"
)

type (
	// PropertyData 属性数据
	PropertyData struct {
		ID        string    `json:"id"`         //属性的id
		Param     any       `json:"property" `  //一个属性的参数
		TimeStamp time.Time `json:"timeStamp" ` //时间戳
	}
	// EventData 事件数据
	EventData struct {
		ID        string         `json:"id" `        //事件id
		Type      string         `json:"type" `      //事件类型: 信息:info  告警alert  故障:fault
		Params    map[string]any `json:"params" `    //事件参数
		TimeStamp time.Time      `json:"timeStamp" ` //时间戳
	}

	FilterOpt struct {
		Page       def.PageInfo2
		ProductID  string
		DeviceName string
		DataID     string
		Order      string //时间排序 aes(默认,从久到近排序) desc(时间从近到久排序)
		Interval   int64  //间隔(单位毫秒) 如果这个值不为零值 则时间的开始和结束必须有效及聚合函数不应该为空
		ArgFunc    string //聚合函数 avg:平均值 first:第一个参数 last:最后一个参数 count:总数 twa: 时间加权平均函数 参考:https://docs.taosdata.com/taos-sql/function
	}

	DeviceDataRepo interface {
		// InsertEventData 插入事件数据
		InsertEventData(ctx context.Context, productID string, deviceName string, event *EventData) error
		// InsertPropertyData 插入一条属性数据
		InsertPropertyData(ctx context.Context, t *schema.Model, productID string, deviceName string, property *PropertyData) error
		// InsertPropertiesData 插入多条属性数据 params key为属性的id,val为属性的值
		InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string, params map[string]any, timestamp time.Time) error
		// GetEventDataWithID 根据事件id获取事件信息
		GetEventDataByID(ctx context.Context, filter FilterOpt) ([]*EventData, error)
		GetEventCountByID(ctx context.Context, filter FilterOpt) (int64, error)
		// GetPropertyDataByID 根据属性id获取属性信息
		GetPropertyDataByID(ctx context.Context, filter FilterOpt) ([]*PropertyData, error)
		GetPropertyCountByID(ctx context.Context, filter FilterOpt) (int64, error)
		// InitProduct 初始化产品的物模型相关表及日志记录表
		InitProduct(ctx context.Context, t *schema.Model, productID string) error
		// DropProduct 删除产品时需要删除产品下的所有表
		DropProduct(ctx context.Context, t *schema.Model, productID string) error
		// InitDevice 创建设备时为设备创建单独的表
		InitDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		// DropDevice 删除设备时需要删除设备的所有表
		DropDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		// ModifyProduct 修改产品物模型 只支持新增和删除,不支持修改数据类型
		ModifyProduct(ctx context.Context, oldT *schema.Model, newt *schema.Model, productID string) error
	}
)

func (f FilterOpt) Check() error {
	if f.Interval != 0 && f.ArgFunc == "" {
		return errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	return nil
}
