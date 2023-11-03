package timer

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Timed struct {
	SvcCtx *svc.ServiceContext
}

func (t Timed) ProcessTask(ctx context.Context, Task *asynq.Task) error {
	defer func() {
		utils.Recover(ctx)
	}()
	err := func() error {
		ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
		defer cancel()
		ctx, span := ctxs.StartSpan(ctx, "timedJob.ProcessTask", "")
		defer span.End()
		utils.Recover(ctx)
		var taskInfo domain.TaskInfo
		json.Unmarshal(Task.Payload(), &taskInfo)
		tr := relationDB.NewTaskRepo(ctx)
		task, err := tr.FindOneByFilter(ctx, relationDB.TaskFilter{
			IDs:       []int64{taskInfo.ID},
			WithGroup: true,
		})
		if err != nil {
			return err
		}
		if task.Type == domain.TaskTypeTiming && task.Status != def.StatusRunning { //如果没有处于运行中,任务不能执行
			return nil
		}
		err = FillTaskInfoDo(&taskInfo, task)
		if err != nil {
			return err
		}
		logx.WithContext(ctx).Infof("timedJob ProcessTask task:%v", utils.Fmt(taskInfo))

		switch task.Group.Type {
		case domain.TaskGroupTypeQueue:
			return t.Queue(ctx, &taskInfo)
		case domain.TaskGroupTypeSql:
			return t.SqlExec(ctx, &taskInfo)
		}
		logx.WithContext(ctx).Errorf("not support job type:%v", task.Group.Type)
		return nil
	}()
	if err != nil {
		logx.WithContext(ctx).Errorf("ProcessTask err:%v task:%v", err, utils.Fmt(Task))
	}
	return nil
}
