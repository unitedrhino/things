package svc

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/domain/schema"
	"gitee.com/i-Things/core/shared/eventBus"
	"gitee.com/i-Things/core/shared/stores"
	deviceinteract "github.com/i-Things/things/src/dmsvr/client/deviceinteract"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	devicemsg "github.com/i-Things/things/src/dmsvr/client/devicemsg"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/cache"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/dataUpdate"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/timer"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config config.Config
	Repo
	SvrClient
	SceneTimerControl timer.SceneControl
	Bus               eventBus.Bus
	DataUpdate        dataUpdate.DataUpdate
}
type Repo struct {
	Store           kv.Store
	SceneDeviceRepo scene.DeviceRepo
	SchemaRepo      schema.ReadRepo
}
type SvrClient struct {
	ProductM       productmanage.ProductManage
	DeviceInteract deviceinteract.DeviceInteract
	DeviceMsg      devicemsg.DeviceMsg
	DeviceM        devicemanage.DeviceManage
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		deviceM        devicemanage.DeviceManage
		productM       productmanage.ProductManage
		deviceInteract deviceinteract.DeviceInteract
		deviceMsg      devicemsg.DeviceMsg
	)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("rulesvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}

	store := kv.NewStore(c.CacheRedis)
	sceneDevice := cache.NewSceneDeviceRepo(relationDB.NewSceneInfoRepo(context.TODO()))
	err = sceneDevice.Init(context.TODO())
	if err != nil {
		logx.Error("设备场景数据初始化失败 err:", err)
		os.Exit(-1)
	}
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

	tr := schema.NewReadRepo(func(ctx context.Context, productID string) (*schema.Model, error) {
		info, err := productM.ProductSchemaTslRead(ctx, &dm.ProductSchemaTslReadReq{ProductID: productID})
		if err != nil {
			return nil, err
		}
		return schema.ValidateWithFmt([]byte(info.Tsl))
	})

	bus := eventBus.NewEventBus()
	du, err := dataUpdate.NewDataUpdate(c.Event)
	logx.Must(err)
	return &ServiceContext{
		Bus:        bus,
		Config:     c,
		DataUpdate: du,
		SvrClient: SvrClient{
			ProductM:       productM,
			DeviceInteract: deviceInteract,
			DeviceMsg:      deviceMsg,
			DeviceM:        deviceM,
		},
		Repo: Repo{
			Store:           store,
			SceneDeviceRepo: sceneDevice,
			SchemaRepo:      tr,
		},
	}
}
