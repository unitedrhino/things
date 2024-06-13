package timerEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	rulelogic "github.com/i-Things/things/service/udsvr/internal/logic/rule"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (l *TimerHandle) SceneThingPropertyReport(in application.PropertyReport) error {
	now := time.Now()
	db := stores.WithNoDebug(l.ctx, relationDB.NewSceneIfTriggerRepo)

	list, err := db.FindByFilter(l.ctx, relationDB.SceneIfTriggerFilter{
		Status: def.True,
		Type:   scene.TriggerTypeDevice,
		Device: &in.Device,
		DataID: in.Identifier,
	}, nil)
	if err != nil {
		return err
	}
	for _, v := range list {
		var po = v
		if po.SceneInfo == nil {
			logx.WithContext(l.ctx).Errorf("trigger device not bind scene, trigger:%v", utils.Fmt(po))
			relationDB.NewSceneIfTriggerRepo(l.ctx).Delete(l.ctx, po.ID)
			continue
		}
		po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
		do := rulelogic.PoToSceneInfoDo(po.SceneInfo)

		ps, err := l.svcCtx.ProductSchemaCache.GetData(l.ctx, in.Device.ProductID)
		if err != nil {
			l.Error(err)
			continue
		}
		if !do.If.Triggers[0].Device.IsHit(ps, in.Identifier, in.Param) {
			continue
		}
		ctx := ctxs.BindTenantCode(l.ctx, string(v.SceneInfo.TenantCode), int64(v.SceneInfo.ProjectID))
		if !do.When.IsHit(ctx, now, nil) {
			continue
		}
		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
			var err error
			func() {
				po.LastRunTime = utils.GetEndTime(now)
				err = db.Update(ctx, po)
				if err != nil { //如果失败了下次还可以执行
					l.Error(err)
					return
				}
			}()
			if err != nil { //如果失败了下次还可以执行
				l.Error(err)
				return
			}
			l.SceneExec(ctx, do)
		})
	}
	return nil
}
