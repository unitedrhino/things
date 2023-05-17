package middlewares

import "context"

type (
	HandleFunc func(ctx context.Context, value any)
	Middleware func(next HandleFunc) HandleFunc
)

var (
	MiddleHandle = map[string]HandleFunc{}
)

func WithMiddlewares(name string, middlewares ...Middleware) {
	var HandleFuncs HandleFunc = func(ctx context.Context, value any) {}
	if MiddleHandle[name] != nil {
		HandleFuncs = MiddleHandle[name]
	}
	for _, middleware := range middlewares {
		HandleFuncs = middleware(HandleFuncs)
	}
	MiddleHandle[name] = HandleFuncs
}
func Execute(name string, ctx context.Context, value any) {
	if handle := MiddleHandle[name]; handle != nil {
		handle(ctx, value)
	}
}
