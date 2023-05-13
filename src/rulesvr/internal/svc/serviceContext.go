package svc

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/schema"
	deviceinteract "github.com/i-Things/things/src/disvr/client/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/client/devicemsg"
	"github.com/i-Things/things/src/disvr/didirect"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/cache"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/internal/timer"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config config.Config
	Repo
	SvrClient
	SceneTimerControl timer.SceneControl
}
type Repo struct {
	Store               kv.Store
	SceneRepo           scene.Repo
	SceneDeviceRepo     scene.DeviceRepo
	SchemaRepo          schema.ReadRepo
	SceneInfoRepo       mysql.RuleSceneInfoModel
	AlarmInfoRepo       mysql.RuleAlarmInfoModel
	AlarmRecordRepo     mysql.RuleAlarmRecordModel
	AlarmSceneRepo      mysql.RuleAlarmSceneModel
	AlarmDealRecordRepo mysql.RuleAlarmDealRecordModel
	AlarmLogRepo        mysql.RuleAlarmLogModel
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
	store := kv.NewStore(c.CacheRedis)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	SceneInfoRepo := mysql.NewRuleSceneInfoModel(conn)
	AlarmInfoRepo := mysql.NewRuleAlarmInfoModel(conn)
	AlarmSceneRepo := mysql.NewRuleAlarmSceneModel(conn)
	AlarmDealRecordRepo := mysql.NewRuleAlarmDealRecordModel(conn)
	AlarmLogRepo := mysql.NewRuleAlarmLogModel(conn)
	AlarmRecordRepo := mysql.NewRuleAlarmRecordModel(conn)
	SceneRepo := mysql.NewRuleSceneInfoModel(conn)
	sceneDevice := cache.NewSceneDeviceRepo(SceneRepo)
	err := sceneDevice.Init(context.TODO())
	if err != nil {
		logx.Error("设备场景数据初始化失败 err:", err)
		os.Exit(-1)
	}
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
		deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
	}

	tr := schema.NewReadRepo(func(ctx context.Context, productID string) (*schema.Model, error) {
		info, err := productM.ProductSchemaTslRead(ctx, &dm.ProductSchemaTslReadReq{ProductID: productID})
		if err != nil {
			return nil, err
		}
		return schema.ValidateWithFmt([]byte(info.Tsl))
	})
	if c.DiRpc.Mode == conf.ClientModeGrpc {
		deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DiRpc.Conf))
		deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DiRpc.Conf))
	} else {
		deviceMsg = didirect.NewDeviceMsg(c.DiRpc.RunProxy)
		deviceInteract = didirect.NewDeviceInteract(c.DiRpc.RunProxy)
	}

	return &ServiceContext{
		Config: c,
		SvrClient: SvrClient{
			ProductM:       productM,
			DeviceInteract: deviceInteract,
			DeviceMsg:      deviceMsg,
			DeviceM:        deviceM,
		},
		Repo: Repo{
			Store:               store,
			SceneRepo:           SceneRepo,
			SceneDeviceRepo:     sceneDevice,
			SchemaRepo:          tr,
			SceneInfoRepo:       SceneInfoRepo,
			AlarmInfoRepo:       AlarmInfoRepo,
			AlarmSceneRepo:      AlarmSceneRepo,
			AlarmDealRecordRepo: AlarmDealRecordRepo,
			AlarmLogRepo:        AlarmLogRepo,
			AlarmRecordRepo:     AlarmRecordRepo,
		},
	}
}
