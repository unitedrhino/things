package scene

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	deviceinteract "github.com/i-Things/things/src/disvr/client/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/client/devicemsg"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
)
import "context"

type TriggerDeviceFilter struct {
	ProductID string `json:"productID"` //产品id
	//Operator  OperationSchema `json:"operator"`  //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
}

type InfoFilter struct {
	Name          string `json:"name"`
	State         int64
	TriggerType   TriggerType
	TriggerDevice *TriggerDeviceFilter //只有设备类型可以查询
}

type Repo interface {
	Insert(ctx context.Context, info *Info) error
	Update(ctx context.Context, info *Info) error
	Delete(ctx context.Context, id int64) error
	FindOne(ctx context.Context, id int64) (*Info, error)
	FindOneByName(ctx context.Context, name string) (*Info, error)
	FindByFilter(ctx context.Context, filter InfoFilter, page *def.PageInfo) (Infos, error)
	CountByFilter(ctx context.Context, filter InfoFilter) (size int64, err error)
}

type DeviceRepo interface {
	Insert(ctx context.Context, info *Info) error
	Update(ctx context.Context, info *Info) error
	Delete(ctx context.Context, id int64) error
	GetInfos(ctx context.Context, device devices.Core, operator DeviceOperationOperator, dataID string) (Infos, error)
}

// TermRepo 场景运行需要用到的对外仓储
type TermRepo struct {
	DeviceMsg  devicemsg.DeviceMsg
	SchemaRepo schema.ReadRepo
}

type ActionRepo struct {
	DeviceInteract deviceinteract.DeviceInteract
	DeviceM        devicemanage.DeviceManage
}
