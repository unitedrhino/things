package middlewares

import (
	"context"
	"fmt"
	"testing"
)

func TestWithMiddlewares(t *testing.T) {
	WithMiddlewares("123", func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, value any) {
			fmt.Println("first 1")
			fmt.Println(value)
			next(ctx, value)
		}
	}, func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, value any) {
			fmt.Println("first 2")
			fmt.Println(value)
			next(ctx, value)
		}
	}, func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, value any) {
			fmt.Println("first 3")
			fmt.Println(value)
			next(ctx, value)
		}
	}, func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, value any) {
			fmt.Println("first 4")
			fmt.Println(value)
			next(ctx, value)
		}
	})
	Execute("123", context.TODO(), "test1")
}
