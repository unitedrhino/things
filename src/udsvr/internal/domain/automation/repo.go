package automation

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	deviceinteract "github.com/i-Things/things/src/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	devicemsg "github.com/i-Things/things/src/dmsvr/client/devicemsg"
)
import "context"

type InfoFilter struct {
	Name        string `json:"name"`
	Status      int64
	TriggerType TriggerType
	AlarmID     int64 //绑定的告警id
}

type Repo interface {
	Insert(ctx context.Context, info *Info) (id int64, err error)
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
	Alarm          AlarmRepo
	Device         devices.Core
	Serial         Serial
	Scene          *Info
}

type AlarmRepo interface {
	//告警触发
	AlarmTrigger(ctx context.Context, in TriggerSerial) error
	//告警解除
	AlarmRelieve(ctx context.Context, in AlarmRelieve) error
}
type Serial interface {
	GenSerial() string
}
