package schemaDataRepo

import (
	"context"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type latestAggPartitionConfig struct {
	BaseSelects  []string
	FinalSelects []string
	Groups       []string
}

func (d *DeviceDataRepo) getPropertyLatestAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterLatestAggOpt) ([]*msgThing.PropertyLatestData, error) {
	if len(filter.Aggs) == 0 {
		return nil, errors.Parameter.AddMsg("至少填写一个聚合函数")
	}

	var retMap = map[string]msgThing.PropertyLatestData{}
	for _, agg := range filter.Aggs {
		p, err := d.getPropertyByAggDataID(m, agg.DataID)
		if err != nil {
			return nil, err
		}

		db, dataType, err := d.getPropertyLatestAggQuery(ctx, p, agg, filter)
		if err != nil {
			logx.WithContext(ctx).Errorf(err.Error())
			return nil, err
		}

		var rows []map[string]any
		if err = db.Find(&rows).Error; err != nil {
			logx.WithContext(ctx).Error(err)
			return nil, stores.ErrFmt(err)
		}

		retMap = d.toPropertyLatestData(agg, &dataType, rows, retMap)
	}

	return utils.MapVToSlice2(retMap), nil
}

func (d *DeviceDataRepo) getPropertyByAggDataID(m *schema.Model, dataID string) (*schema.Property, error) {
	if p, ok := m.Property[dataID]; ok {
		return p, nil
	}
	id, _, ok := schema.GetArray(dataID)
	if !ok {
		return nil, errors.Parameter.AddMsg("标识符未找到")
	}
	p, ok := m.Property[id]
	if !ok {
		return nil, errors.Parameter.AddMsg("标识符未找到")
	}
	return p, nil
}

func (d *DeviceDataRepo) getPropertyLatestAggQuery(
	ctx context.Context, p *schema.Property, agg msgThing.PropertyAgg, filter msgThing.FilterLatestAggOpt,
) (*stores.DB, schema.Define, error) {
	baseDataID, pos, hasArrayPos := schema.GetArray(agg.DataID)
	if baseDataID == "" {
		baseDataID = agg.DataID
	}

	tableName := getTableName(p.Define)
	dataType := p.Define
	if dataType.Type == schema.DataTypeArray {
		dataType = *dataType.ArrayInfo
	}
	if dataType.Type == schema.DataTypeStruct {
		dd, _ := schema.ParseDataID(agg.DataID)
		if dd != nil && dd.Column != "" {
			if spec := dataType.Spec[dd.Column]; spec != nil {
				dataType = spec.DataType
			}
		}
	}

	baseQuery := d.db.WithContext(ctx).Table(tableName + " as tb")
	baseQuery, partitionCfg := d.applyLatestAggPartitions(filter.PartitionBy, baseQuery)

	selectParts := []string{
		"tb.product_id as product_id",
		"tb.device_name as device_name",
		"tb.ts as ts",
		"tb.param as param",
		"ROW_NUMBER() OVER (PARTITION BY tb.product_id, tb.device_name ORDER BY tb.ts DESC) as rn",
	}
	selectParts = append(selectParts, partitionCfg.BaseSelects...)
	baseQuery = baseQuery.Select(strings.Join(distinctStrings(selectParts), ", "))
	baseQuery = baseQuery.Where("tb.identifier = ?", baseDataID)
	if hasArrayPos {
		baseQuery = baseQuery.Where("tb.pos = ?", pos)
	}
	baseQuery = d.fillFilter(ctx, baseQuery, filter.Filter)

	latestRows := d.db.WithContext(ctx).Table("(?) as latest", baseQuery).Where("rn = 1")

	finalSelects := append([]string{}, partitionCfg.FinalSelects...)
	for _, argFunc := range agg.ArgFuncs {
		finalSelects = append(finalSelects, d.buildLatestAggExpr(argFunc, dataType))
	}
	finalSelects = distinctStrings(finalSelects)

	finalQuery := latestRows.Select(strings.Join(finalSelects, ", "))
	if len(partitionCfg.Groups) > 0 {
		finalQuery = finalQuery.Group(strings.Join(distinctStrings(partitionCfg.Groups), ", "))
	}
	return finalQuery, dataType, nil
}

func (d *DeviceDataRepo) applyLatestAggPartitions(partitionBy string, db *stores.DB) (*stores.DB, latestAggPartitionConfig) {
	if partitionBy == "" {
		return db, latestAggPartitionConfig{}
	}

	cfg := latestAggPartitionConfig{}
	parts := strings.Split(utils.CamelCaseToUdnderscore(partitionBy), ",")
	var hasDeviceJoin bool
	var hasGroupJoin bool

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case part == "device_name":
			cfg.FinalSelects = append(cfg.FinalSelects, "latest.device_name as device_name")
			cfg.Groups = append(cfg.Groups, "latest.device_name")
		case strings.Contains(part, "area_id"):
			if !hasDeviceJoin {
				db = db.Joins("left join dm_device_info as di on tb.product_id = di.product_id and tb.device_name = di.device_name")
				hasDeviceJoin = true
			}
			cfg.BaseSelects = append(cfg.BaseSelects,
				"di.tenant_code as tenant_code",
				"di.project_id as project_id",
				"di.area_id as area_id",
				"di.area_id_path as area_id_path",
			)
			cfg.FinalSelects = append(cfg.FinalSelects,
				"latest.tenant_code as tenant_code",
				"latest.project_id as project_id",
				"latest.area_id as area_id",
				"latest.area_id_path as area_id_path",
			)
			cfg.Groups = append(cfg.Groups,
				"latest.tenant_code",
				"latest.project_id",
				"latest.area_id",
				"latest.area_id_path",
			)
		case strings.Contains(part, "group"):
			if !hasGroupJoin {
				db = db.Joins("left join dm_group_device as gd on tb.product_id = gd.product_id and tb.device_name = gd.device_name")
				hasGroupJoin = true
			}
			cfg.BaseSelects = append(cfg.BaseSelects,
				"gd.purpose as group_purpose",
				"gd.group_id as group_id",
				"gd.group_id_path as group_id_path",
			)
			cfg.FinalSelects = append(cfg.FinalSelects,
				"latest.group_purpose as group_purpose",
				"latest.group_id as group_id",
				"latest.group_id_path as group_id_path",
			)
			cfg.Groups = append(cfg.Groups,
				"latest.group_purpose",
				"latest.group_id",
				"latest.group_id_path",
			)
		}
	}
	return db, cfg
}

func (d *DeviceDataRepo) buildLatestAggExpr(argFunc string, dataType schema.Define) string {
	paramExpr := "latest.param"
	switch argFunc {
	case "first":
		return fmt.Sprintf("(ARRAY_AGG(%s ORDER BY latest.ts ASC))[1] AS first_param", paramExpr)
	case "last":
		return fmt.Sprintf("(ARRAY_AGG(%s ORDER BY latest.ts DESC))[1] AS last_param", paramExpr)
	case "max":
		return fmt.Sprintf("%s AS max_param", d.wrapLatestAggParam("MAX", dataType))
	case "min":
		return fmt.Sprintf("%s AS min_param", d.wrapLatestAggParam("MIN", dataType))
	case "avg":
		return fmt.Sprintf("%s AS avg_param", d.wrapLatestAggParam("AVG", dataType))
	case "sum":
		return fmt.Sprintf("%s AS sum_param", d.wrapLatestAggParam("SUM", dataType))
	case "count":
		return "COUNT(latest.param) AS count_param"
	default:
		return fmt.Sprintf("%s(%s) AS %s_param", strings.ToUpper(argFunc), paramExpr, argFunc)
	}
}

func (d *DeviceDataRepo) wrapLatestAggParam(fn string, dataType schema.Define) string {
	switch dataType.Type {
	case schema.DataTypeBool:
		return fmt.Sprintf("%s(CASE WHEN latest.param THEN 1 ELSE 0 END)", fn)
	default:
		return fmt.Sprintf("%s(latest.param)", fn)
	}
}

func (d *DeviceDataRepo) toPropertyLatestData(
	agg msgThing.PropertyAgg, define *schema.Define, rows []map[string]any, retMap map[string]msgThing.PropertyLatestData,
) map[string]msgThing.PropertyLatestData {
	dd, _ := schema.ParseDataID(agg.DataID)
	for _, row := range rows {
		data := msgThing.PropertyLatestData{
			DeviceName: cast.ToString(row["device_name"]),
			TenantCode: dataType2TenantCode(row["tenant_code"]),
			ProjectID:  dataType.ProjectID(cast.ToInt64(row["project_id"])),
			AreaID:     dataType.AreaID(cast.ToInt64(row["area_id"])),
			AreaIDPath: dataType.AreaIDPath(cast.ToString(row["area_id_path"])),
		}
		if cast.ToString(row["group_purpose"]) != "" {
			groupID := cast.ToInt64(row["group_id"])
			groupIDPath := cast.ToString(row["group_id_path"])
			info := def.IDsInfo{}
			if groupID != 0 {
				info.IDs = []int64{groupID}
			}
			if groupIDPath != "" {
				info.IDPaths = []string{groupIDPath}
			}
			data.BelongGroup = map[string]def.IDsInfo{
				cast.ToString(row["group_purpose"]): info,
			}
		}

		key := utils.MarshalNoErr(data)
		ret := retMap[key]
		if ret.DeviceName == "" && ret.ProjectID == 0 && ret.AreaID == 0 && len(ret.Values) == 0 {
			ret = data
		}

		value := msgThing.PropertyLatestAggData{
			Identifier: agg.DataID,
			Values:     map[string]any{},
		}
		if define.Type == schema.DataTypeStruct && dd != nil && dd.Column == "" {
			argCache := map[string]map[string]any{}
			for k, v := range row {
				if !strings.HasSuffix(k, "_param") {
					continue
				}
				keys := strings.Split(k, "_")
				if len(keys) < 3 {
					continue
				}
				argFunc := keys[len(keys)-2]
				if argCache[argFunc] == nil {
					argCache[argFunc] = map[string]any{}
				}
				fieldID := strings.Join(keys[:len(keys)-2], "_")
				spec := define.Spec[fieldID]
				if spec == nil {
					continue
				}
				if formatted, err := spec.DataType.FmtValue(v); err == nil {
					v = formatted
				}
				if b, ok := v.(bool); ok {
					v = cast.ToInt64(b)
				}
				argCache[argFunc][fieldID] = v
			}
			for argFunc, values := range argCache {
				value.Values[argFunc] = values
			}
		} else {
			for k, v := range row {
				if !strings.HasSuffix(k, "_param") {
					continue
				}
				argFunc := strings.TrimSuffix(k, "_param")
				if formatted, err := define.FmtValue(v); err == nil {
					v = formatted
				}
				if b, ok := v.(bool); ok {
					v = cast.ToInt64(b)
				}
				value.Values[argFunc] = v
			}
		}

		ret.Values = append(ret.Values, value)
		retMap[key] = ret
	}
	return retMap
}

func dataType2TenantCode(v any) dataType.TenantCode {
	return dataType.TenantCode(cast.ToString(v))
}

func distinctStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
