package tdengine

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func AffiliationToMap(in devices.Affiliation) map[string]any {
	var ret = make(map[string]any)
	if in.TenantCode != "" {
		ret["tenant_code"] = in.TenantCode
	}
	if in.ProjectID != 0 {
		ret["project_id"] = in.ProjectID
	}
	if in.AreaIDPath != "" {
		ret["area_id_path"] = in.AreaIDPath
	}
	if in.AreaID != 0 {
		ret["area_id"] = in.AreaID
	}
	return ret
}

func AlterTag(ctx context.Context, t *clients.Td, tables []string, tags map[string]any) error {
	for _, table := range tables {
		for k, v := range tags {
			_, err := t.ExecContext(ctx, fmt.Sprintf(" ALTER TABLE %s SET TAG `%s`='%v'; ",
				table, k, v))
			if err != nil {
				if strings.Contains(err.Error(), "Table does not exist") {
					break
				}
				logx.Error(err)
			}
		}
	}
	return nil
}

type Tag struct {
	Table string
	Tags  map[string]any
}

func AlterTags(ctx context.Context, t *clients.Td, tags []Tag) error {
	for _, tag := range tags {
		for k, v := range tag.Tags {
			_, err := t.ExecContext(ctx, fmt.Sprintf(" ALTER TABLE %s SET TAG `%s`='%v'; ",
				tag.Table, k, v))
			if err != nil {
				if strings.Contains(err.Error(), "Table does not exist") {
					break
				}
				logx.Error(err)
			}
		}

	}
	return nil
}

func GetTdType(define schema.Define) string {
	switch define.Type {
	case schema.DataTypeBool:
		return "BOOL"
	case schema.DataTypeInt:
		return "BIGINT"
	case schema.DataTypeString:
		return "BINARY(5000)"
	case schema.DataTypeStruct:
		return "BINARY(5000)"
	case schema.DataTypeFloat:
		return "DOUBLE"
	case schema.DataTypeTimestamp:
		return "TIMESTAMP"
	case schema.DataTypeArray:
		return "BINARY(5000)"
	case schema.DataTypeEnum:
		return "SMALLINT"
	default: //走到这里说明前面没有进行校验需要检查是否是前面有问题
		panic(fmt.Sprintf("%v not support", define.Type))
	}
}
