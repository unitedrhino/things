package clients

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"
)

func NewAsynqClient(c cache.ClusterConf) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: c[0].Host, Password: c[0].Pass})
}

func NewAsynqServer(c cache.ClusterConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c[0].Host, Password: c[0].Pass},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>>  err : %+v \n", err)
				return true
			},
			Concurrency: 20, //max concurrent process job task num
			Queues: map[string]int{
				"critical": 3,
				"default":  2,
				"low":      1,
			},
			StrictPriority: true,
		},
	)
}

// create scheduler
func NewAsynqScheduler(c cache.ClusterConf) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c[0].Host,
			Password: c[0].Pass,
		}, &asynq.SchedulerOpts{
			Location: location,
			PostEnqueueFunc: func(task *asynq.TaskInfo, err error) {
				if err == nil {
					return
				}
				fmt.Printf("Scheduler PostEnqueueFunc  err : %+v , task : %+v", err, task)
			},
		})
}
