package utils

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime"
	"runtime/debug"
)

func Recover(ctx context.Context) {
	if p := recover(); p != nil {
		HandleThrow(ctx, p)
	}
}

func HandleThrow(ctx context.Context, p any) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	logx.WithContext(ctx).Errorf("HandleThrow|func=%s|error=%#v|stack=%s\n", f, p, string(debug.Stack()))
	//os.Exit(-1)
}

func Go(ctx context.Context, f func()) {
	go func() {
		defer Recover(ctx)
		f()
	}()
}
