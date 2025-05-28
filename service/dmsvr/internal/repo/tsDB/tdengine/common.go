package tdengine

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func ToBelongGroup(in map[string]*dm.IDsInfo) (out map[string]def.IDsInfo) {
	if in == nil {
		return
	}
	out = make(map[string]def.IDsInfo)
	for k, v := range in {
		out[k] = utils.Copy2[def.IDsInfo](v)
	}
	return
}

func AffiliationToMap(in devices.Affiliation, groupConfigs []*deviceGroup.GroupDetail) map[string]any {
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
	if in.BelongGroup != nil && groupConfigs != nil {
		for _, groupConfig := range groupConfigs {
			if !groupConfig.TsIndex {
				continue
			}
			ret[fmt.Sprintf("group_%s_ids", groupConfig.Value)] = utils.GenSliceStr(in.BelongGroup[groupConfig.Value].IDs)
			ret[fmt.Sprintf("group_%s_id_paths", groupConfig.Value)] = utils.GenSliceStr(in.BelongGroup[groupConfig.Value].IDPaths)
		}
	}
	return ret
}

func GenTagsDef(ts string, groupConfigs []*deviceGroup.GroupDetail) (tagNames string) {
	var tags []string
	for _, groupConfig := range groupConfigs {
		if groupConfig.TsIndex == false {
			continue
		}
		tags = append(tags, fmt.Sprintf("`group_%s_ids`  BINARY(250) ", groupConfig.Value), fmt.Sprintf("`group_%s_id_paths`  BINARY(250) ", groupConfig.Value))
	}
	if len(tags) == 0 {
		return ts
	}
	return ts + "," + strings.Join(tags, ",")
}

func GroupFilter(sql sq.SelectBuilder, groupConfigs []*deviceGroup.GroupDetail, BelongGroup map[string]def.IDsInfo) sq.SelectBuilder {
	if len(BelongGroup) == 0 {
		return sql
	}
	for _, groupConfig := range groupConfigs {
		if groupConfig.TsIndex == false {
			continue
		}
		bg, ok := BelongGroup[groupConfig.Value]
		if !ok {
			continue
		}
		if len(bg.IDs) > 0 {
			sql = sql.Where(stores.ArrayEqToSql("group_"+groupConfig.Value+"_ids", bg.IDs))
		}
		if len(bg.IDPaths) > 0 {
			sql = sql.Where(stores.ArrayEqToSql("group_"+groupConfig.Value+"_id_paths", bg.IDPaths))
		}
	}
	return sql
}

func GenTagsParams(ts string, groupConfigs []*deviceGroup.GroupDetail, BelongGroup map[string]def.IDsInfo) (tagNames string, vals string) {
	var tags []string
	var vv []string
	for _, groupConfig := range groupConfigs {
		if groupConfig.TsIndex == false {
			continue
		}
		tags = append(tags, fmt.Sprintf("`group_%s_ids`", groupConfig.Value), fmt.Sprintf("`group_%s_id_paths`", groupConfig.Value))
		vv = append(vv, "'"+utils.GenSliceStr(BelongGroup[groupConfig.Value].IDs)+"'", "'"+utils.GenSliceStr(BelongGroup[groupConfig.Value].IDPaths)+"'")
	}
	if len(tags) == 0 {
		return ts, ""
	}
	return ts + "," + strings.Join(tags, ","), "," + strings.Join(vv, ",")
}

func AlterTag(ctx context.Context, t *clients.Td, tables []string, tags map[string]any) error {
	for _, table := range tables {
		if table[0] != '`' {
			table = "`" + table + "`"
		}
		var vals []string
		for k, v := range tags {
			vals = append(vals, fmt.Sprintf(" `%s`='%v' ", k, v))
		}
		for i := 3; i > 0; i-- { //重试三次
			val := strings.Join(vals, ",")
			_, err := t.ExecContext(ctx, fmt.Sprintf(" ALTER TABLE %s SET TAG %s; ",
				table, val))
			if err != nil {
				if strings.Contains(err.Error(), "Table does not exist") {
					break
				}
				logx.WithContext(ctx).Error(err)
				continue
			}
			break
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
		var vals []string
		for k, v := range tag.Tags {
			vals = append(vals, fmt.Sprintf(" `%s`='%v' ", k, v))
		}
		for i := 3; i > 0; i-- { //重试三次
			val := strings.Join(vals, ",")
			if tag.Table[0] != '`' {
				tag.Table = "`" + tag.Table + "`"
			}
			_, err := t.ExecContext(ctx, fmt.Sprintf(" ALTER TABLE %s SET TAG %s; ",
				tag.Table, val))
			if err != nil {
				if strings.Contains(err.Error(), "Table does not exist") {
					break
				}
				logx.WithContext(ctx).Error(err)
				continue
			}
			break
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
