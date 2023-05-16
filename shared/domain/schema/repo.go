// Package schema 这个文件定义和模板相关的接口及dto定义
package schema

import (
	"context"
	"time"
)

type (
	Info struct {
		ProductID   string // 产品id
		Schema      string // 数据模板
		CreatedTime time.Time
	}
	Repo interface {
		ReadRepo
		TslImport(ctx context.Context, productID string, template *Model) error
		//Update(ctx context.Context, productID string, template *Model) error
		TslRead(ctx context.Context, productID string) (*Model, error)
		Delete(ctx context.Context, productID string) error
	}
	ReadRepo interface {
		GetSchemaModel(ctx context.Context, productID string) (*Model, error)
		//GetReadInfo(ctx context.Context, productID string) (*Info, error)
		ClearCache(ctx context.Context, productID string) error
	}
	GetSchemaModel func(ctx context.Context, productID string) (*Model, error)
)
