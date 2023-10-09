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
	"time"
)

type Timed struct {
	SvcCtx *svc.ServiceContext
}

func (t Timed) ProcessTask(ctx context.Context, Task *asynq.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	utils.Recover(ctx)
	var jb task.Info
	json.Unmarshal(Task.Payload(), &jb)
	ctx, span := ctxs.StartSpan(ctx, fmt.Sprintf("timedJob_%s", jb.Code), "")
	defer span.End()
	err := jb.Init()
	if err != nil {
		return err
	}
	switch jb.Type {
	case task.TaskTypeQueue:
		return t.SvcCtx.PubJob.Publish(ctx, jb.SubType, jb.Queue.Topic, []byte(jb.Queue.Payload))
	}
	return nil
}
