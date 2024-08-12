package scene

import (
	"gitee.com/i-Things/core/service/syssvr/client/common"
	"gitee.com/i-Things/core/service/syssvr/client/notifymanage"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/dmsvr/client/devicegroup"
	deviceinteract "github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/client/devicemanage"
	devicemsg "github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/i-Things/things/service/dmsvr/dmExport"
)
import "context"

type CheckRepo struct {
	Ctx                context.Context
	DeviceCache        dmExport.DeviceCacheT
	UserShareCache     dmExport.UserShareCacheT
	ProductCache       dmExport.ProductCacheT
	ProductSchemaCache dmExport.SchemaCacheT
	ProjectCache       sysExport.ProjectCacheT
	DeviceMsg          devicemsg.DeviceMsg
	Common             common.Common
	GetSceneInfo       func(ctx context.Context, sceneID int64) (*Info, error)
	Info               *Info
}

type ActionRepo struct {
	Info           *Info
	UserID         int64
	DeviceInteract deviceinteract.DeviceInteract
	DeviceM        devicemanage.DeviceManage
	DeviceCache    dmExport.DeviceCacheT
	ProductCache   dmExport.ProductCacheT
	DeviceG        devicegroup.DeviceGroup
	NotifyM        notifymanage.NotifyManage
	SceneExec      func(ctx context.Context, sceneID int64, status def.Bool) error
	AlarmExec      func(ctx context.Context, in AlarmSerial) error
	SaveLog        func(ctx context.Context, log *Log) error
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
