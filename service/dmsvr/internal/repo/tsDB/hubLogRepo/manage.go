package hubLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
)

func (h *HubLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	return
}

func (h *HubLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	err := h.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&Hub{}).Error
	return stores.ErrFmt(err)
}

func (h *HubLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	err := h.db.WithContext(ctx).Where("product_id = ?", productID).
		Where("device_name = ?", deviceName).Delete(&Hub{}).Error
	return stores.ErrFmt(err)
}

func (h *HubLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	return nil
}
