package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/domain/job"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/svc"
)

type Timed struct {
	SvcCtx *svc.ServiceContext
}

func (t Timed) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var jb job.Job
	json.Unmarshal(task.Payload(), &jb)
	fmt.Println(jb)
	err := jb.Init()
	if err != nil {
		return err
	}
	switch jb.Type {
	case job.JobTypeQueue:
		return t.SvcCtx.PubJob.Publish(ctx, jb.SubType, jb.Queue.Topic, []byte(jb.Queue.Payload))
	}
	return nil
}
