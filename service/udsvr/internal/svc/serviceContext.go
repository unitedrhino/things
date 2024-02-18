package svc

import (
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/dmdirect"
	"github.com/i-Things/things/service/udsvr/internal/config"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/zrpc"
)

type SvrClient struct {
	ProductM       productmanage.ProductManage
	DeviceInteract deviceinteract.DeviceInteract
	DeviceMsg      devicemsg.DeviceMsg
	DeviceM        devicemanage.DeviceManage

	TimedM   timedmanage.TimedManage
	AreaM    areamanage.AreaManage
	ProjectM projectmanage.ProjectManage
}

type ServiceContext struct {
	Config config.Config
	SvrClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM         timedmanage.TimedManage
		areaM          areamanage.AreaManage
		projectM       projectmanage.ProjectManage
		deviceM        devicemanage.DeviceManage
		productM       productmanage.ProductManage
		deviceInteract deviceinteract.DeviceInteract
		deviceMsg      devicemsg.DeviceMsg
	)
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database)
	timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
	projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
		deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
		deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
		deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
	}
	return &ServiceContext{
		Config: c,
		SvrClient: SvrClient{
			TimedM:         timedM,
			AreaM:          areaM,
			ProjectM:       projectM,
			ProductM:       productM,
			DeviceInteract: deviceInteract,
			DeviceMsg:      deviceMsg,
			DeviceM:        deviceM,
		},
	}
}
