package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timedjobsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Timed struct {
	SvcCtx *svc.ServiceContext
}

func (t Timed) ProcessTask(ctx context.Context, Task *asynq.Task) error {
	err := func() error {
		ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
		defer cancel()
		utils.Recover(ctx)
		var jb task.Info
		json.Unmarshal(Task.Payload(), &jb)
		ctx, span := ctxs.StartSpan(ctx, fmt.Sprintf("timedJob_%s", jb.Code), "")
		defer span.End()
		logx.WithContext(ctx).Infof("timedJob ProcessTask task:%v", jb)
		err := jb.Init()
		if err != nil {
			return err
		}
		switch jb.Type {
		case task.TaskTypeQueue:
			return t.Queue(ctx, &jb)
		case task.TaskTypeSql:
			return t.SqlExec(ctx, &jb)
		}
		logx.WithContext(ctx).Errorf("not support job type:%v", jb.Type)
		return nil
	}()
	if err != nil {
		logx.WithContext(ctx).Errorf("ProcessTask err:%v task:%v", err, utils.Fmt(Task))
	}
	return nil
}
