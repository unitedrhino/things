// Package repo 本文件是提供设备模型数据存储的信息
package msgThing

import (
	"context"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/jsonx"
	"time"
)

type (
	// PropertyData 属性数据
	PropertyData struct {
		TenantCode  dataType.TenantCode    `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                        // 租户编码
		ProjectID   dataType.ProjectID     `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
		AreaID      dataType.AreaID        `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`    // 项目区域ID(雪花ID)
		AreaIDPath  dataType.AreaIDPath    `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`                 // 项目区域ID路径(雪花ID)
		BelongGroup map[string]def.IDsInfo `gorm:"column:belong_group;type:json;serializer:json;default:'{}'"`

		DeviceName string    `gorm:"column:device_name;type:varchar(50);NOT NULL" json:"deviceName"`
		Identifier string    `gorm:"column:identifier;type:varchar(50);NOT NULL" json:"identifier"` //标识符
		Param      any       `gorm:"column:param;type:varchar(256);NOT NULL" json:"param" `         //一个属性的参数
		TimeStamp  time.Time `gorm:"column:ts;NOT NULL;" json:"timeStamp"`                          //时间戳
	}
	// EventData 事件数据
	EventData struct {
		TenantCode  dataType.TenantCode    `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                        // 租户编码
		ProjectID   dataType.ProjectID     `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
		AreaID      dataType.AreaID        `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`    // 项目区域ID(雪花ID)
		AreaIDPath  dataType.AreaIDPath    `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`                 // 项目区域ID路径(雪花ID)
		BelongGroup map[string]def.IDsInfo `gorm:"column:belong_group;type:json;serializer:json;default:'{}'"`

		Identifier string         `gorm:"column:identifier;type:varchar(50);NOT NULL" json:"identifier"` //标识符
		Type       string         `gorm:"column:type;type:varchar(20);NOT NULL" json:"type" `            //事件类型: 信息:info  告警alert  故障:fault
		Params     map[string]any `gorm:"column:param;type:varchar(256);NOT NULL" json:"params" `        //事件参数
		TimeStamp  time.Time      `gorm:"column:ts;NOT NULL;" json:"timeStamp"`                          //时间戳
		DeviceName string         `gorm:"column:device_name;type:varchar(50);NOT NULL" json:"device_name" `
	}
	/*
	   FILL 语句指定某一窗口区间数据缺失的情况下的填充模式。填充模式包括以下几种：
	   不进行填充：NONE（默认填充模式）。
	   VALUE 填充：固定值填充，此时需要指定填充的数值。例如：FILL(VALUE, 1.23)。这里需要注意，最终填充的值受由相应列的类型决定，如 FILL(VALUE, 1.23)，相应列为 INT 类型，则填充值为 1。
	   PREV 填充：使用前一个非 NULL 值填充数据。例如：FILL(PREV)。
	   NULL 填充：使用 NULL 填充数据。例如：FILL(NULL)。
	   LINEAR 填充：根据前后距离最近的非 NULL 值做线性插值填充。例如：FILL(LINEAR)。
	   NEXT 填充：使用下一个非 NULL 值填充数据。例如：FILL(NEXT)。
	*/
	FilterOpt struct {
		Page def.PageInfo2

		TenantCode  string
		ProjectID   int64  `json:"projectID,omitempty"`
		AreaID      int64  `json:"areaID,omitempty"`
		AreaIDPath  string `json:"areaIDPath,omitempty"`
		BelongGroup map[string]def.IDsInfo
		AreaIDs     []int64 `json:"areaIDs"`

		ProductID  string
		ProductIDs []string
		//DeviceName  string
		DeviceNames  []string
		DataID       string
		Types        []string     //事件类型: 信息:info  告警alert  故障:fault
		Order        stores.Order //0:aes(默认,从久到近排序) 1:desc(时间从近到久排序)
		Interval     int64        //间隔(单位毫秒) 如果这个值不为零值 则时间的开始和结束必须有效及聚合函数不应该为空
		IntervalUnit def.TimeUnit //间隔单位 a (毫秒,默认), d (天), h (小时), m (分钟), n (月), s (秒), u (微秒), w (周), y (年)
		Fill         string       //指定窗口区间数据缺失的情况下的填充模式
		ArgFunc      string       //聚合函数 avg:平均值 first:第一个参数 last:最后一个参数 count:总数 twa: 时间加权平均函数 参考:https://docs.taosdata.com/taos-sql/function
		PartitionBy  string       //切分数据,可以填写deviceName
	}
	LatestFilter struct {
		ProductID  string
		DeviceName string
		DataID     string
	}
	Optional struct {
		Sync      bool //同步执行
		OnlyCache bool //只记录到缓存中

		TenantCode  dataType.TenantCode // 租户编码
		ProjectID   dataType.ProjectID  // 项目ID(雪花ID)
		AreaID      dataType.AreaID     // 项目区域ID(雪花ID)
		AreaIDPath  dataType.AreaIDPath // 项目区域ID路径(雪花ID)
		BelongGroup map[string]def.IDsInfo
	}

	SchemaDataRepo interface {
		Init(ctx context.Context) error
		VersionUpdate(ctx context.Context, version string, dc *caches.Cache[dm.DeviceInfo, devices.Core]) error
		// InsertEventData 插入事件数据
		InsertEventData(ctx context.Context, productID string, deviceName string, event *EventData) error
		// InsertPropertyData 插入一条属性数据
		InsertPropertyData(ctx context.Context, t *schema.Property, productID string, deviceName string, property *Param, timestamp time.Time, optional Optional) error
		// InsertPropertiesData 插入多条属性数据 params key为属性的id,val为属性的值
		InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string, params map[string]Param, timestamp time.Time, optional Optional) error
		// GetEventDataWithID 根据事件id获取事件信息
		GetEventDataByFilter(ctx context.Context, filter FilterOpt) ([]*EventData, error)
		GetEventCountByFilter(ctx context.Context, filter FilterOpt) (int64, error)
		// GetPropertyDataByID 根据属性id获取属性信息
		GetPropertyDataByID(ctx context.Context, p *schema.Property, filter FilterOpt) ([]*PropertyData, error)
		GetLatestPropertyDataByID(ctx context.Context, p *schema.Property, filter LatestFilter) (*PropertyData, error)
		GetPropertyCountByID(ctx context.Context, p *schema.Property, filter FilterOpt) (int64, error)
		// InitProduct 初始化产品的物模型相关表及日志记录表
		InitProduct(ctx context.Context, t *schema.Model, productID string) error
		// DeleteProduct 删除产品时需要删除产品下的所有表
		DeleteProduct(ctx context.Context, t *schema.Model, productID string) error
		// InitDevice 创建设备时为设备创建单独的表
		InitDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		// DeleteDevice 删除设备时需要删除设备的所有表
		DeleteDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		DeleteDeviceProperty(ctx context.Context, productID string, deviceName string, s []schema.Property) error
		// UpdateProduct 修改产品物模型 只支持新增和删除,不支持修改数据类型
		//UpdateProduct(ctx context.Context, oldT *schema.Model, newt *schema.Model, productID string) error
		CreateProperty(ctx context.Context, p *schema.Property, productID string) error
		DeleteProperty(ctx context.Context, p *schema.Property, productID string, identifier string) error
		UpdateProperty(ctx context.Context, oldP *schema.Property, newP *schema.Property, productID string) error

		UpdateDevice(ctx context.Context, dev devices.Core, t *schema.Model, affiliation devices.Affiliation) error
	}
)

func (p *PropertyData) String() string {

	v, _ := jsonx.Marshal(p)
	return string(v)
}

func (p *PropertyData) Fmt() *PropertyData {
	switch param := p.Param.(type) {
	case map[string]any:
		for k, v := range param {
			param[k] = utils.BoolToInt(v)
		}
		p.Param = param
	default:
		p.Param = utils.BoolToInt(p.Param)
	}
	return p
}

func (f FilterOpt) Check() error {
	if f.Interval != 0 && f.ArgFunc == "" {
		return errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	return nil
}
