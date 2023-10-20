package timedDelay

import (
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"
)

type Option struct {
	Priority string //优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	//以下两个参数优先使用ProcessIn
	ProcessIn time.Duration //多久之后发
	ProcessAt time.Time     // 固定时间发
	Timeout   time.Duration //超时时间 优先使用
	Deadline  time.Time     //截止时间
}

type Timed struct {
	asynqClient *asynq.Client
	serverName  string //服务名
}

func NewTimed(c cache.ClusterConf, serverName string) *Timed {
	return &Timed{asynqClient: clients.NewAsynqClient(c), serverName: serverName}
}
func (t Timed) Enqueue(j *task.Info, option *Option) error {
	if option != nil {
		j.Priority = option.Priority
	}
	err := j.Init()
	if err != nil {
		return err
	}
	var opts []asynq.Option
	if option != nil {
		var opt = asynq.ProcessAt(option.ProcessAt)
		if option.ProcessIn != 0 {
			opt = asynq.ProcessIn(option.ProcessIn)
		}
		opts = append(opts, opt)
		if option.Timeout != 0 {
			opts = append(opts, asynq.Timeout(option.Timeout))
		}
		if !option.Deadline.IsZero() {
			opts = append(opts, asynq.Deadline(option.Deadline))
		}
	}
	_, err = t.asynqClient.Enqueue(j.ToTask(), opts...)
	return err
}
