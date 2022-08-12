// Package schema 这个文件定义和模板相关的接口及dto定义
package schema

import (
	"context"
	"time"
)

type (
	SchemaInfo struct {
		ProductID   string // 产品id
		Template    string // 数据模板
		CreatedTime time.Time
	}
	SchemaRepo interface {
		Insert(ctx context.Context, productID string, template *Model) error
		GetSchemaModel(ctx context.Context, productID string) (*Model, error)
		GetSchemaInfo(ctx context.Context, productID string) (*SchemaInfo, error)
		Update(ctx context.Context, productID string, template *Model) error
		Delete(ctx context.Context, productID string) error
		ClearCache(ctx context.Context, productID string) error
	}
	GetSchemaModel func(ctx context.Context, productID string) (*Model, error)
)
