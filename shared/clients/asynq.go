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
func NewTimedScheduler(c cache.ClusterConf) *TimedScheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return &TimedScheduler{Asynq: asynq.NewScheduler(
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
		}), run: make(map[string]string, 100)}
}

type TimedScheduler struct {
	Asynq *asynq.Scheduler
	run   map[string]string //key是任务code,value 是entryID
}

func (s *TimedScheduler) Register(cronspec string, taskCode string, payload []byte, opts ...asynq.Option) (err error) {
	t, ok := s.run[taskCode]
	if ok { //如果正在运行,需要先删除再注册
		err = s.Unregister(t)
		if err != nil {
			return err
		}
	}
	task := asynq.NewTask(taskCode, payload, opts...)
	entryID, err := s.Asynq.Register(cronspec, task)
	if err != nil {
		return err
	}
	s.run[taskCode] = entryID
	return
}

func (s *TimedScheduler) Unregister(taskCode string) error {
	t, ok := s.run[taskCode]
	if ok { //如果正在运行,需要先删除再注册
		err := s.Asynq.Unregister(t)
		if err != nil {
			return err
		}
		delete(s.run, taskCode)
	}
	return nil
}

func (s *TimedScheduler) Run() error {
	return s.Asynq.Run()
}
