package hubLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/devices"
)

func (h *HubLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	h.once.Do(func() {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`content` BINARY(5000),`topic` BINARY(500), `action` BINARY(100),"+
			" `request_id` BINARY(100), `trace_id` BINARY(100), `result_type` BIGINT,`resp_payload` BINARY(5000))"+
			"TAGS (`product_id` BINARY(50),`device_name`  BINARY(50));",
			h.GetLogStableName())
		_, err = h.t.ExecContext(ctx, sql)
	})
	return
}

func (h *HubLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	return nil
}

func (h *HubLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", h.GetLogTableName(productID, deviceName))
	if _, err := h.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (h *HubLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	return nil
}
