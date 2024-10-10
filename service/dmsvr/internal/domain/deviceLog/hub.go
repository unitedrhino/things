package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/domain/application"
	"time"
)

type ActionType = string

const (
	ActionTypeGateway  ActionType = "gateway"  //网关操作子设备
	ActionTypeOta      ActionType = "ota"      //ota升级消息
	ActionTypeProperty ActionType = "property" //物模型属性消息
	ActionTypeEvent    ActionType = "event"    //事件消息
	ActionTypeAction   ActionType = "action"   //行为消息
	ActionTypeExt      ActionType = "ext"      //拓展消息
	ActionTypeNtp      ActionType = "ntp"      //获取时间
)

type (
	Hub struct {
		ProductID   string     `json:"productID"`             // 产品id
		DeviceName  string     `json:"deviceName"`            // 设备名称
		Content     string     `json:"content,omitempty"`     // 具体信息
		Topic       string     `json:"topic,omitempty"`       // 主题
		Action      ActionType `json:"action,omitempty"`      // 操作类型
		Timestamp   time.Time  `json:"timestamp,omitempty"`   // 操作时间
		RequestID   string     `json:"requestID,omitempty"`   // 请求ID
		TraceID     string     `json:"traceID,omitempty"`     // 服务器端事务id
		ResultCode  int64      `json:"resultCode,omitempty"`  // 请求结果状态,200为成功
		RespPayload string     `json:"respPayload,omitempty"` //返回的内容
	}
	HubFilter struct {
		ProductID  string   // 产品id
		DeviceName string   // 设备名称
		Actions    []string //过滤操作类型 property:属性 event:事件 action:操作 thing:物模型提交的操作为匹配的日志
		Topics     []string //过滤主题
		Content    string   //过滤内容
		RequestID  string   //过滤请求ID
	}
	HubRepo interface {
		ManageRepo
		GetDeviceLog(ctx context.Context, filter HubFilter, page def.PageInfo2) ([]*Hub, error)
		GetCountLog(ctx context.Context, filter HubFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Hub) error
	}
)

func (h Hub) ToApp() application.Hub {
	return application.Hub{
		ProductID:   h.ProductID,
		DeviceName:  h.DeviceName,
		Content:     h.Content,
		Topic:       h.Topic,
		Action:      h.Action,
		Timestamp:   h.Timestamp.UnixMilli(),
		RequestID:   h.RequestID,
		TraceID:     h.TraceID,
		ResultCode:  h.ResultCode,
		RespPayload: h.RespPayload,
	}
}
