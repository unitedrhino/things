package timerEvent

import (
	"context"
	"fmt"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	rulelogic "github.com/i-Things/things/service/udsvr/internal/logic/rule"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TimerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewSceneHandle(ctx context.Context, svcCtx *svc.ServiceContext) *TimerHandle {
	return &TimerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}
func (l *TimerHandle) LockRunning(ctx context.Context, Type string /*scene deviceTimer*/, triggerID int64) (deferF func()) {
	key := fmt.Sprintf("things:rule:%s:trigger:%d", Type, triggerID)
	ok, err := l.svcCtx.Store.SetnxExCtx(ctx, key, time.Now().Format("2006-01-02 15:04:05.999"), 5)
	if err != nil || !ok {
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		return nil
	}
	//抢到锁了
	return func() {
		_, err := l.svcCtx.Store.DelCtx(ctx, key)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}

}

func (l *TimerHandle) SceneExec(ctx context.Context, do *scene.Info) error {
	err := do.Then.Execute(ctx, scene.ActionRepo{
		Info:           do,
		DeviceInteract: l.svcCtx.DeviceInteract,
		DeviceM:        l.svcCtx.DeviceM,
		DeviceG:        l.svcCtx.DeviceG,
		SceneExec: func(ctx context.Context, sceneID int64) error {
			_, err := rulelogic.NewSceneManuallyTriggerLogic(ctx, l.svcCtx).SceneManuallyTrigger(&ud.WithID{Id: sceneID})
			return err
		},
		AlarmExec: func(ctx context.Context, in scene.AlarmSerial) error {
			if len(in.Scene.If.Triggers) == 0 {
				logx.WithContext(ctx).Error("没有触发器")
				return nil
			}
			trigger := in.Scene.If.Triggers[0]
			req := ud.AlarmRecordCreateReq{
				TriggerType: trigger.Type,
				SceneName:   in.Scene.Name,
				SceneID:     in.Scene.ID,
				Mode:        scene.ActionAlarmModeTrigger,
			}
			if trigger.Type == scene.TriggerTypeDevice && trigger.Device != nil {
				req.ProductID = trigger.Device.ProductID
				req.DeviceName = trigger.Device.DeviceName
			}
			_, err := rulelogic.NewAlarmRecordCreateLogic(l.ctx, l.svcCtx).AlarmRecordCreate(&req)
			return err
		},
	})
	return err
}
