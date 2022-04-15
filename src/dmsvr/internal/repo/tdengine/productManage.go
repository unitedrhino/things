package tdengine

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func (d *DeviceDataRepo) DropProduct(ctx context.Context, t *deviceTemplate.Template, productID string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) InitProduct(ctx context.Context, t *deviceTemplate.Template, productID string) error {
	for _, p := range t.Properties {
		err := d.createPropertyStable(ctx, p, productID)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s|createPropertyStable|prodecutID:%v,properties:%v,err:%v",
				utils.FuncName(), productID, p, err)
			return err
		}
	}
	sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
		"(ts timestamp,event_id BINARY(50),event_type BINARY(20), param BINARY(5000)) TAGS (device_name BINARY(50));",
		getEventStableName(productID))
	if _, err := d.t.Exec(sql); err != nil {
		return err
	}
	sql = fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
		"(ts timestamp,action_id BINARY(50),input BINARY(5000),output BINARY(5000)) TAGS (device_name BINARY(50));", getActionStableName(productID))
	if _, err := d.t.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (d *DeviceDataRepo) ModifyProduct(
	ctx context.Context, oldT *deviceTemplate.Template, newt *deviceTemplate.Template, productID string) error {
	//todo 这部分的逻辑比较复杂,后面有大块时间再处理
	return nil
}

func (d *DeviceDataRepo) createPropertyStable(
	ctx context.Context, p deviceTemplate.Property, productID string) error {
	var sql string
	if p.Define.Type != deviceTemplate.STRUCT {
		sql = fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (ts timestamp,param %s) TAGS (device_name BINARY(50));",
			getPropertyStableName(productID, p.ID), getTdType(p.Define))
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	} else {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (ts timestamp, %s) TAGS (device_name BINARY(50));",
			getPropertyStableName(productID, p.ID), getSpecsColumn(p.Define.Specs))
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func getSpecsColumn(s deviceTemplate.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("%s %s", v.ID, getTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}
