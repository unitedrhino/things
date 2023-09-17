package timers

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"go.uber.org/atomic"
	"time"
)

type RunFunc func(ctx context.Context) error

type TimerControl interface {
	Create(keys []string, cron string, runFunc RunFunc) error
	IsRunning() bool
	Delete(key string) error
	Start(startFunc func(ctx context.Context))
}

type Timer struct {
	ctx       context.Context
	scheduler *gocron.Scheduler
	logx.Logger
	isRunning    *atomic.Bool
	singletonKey string //锁定的key
	store        kv.Store
}

func NewTimer(ctx context.Context, singletonKey string, store kv.Store) *Timer {
	timer := Timer{
		ctx:          ctx,
		scheduler:    gocron.NewScheduler(time.Local),
		Logger:       logx.WithContext(ctx),
		isRunning:    atomic.NewBool(false),
		singletonKey: singletonKey,
		store:        store,
	}
	return &timer
}

func (s *Timer) Start(startFunc func(ctx context.Context)) {
	utils.Go(s.ctx, func() {
		for true { //定时任务为单例执行模式,有效期15秒,如果服务挂了,其他服务每隔10秒检测到就抢到执行
			ok, err := s.store.SetnxExCtx(s.ctx, s.singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
			if err != nil {
				s.Errorf("%s.Store.SetnxExCtx singletonKey:%v err:%v", utils.FuncName(), s.singletonKey, err)
				time.Sleep(time.Second * 10)
				continue
			}
			if ok { //抢到锁了
				break
			}
			//没抢到锁,10秒钟后继续
			time.Sleep(time.Second * 10)
		}
		s.Infof("Timer start running singletonKey:%v", s.singletonKey)
		//抢到锁需要维系锁
		s.keepSingleton()
		startFunc(s.ctx)
		s.isRunning.Store(true)
		s.scheduler.StartAsync()
	})
}

func (s *Timer) IsRunning() bool {
	return s.isRunning.Load()
}

func (s *Timer) Create(keys []string, cron string, runFunc RunFunc) error {
	if !s.isRunning.Load() {
		return nil
	}
	_, err := s.scheduler.Tag(keys...).CronWithSeconds(cron).Do(func() (err error) {
		ctx, span := ctxs.StartSpan(s.ctx, fmt.Sprintf("SceneTimer.jobRun"), "")
		defer span.End()
		defer utils.Recover(ctx)
		startTime := time.Now().UnixMilli()
		defer logx.WithContext(ctx).Infof("%s.timer.jobRun end use:%vms singletonKey:%v keys:%v cron:%v err:%v",
			utils.FuncName(), time.Now().UnixMilli()-startTime, s.singletonKey, keys, cron, err)
		err = runFunc(ctx)
		return err
	})
	return err
}

func (s *Timer) Delete(key string) error {
	if !s.isRunning.Load() {
		return nil
	}
	jobs, err := s.scheduler.FindJobsByTag(key)
	if err != nil {
		s.Errorf("%s.FindJobsByTag err:%v", utils.FuncName(), err)
		return nil
	}
	for _, job := range jobs {
		s.scheduler.Remove(job)
	}
	return nil
}

func (s *Timer) keepSingleton() {
	//每隔10秒刷新锁,如果服务挂了,锁才能退出
	utils.Go(s.ctx, func() {
		defer utils.Recover(s.ctx)
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			s.store.SetexCtx(s.ctx, s.singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
		}
	})
}
