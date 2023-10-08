package svc

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/src/timedschedulersvr/internal/config"
	"time"
)

// create scheduler
func newScheduler(c config.Config) *asynq.Scheduler {

	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Host,
			Password: c.Redis.Pass,
		}, &asynq.SchedulerOpts{
			Location: location,
			PostEnqueueFunc: func(task *asynq.TaskInfo, err error) {
				if err == nil {
					return
				}
				fmt.Printf("Scheduler EnqueueErrorHandler <<<<<<<===>>>>> err : %+v , task : %+v", err, task)
			},
		})
}
