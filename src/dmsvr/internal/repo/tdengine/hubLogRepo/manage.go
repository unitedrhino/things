package hubLogRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsgManage"
)

func (d HubLogRepo) InitProduct(ctx context.Context, productID string) error {
	sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
		"(`ts` timestamp,`content` BINARY(5000),`topic` BINARY(500), `action` BINARY(100),"+
		" `request_id` BINARY(100), `trance_id` BINARY(100), `result_type` BIGINT)"+
		"TAGS (`product_id` BINARY(50),`device_name`  BINARY(50));",
		d.GetLogStableName())
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (d HubLogRepo) DropProduct(ctx context.Context, productID string) error {
	return nil
}

func (d HubLogRepo) DropDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", d.GetLogTableName(productID, deviceName))
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}
func (d HubLogRepo) Insert(ctx context.Context, data *deviceMsgManage.HubLog) error {
	sql := fmt.Sprintf("insert into %s using %s tags('%s','%s')(`ts`, `content`, `topic`, `action`,"+
		" `requestID`, `trance_id`, `result_type`) values (?,?,?,?,?,?,?);",
		d.GetLogTableName(data.ProductID, data.DeviceName), d.GetLogStableName(), data.ProductID, data.DeviceName)
	if _, err := d.t.ExecContext(ctx, sql, data.Timestamp, data.Content, data.Topic, data.Action,
		data.RequestID, data.TranceID, data.ResultType); err != nil {
		return err
	}
	return nil
}
