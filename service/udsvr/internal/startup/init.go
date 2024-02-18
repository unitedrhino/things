package startup

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	TimerInit(svcCtx)
	Subscribe(svcCtx)
	InitEventBus(svcCtx)
}

func Subscribe(svcCtx *svc.ServiceContext) {

}

func InitEventBus(svcCtx *svc.ServiceContext) {
}

func TimerInit(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
		GroupCode: def.TimedIThingsQueueGroupCode,                    //组编码
		Type:      1,                                                 //任务类型 1 定时任务 2 延时任务
		Name:      "iThings规则引擎定时任务",                                 // 任务名称
		Code:      "iThingsRuleTimer",                                //任务编码
		Params:    `{"topic":"server.things.ud.timer","payload":""}`, // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
		CronExpr:  "@every 1s",                                       // cron执行表达式
		Status:    def.StatusWaitRun,                                 // 状态
		Priority:  3,                                                 //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	})
	if !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
		logx.Must(err)
	}
}
