package clients

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

func NewAsynqClient(c redis.RedisKeyConf) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: c.Host, Password: c.Pass})
}

func NewAsynqServer(c redis.RedisKeyConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Host, Password: c.Pass},
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
func NewAsynqScheduler(c redis.RedisKeyConf) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Host,
			Password: c.Pass,
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
