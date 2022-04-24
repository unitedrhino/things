// Package templateModel 这个文件定义和模板相关的接口及dto定义
package templateModel

import (
	"context"
	"time"
)

type (
	TemplateInfo struct {
		ProductID   string // 产品id
		Template    string // 数据模板
		CreatedTime time.Time
	}
	TemplateRepo interface {
		Insert(ctx context.Context, productID string, template *Template) error
		GetTemplate(ctx context.Context, productID string) (*Template, error)
		GetTemplateInfo(ctx context.Context, productID string) (*TemplateInfo, error)
		Update(ctx context.Context, productID string, template *Template) error
		Delete(ctx context.Context, productID string) error
		ClearCache(ctx context.Context, productID string) error
	}
)
