package viddirect

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
	"github.com/i-Things/things/src/vidsvr/internal/handler"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
	"time"
)

// vidsvr构建一个http server 以便于zlmediakit hooks访问

type (
	ServiceContext = svc.ServiceContext
	ApiCtx         struct {
		Server *rest.Server
		SvcCtx *ServiceContext
	}
)

var (
	timeObj    timedmanage.TimedManage
	runApiOnce sync.Once
	//threadOnce  sync.Once
	initSvrOnce sync.Once
	cronTask    *cron.Cron
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	var server = apiCtx.Server
	ctx := GetSvcCtx()
	apiCtx.SvcCtx = ctx
	if apiCtx.Server == nil {
		server = rest.MustNewServer(ctx.Config.Restconf)
		apiCtx.Server = server
	}
	handler.RegisterHandlers(server, ctx)
	return apiCtx
}

func ApiDirectRun() {
	runApiOnce.Do(func() {
		utils.Go(context.Background(), ApiRun) //golang 后台执行
	})
}

func ApiRun() {
	apiCtx := NewApi(ApiCtx{})
	apiCtx.Server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n",
		apiCtx.SvcCtx.Config.Restconf.Host, apiCtx.SvcCtx.Config.Restconf.Port)
	defer apiCtx.Server.Stop()
	initSvrOnce.Do(func() {
		//InitData()
		utils.Go(context.Background(), DelayTask)
		utils.Go(context.Background(), CronTask)
	})

	//初始化第一个流服务
	apiCtx.Server.Start()
}

func CronTask() {
	cronTask = cron.New(cron.WithSeconds())
	fmt.Println(time.Now())
	//cronTask.AddFunc("*/30 * * * * *",func)
	cronTask.AddFunc("*/30 * * * * *", func() {
		media.SrvInfoStatusCheck()
	})
	cronTask.Start()
	select {}
}

func InitData() {
	ctx := GetSvcCtx()
	//sendTime := time.Now()
	fmt.Printf("ctx.Config.TimedJobRpc.Enable: %v ...\n", ctx.Config.TimedJobRpc.Enable)
	fmt.Printf("InitData send nats: %s ...\n", topics.VidInfoInitDatabase)
	if ctx.Config.TimedJobRpc.Enable {

		if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
			timeObj = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
		} else {
			timeObj = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
		}
		//发步一个延时任务  初始化docker-zlmediakit
		timeObj.TaskSend(context.Background(), &timedjob.TaskSendReq{
			GroupCode: def.TimedIThingsQueueGroupCode,
			Code:      "VidInfoInitDatabase",
			Option: &timedjob.TaskSendOption{
				ProcessIn: 2,
				Deadline:  0,
			},
			ParamQueue: &timedjob.TaskParamQueue{
				Topic:   topics.VidInfoInitDatabase,
				Payload: string(topics.VidInfoInitDatabase),
			},
		})
	}
}

func DelayTask() {
	time.Sleep(5 * time.Second) //5 秒后执行
	id := deviceAuth.GetStrProductID(svcCtx.VidmgrID.GetSnowflakeId())
	media.InitDockerSrv(c, id)
}
