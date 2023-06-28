package sceneTimer

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/timers"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/repoComplex"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/internal/timer"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type SceneTimer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	timer timers.TimerControl
}

var (
	sceneTimer SceneTimer
	once       sync.Once
)

const (
	singletonKey = "iThingsSceneTimer"
)

func NewSceneTimer(ctx context.Context, svcCtx *svc.ServiceContext) *SceneTimer {
	once.Do(func() {
		sceneTimer.timer = timers.NewTimer(ctx, singletonKey, svcCtx.Store)
		sceneTimer.ctx = ctx
		sceneTimer.svcCtx = svcCtx
		sceneTimer.Logger = logx.WithContext(ctx)
	})
	return &sceneTimer
}

func NewSceneTimerControl() timer.SceneControl {
	return &sceneTimer
}

func (s *SceneTimer) Start() {
	s.timer.Start(func(ctx context.Context) {
		infos, err := s.svcCtx.SceneRepo.FindByFilter(s.ctx, scene.InfoFilter{
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
	})
}

func (s *SceneTimer) Create(info *scene.Info) error {
	keys := []string{genInfoIDKey(info.ID), genInfoNameKey(info.Name)}
	return s.timer.Create(keys, info.Trigger.Timer.Cron, func(ctx context.Context) error {
		return s.jobRun(ctx, info)
	})
}

func (s *SceneTimer) IsRunning() bool {
	return s.timer.IsRunning()
}

func (s *SceneTimer) Update(info *scene.Info) error {
	if !s.timer.IsRunning() {
		return nil
	}
	err := s.Delete(info.ID)
	if err != nil {
		return err
	}
	return s.Create(info)
}

func (s *SceneTimer) Delete(id int64) error {
	if !s.timer.IsRunning() {
		return nil
	}
	return s.timer.Delete(genInfoIDKey(id))
}

func genInfoIDKey(id int64) string {
	return fmt.Sprintf("id:%v", id)
}

func genInfoNameKey(name string) string {
	return "name:" + name
}

func (s *SceneTimer) jobRun(ctx context.Context, info *scene.Info) (err error) {
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
