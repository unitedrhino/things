package startup

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/timer"
)

func Init(svcCtx *svc.ServiceContext) error {
	utils.Go(context.Background(), func() {
		utils.SingletonRun(context.Background(), svcCtx.Store, "svr:timedschedulersvr", func(ctx2 context.Context) {
			svcCtx.SchedulerRun = true
			timer.Run(svcCtx)
		})
	})
	return nil
}
