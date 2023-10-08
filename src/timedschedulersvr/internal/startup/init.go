package startup

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/job"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) error {
	return InitTimer(svcCtx)
}

func InitTimer(svcCtx *svc.ServiceContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//ddsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "InitTimer", "")
	defer span.End()
	jobDB := relationDB.NewJobRepo(ctx)
	//先把状态全部改成暂停状态
	err := jobDB.UpdateByFilter(ctx, &relationDB.TimedQueueJob{Status: relationDB.JobStatusPause},
		relationDB.JobFilter{Status: relationDB.JobStatusRun})
	if err != nil {
		return err
	}
	js, err := jobDB.FindByFilter(ctx, relationDB.JobFilter{Status: relationDB.JobStatusPause}, &def.PageInfo{
		Orders: []def.OrderBy{{Filed: "priority", Sort: def.OrderDesc}},
	})
	if err != nil {
		return err
	}
	for _, j := range js {
		jb := job.Job{
			Group:    j.Group,
			Type:     j.Type,
			SubType:  j.SubType,
			Name:     j.Name,
			Code:     j.Code,
			Params:   j.Params,
			Priority: j.Priority,
		}
		err := jb.Init()
		if err != nil {
			logx.WithContext(ctx).Errorf("job init  err:%+v , job:%+v", err, jb)
			continue
		}
		task := jb.ToTask()
		// every one minute exec
		entryID, err := svcCtx.Scheduler.Register(j.CronExpression, task)
		if err != nil {
			logx.WithContext(ctx).Errorf("Scheduler.Register  err:%+v , task:%+v", err, task)
			continue
		}
		j.EntryID = entryID
		j.Status = relationDB.JobStatusRun
		err = jobDB.Update(ctx, j)
		if err != nil {
			logx.WithContext(ctx).Errorf("Scheduler.Update  err:%+v , task:%+v", err, task)
			svcCtx.Scheduler.Unregister(entryID)
		}
	}
	utils.Go(ctx, func() {
		err := svcCtx.Scheduler.Run()
		logx.Error(err)
	})
	return nil
}
