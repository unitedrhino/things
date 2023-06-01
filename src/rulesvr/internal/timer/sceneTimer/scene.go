package sceneTimer

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/repoComplex"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/internal/timer"
	"github.com/zeromicro/go-zero/core/logx"
	atomic2 "go.uber.org/atomic"
	"sync"
	"time"
)

type SceneTimer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	scheduler *gocron.Scheduler
}

var (
	sceneTimer SceneTimer
	once       sync.Once
	isRunning  = atomic2.NewBool(false)
)

const (
	singletonKey = "iThingsSceneTimer"
)

func NewSceneTimer(ctx context.Context, svcCtx *svc.ServiceContext) *SceneTimer {
	once.Do(func() {
		sceneTimer.ctx = ctx
		sceneTimer.svcCtx = svcCtx
		sceneTimer.Logger = logx.WithContext(ctx)
		sceneTimer.scheduler = gocron.NewScheduler(time.Local)
	})
	return &sceneTimer
}
func NewSceneTimerControl() timer.SceneControl {
	return &sceneTimer
}

func (s *SceneTimer) Start() {
	go func() {
		defer utils.Recover(s.ctx)
		for true { //定时任务为单例执行模式,有效期15秒,如果服务挂了,其他服务每隔10秒检测到就抢到执行
			ok, err := s.svcCtx.Store.SetnxExCtx(s.ctx, singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
			if err != nil {
				s.Errorf("%s.Store.SetnxExCtx err:%v", utils.FuncName(), err)
				time.Sleep(time.Second * 10)
				continue
			}
			if ok { //抢到锁了
				break
			}
			//没抢到锁,10秒钟后继续
			time.Sleep(time.Second * 10)
		}
		s.Infof("SceneTimer start running")
		//抢到锁需要维系锁
		s.keepSingleton()
		s.run()
	}()

}

func (s *SceneTimer) keepSingleton() {
	//每隔10秒刷新锁,如果服务挂了,锁才能退出
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			s.svcCtx.Store.SetexCtx(s.ctx, singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
		}
	}()
}
func (s *SceneTimer) run() {
	isRunning.Store(true)
	infos, err := s.svcCtx.SceneRepo.FindByFilter(s.ctx, scene.InfoFilter{
		Name:        "",
		Status:      def.Enable,
		TriggerType: scene.TriggerTypeTimer,
	}, nil)
	if err != nil {
		s.Errorf("%s.SceneRepo.FindByFilter err:%v", utils.FuncName(), err)
		return
	}
	for _, info := range infos {
		s.Create(info)
	}
	s.scheduler.StartAsync()
}
func (s *SceneTimer) Create(info *scene.Info) error {
	if !isRunning.Load() {
		return nil
	}
	_, err := s.scheduler.Tag(genInfoIDKey(info.ID), genInfoNameKey(info.Name)).CronWithSeconds(info.Trigger.Timer.Cron).Do(func(info *scene.Info) {
		s.jobRun(info)
	}, info)
	return err
}
func (s *SceneTimer) Update(info *scene.Info) error {
	if !isRunning.Load() {
		return nil
	}
	err := s.Delete(info.ID)
	if err != nil {
		return err
	}
	_, err = s.scheduler.Tag(genInfoIDKey(info.ID), genInfoNameKey(info.Name)).CronWithSeconds(info.Trigger.Timer.Cron).Do(func(info *scene.Info) {
		s.jobRun(info)
	}, info)
	return err
}
func (s *SceneTimer) Delete(id int64) error {
	if !isRunning.Load() {
		return nil
	}
	jobs, err := s.scheduler.FindJobsByTag(genInfoIDKey(id))
	if err != nil {
		s.Errorf("%s.FindJobsByTag err:%v", err)
		return err
	}
	for _, job := range jobs {
		s.scheduler.Remove(job)
	}
	return nil
}
func genInfoIDKey(id int64) string {
	return fmt.Sprintf("id:%v", id)
}
func genInfoNameKey(name string) string {
	return "name:" + name
}
func (s *SceneTimer) jobRun(info *scene.Info) (err error) {
	ctx, span := traces.StartSpan(s.ctx, fmt.Sprintf("SceneTimer.jobRun"), "")
	defer span.End()
	if len(info.When) != 0 {
		if !info.When.IsHit(ctx, scene.TermRepo{
			DeviceMsg:  s.svcCtx.DeviceMsg,
			SchemaRepo: s.svcCtx.SchemaRepo,
		}) {
			s.Infof("%s timer when not commit scene name:%v id:%v", utils.FuncName(), info.Name, info.ID)
			return
		}
	}
	startTime := time.Now().UnixMilli()
	defer logx.WithContext(ctx).Infof("%s.timer.jobRun end use:%vms sceneName:%v err:%v",
		utils.FuncName(), time.Now().UnixMilli()-startTime, info.Name, err)
	logx.WithContext(ctx).Infof("%s.timer.jobRun start sceneID:%v sceneName:%v", utils.FuncName(), info.ID, info.Name)
	err = info.Then.Execute(ctx, scene.ActionRepo{
		DeviceInteract: s.svcCtx.DeviceInteract,
		DeviceM:        s.svcCtx.DeviceM,
		Alarm:          repoComplex.NewSceneAlarm(s.svcCtx),
		Scene:          info,
	})
	return err
}
