// Package repo 本文件是提供设备模型数据存储的信息
package deviceMsgManage

import (
	"context"
	"github.com/i-Things/things/shared/domain/schema"
)

type (
	SchemaDataRepo interface {
		// InitProduct 初始化产品的物模型相关表及日志记录表
		InitProduct(ctx context.Context, t *schema.Model, productID string) error
		// DropProduct 删除产品时需要删除产品下的所有表
		DropProduct(ctx context.Context, t *schema.Model, productID string) error
		// InitDevice 创建设备时为设备创建单独的表
		InitDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		// DropDevice 删除设备时需要删除设备的所有表
		DropDevice(ctx context.Context, t *schema.Model, productID string, deviceName string) error
		// ModifyProduct 修改产品物模型 只支持新增和删除,不支持修改数据类型
		ModifyProduct(ctx context.Context, oldT *schema.Model, newt *schema.Model, productID string) error
	}
)
