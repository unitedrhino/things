package scene

import (
	"gitee.com/i-Things/core/service/syssvr/client/notifymanage"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"github.com/i-Things/things/service/dmsvr/client/devicegroup"
	deviceinteract "github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/client/devicemanage"
	devicemsg "github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/i-Things/things/service/dmsvr/dmExport"
)
import "context"

type InfoFilter struct {
	Name        string `json:"name"`
	Status      int64
	TriggerType SceneType
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
	GetInfos(ctx context.Context, device devices.Core, operator TriggerDeviceType, dataID string) (Infos, error)
}

type ValidateRepo struct {
	Ctx                context.Context
	DeviceCache        dmExport.DeviceCacheT
	ProductCache       dmExport.ProductCacheT
	ProductSchemaCache dmExport.SchemaCacheT
	GetSceneInfo       func(ctx context.Context, sceneID int64) (*Info, error)
	Info               *Info
}

type WhenRepo interface {
}

// TermRepo 场景运行需要用到的对外仓储
type TermRepo struct {
	DeviceMsg  devicemsg.DeviceMsg
	SchemaRepo schema.ReadRepo
}

type ActionRepo struct {
	Info           *Info
	UserID         int64
	DeviceInteract deviceinteract.DeviceInteract
	DeviceM        devicemanage.DeviceManage
	DeviceG        devicegroup.DeviceGroup
	NotifyM        notifymanage.NotifyManage
	SceneExec      func(ctx context.Context, sceneID int64) error
	AlarmExec      func(ctx context.Context, in AlarmSerial) error
}

//type AlarmRepo interface {
//	//告警触发
//	AlarmTrigger(ctx context.Context, in AlarmSerial) error
//	//告警解除
//	AlarmRelieve(ctx context.Context, in AlarmRelieve) error
//}
//type Serial interface {
//	GenSerial() string
//}
