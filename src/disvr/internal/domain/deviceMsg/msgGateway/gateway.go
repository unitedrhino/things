package msgGateway

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
)

type (
	//Msg 请求和回复结构体
	Msg struct {
		*deviceMsg.CommonMsg
		Payload *GateWayData `json:"payload,omitempty"`
	}
	Device struct {
		ProductID  string `json:"productID"`        //产品id
		DeviceName string `json:"deviceName"`       //设备名称
		Result     int64  `json:"result,omitempty"` //子设备绑定结果
	}
	GateWayData struct {
		Status  int       `json:"status,omitempty"`
		Devices []*Device `json:"devices"`
	}
)

const (
	TypeOperation = "operation" //拓扑关系管理
	TypeStatus    = "status"    //代理子设备上下线
)
