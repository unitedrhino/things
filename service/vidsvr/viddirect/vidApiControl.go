package viddirect

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/core/shared/domain/deviceAuth"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsvr/internal/handler"
	"github.com/i-Things/things/service/vidsvr/internal/media"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/rest"
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

func DelayTask() {
	time.Sleep(5 * time.Second) //5 秒后执行
	id := deviceAuth.GetStrProductID(svcCtx.VidmgrID.GetSnowflakeId())
	media.InitDockerSrv(c, id)
}
