package utils

import (
	"context"
	"time"
)

func GetDeadLine(ctx context.Context, defaultDeadLine time.Time) time.Time {
	dead, ok := ctx.Deadline()
	if !ok {
		return defaultDeadLine
	}
	return dead
}
