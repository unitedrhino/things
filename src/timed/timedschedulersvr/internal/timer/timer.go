package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

func Run(svcCtx *svc.ServiceContext) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//ddsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "timedSchedulersvr.taskRun", "")
	defer span.End()
	{ //先初始化数据库状态
		msg := "初始化数据库执行错误"
		jobDB := relationDB.NewTaskRepo(ctx)
		//将运行中的任务修改为等待运行
		err := jobDB.UpdateByFilter(ctx, &relationDB.TimedTask{Status: def.StatusWaitRun},
			relationDB.TaskFilter{Status: []int64{def.StatusRunning}, Types: []int64{domain.TaskTypeTiming}})
		errors.Must(err, msg)
		//将等待暂停的任务调整为已暂停
		err = jobDB.UpdateByFilter(ctx, &relationDB.TimedTask{Status: def.StatusStopped},
			relationDB.TaskFilter{Status: []int64{def.StatusWaitStop}, Types: []int64{domain.TaskTypeTiming}})
		errors.Must(err, msg)
		//删除等待删除的任务
		err = jobDB.DeleteByFilter(ctx, relationDB.TaskFilter{Status: []int64{def.StatusWaitDelete}, Types: []int64{domain.TaskTypeTiming}})
		errors.Must(err, msg)
	}
	utils.Go(ctx, func() {
		ctx := context.Background()
		TaskCheck(svcCtx)
		utils.Go(ctx, func() {
			err := svcCtx.Scheduler.Run()
			errors.Must(err, "Scheduler.Run")
		})
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				TaskCheck(svcCtx)
			}
		}
	})

}

func TaskCheck(svcCtx *svc.ServiceContext) {
	logx.Info("TaskCheck run")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//ddsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "timedSchedulersvr.taskCheck", "")
	defer span.End()
	err := func() error {
		jobDB := relationDB.NewTaskRepo(ctx)
		js, err := jobDB.FindByFilter(ctx, relationDB.TaskFilter{WithGroup: true,
			Status: []int64{def.StatusWaitStop, def.StatusWaitDelete, def.StatusWaitRun},
			Types:  []int64{domain.TaskTypeTiming}},
			&def.PageInfo{
				Orders: []def.OrderBy{{Filed: "priority", Sort: def.OrderDesc}},
			})
		if err != nil {
			return err
		}
		wait := sync.WaitGroup{}
		for _, j := range js {
			wait.Add(1)
			t := j
			utils.Go(ctx, func() {
				err := func() error {
					switch t.Status {
					case def.StatusWaitRun:
						return TaskStatusRunCheck(ctx, svcCtx, &wait, t)
					case def.StatusWaitDelete, def.StatusWaitStop:
						return TaskStatusStopCheck(ctx, svcCtx, &wait, t)
					}
					//其他状态不需要处理
					return nil
				}()
				if err != nil {
					logx.WithContext(ctx).Errorf("TaskCheck.one  err:%+v , task:%+v", err, t)
				}
			})
		}
		wait.Wait()
		return nil
	}()
	if err != nil {
		logx.WithContext(ctx).Errorf("TaskCheck  err:%v", err)
	}
}

func getTaskCode(j *relationDB.TimedTask) string {
	return fmt.Sprintf("timing:%s:%s", j.GroupCode, j.Code)
}

// 需要检查任务是否启动,如果没有启动需要启动
func TaskStatusRunCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTask) error {
	defer wait.Done()
	taskCode := getTaskCode(task)
	taskInfo := domain.TaskInfo{
		ID:     task.ID,
		Params: "",
	}
	payload, _ := json.Marshal(taskInfo)
	err := svcCtx.Scheduler.Register(task.CronExpr, taskCode, payload, asynq.Queue(domain.ToPriority(task.Priority)))
	if err != nil {
		logx.WithContext(ctx).Errorf("TaskStatusRunCheck.Register err:%v task:%v", err, task)
		return errors.System.AddDetail(err)
	}
	jobDB := relationDB.NewTaskRepo(ctx)
	task.Status = def.StatusRunning
	err = jobDB.Update(ctx, task)
	return err
}

// 如果处于运行状态需要停止
func TaskStatusStopCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTask) error {
	defer wait.Done()
	taskCode := getTaskCode(task)
	err := svcCtx.Scheduler.Unregister(taskCode)
	if err != nil {
		logx.WithContext(ctx).Errorf("TaskStatusStopCheck.Unregister err:%v task:%v", err, task)
		return errors.System.AddDetail(err)
	}
	jobDB := relationDB.NewTaskRepo(ctx)
	switch task.Status {
	case def.StatusWaitDelete:
		err = jobDB.Delete(ctx, task.ID)
		if err != nil {
			return err
		}
	case def.StatusWaitStop:
		task.Status = def.StatusStopped
		err = jobDB.Update(ctx, task)
	}
	return err
}
