package cache

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

// UserMultiDeviceShareManager 批量设备分享 Token 管理器
// 在 caches.Cache 基础上增加 Redis Set 索引，支持按用户或设备定位分享 Token
type UserMultiDeviceShareManager struct {
	dataCache *caches.Cache[dm.UserDeviceShareMultiInfo, string]
	store     kv.Store
}

// MultiShareItem 批量分享列表项，包含 Token 及其对应数据
type MultiShareItem struct {
	Token string                       // 分享 Token
	Info  *dm.UserDeviceShareMultiInfo // 分享数据
}

// multiShareContainsDevice 判断分享 Token 是否包含指定设备。
func multiShareContainsDevice(info *dm.UserDeviceShareMultiInfo, productID string, deviceName string) bool {
	if info == nil {
		return false
	}
	for _, device := range info.Devices {
		if device.ProductID == productID && device.DeviceName == deviceName {
			return true
		}
	}
	return false
}

// NewUserMultiDeviceShareManager 创建批量分享 Token 管理器
func NewUserMultiDeviceShareManager(dataCache *caches.Cache[dm.UserDeviceShareMultiInfo, string], store kv.Store) *UserMultiDeviceShareManager {
	return &UserMultiDeviceShareManager{
		dataCache: dataCache,
		store:     store,
	}
}

// genListKey 生成用户维度的 Token 列表索引 key
func (m *UserMultiDeviceShareManager) genListKey(tenantCode string, userID int64) string {
	return fmt.Sprintf("things:device:share:batch:list:%s:%d", tenantCode, userID)
}

// genDeviceListKey 生成设备维度的 Token 索引键。
func (m *UserMultiDeviceShareManager) genDeviceListKey(tenantCode string, productID string, deviceName string) string {
	return fmt.Sprintf("things:device:share:batch:device:%s:%s:%s", tenantCode, productID, deviceName)
}

// multiShareIndexKeys 返回 Token 需要写入的全部索引键。
func (m *UserMultiDeviceShareManager) multiShareIndexKeys(tenantCode string, info *dm.UserDeviceShareMultiInfo) []string {
	if info == nil {
		return nil
	}
	keys := []string{m.genListKey(tenantCode, info.UserID)}
	seen := make(map[string]struct{}, len(info.Devices))
	for _, device := range info.Devices {
		key := m.genDeviceListKey(tenantCode, device.ProductID, device.DeviceName)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	return keys
}

// SetData 写入分享数据，同时把 Token 加入用户和设备索引。
func (m *UserMultiDeviceShareManager) SetData(ctx context.Context, tenantCode, token string, data *dm.UserDeviceShareMultiInfo) error {
	err := m.dataCache.SetData(ctx, token, data)
	if err != nil {
		return err
	}
	if data != nil {
		for _, listKey := range m.multiShareIndexKeys(tenantCode, data) {
			_, err = m.store.SaddCtx(ctx, listKey, token)
			if err != nil {
				return stores.ErrFmt(err)
			}
			// 给 Set 设置与数据相同的 TTL，避免长期残留
			err = m.store.ExpireCtx(ctx, listKey, int(userShared.MultiDeviceShareTokenTTLSeconds))
			if err != nil {
				return stores.ErrFmt(err)
			}
		}
	}
	return nil
}

// GetData 通过 Token 获取分享数据
func (m *UserMultiDeviceShareManager) GetData(ctx context.Context, token string) (*dm.UserDeviceShareMultiInfo, error) {
	return m.dataCache.GetData(ctx, token)
}

// GetList 获取指定用户的批量分享 Token 列表，自动清理过期项
func (m *UserMultiDeviceShareManager) GetList(ctx context.Context, tenantCode string, userID int64) ([]*MultiShareItem, error) {
	listKey := m.genListKey(tenantCode, userID)
	tokens, err := m.store.SmembersCtx(ctx, listKey)
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	var result []*MultiShareItem
	for _, token := range tokens {
		info, err := m.GetData(ctx, token)
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				// 已过期或不存在，从列表中清理
				_, _ = m.store.SremCtx(ctx, listKey, token)
			}
			continue
		}
		result = append(result, &MultiShareItem{Token: token, Info: info})
	}
	return result, nil
}

// DeleteToken 删除指定 Token，并从用户和设备索引中移除。
func (m *UserMultiDeviceShareManager) DeleteToken(ctx context.Context, tenantCode string, userID int64, token string) error {
	info, infoErr := m.GetData(ctx, token)
	if infoErr != nil && !errors.Cmp(infoErr, errors.NotFind) {
		return infoErr
	}
	if err := m.dataCache.SetData(ctx, token, nil); err != nil {
		return err
	}
	indexKeys := []string{m.genListKey(tenantCode, userID)}
	if info != nil {
		indexKeys = m.multiShareIndexKeys(tenantCode, info)
	}
	for _, listKey := range indexKeys {
		if _, err := m.store.SremCtx(ctx, listKey, token); err != nil {
			return stores.ErrFmt(err)
		}
	}
	return nil
}

// DeleteAllDeviceTokens 删除所有包含指定设备的未过期分享 Token。
func (m *UserMultiDeviceShareManager) DeleteAllDeviceTokens(
	ctx context.Context,
	tenantCode string,
	productID string,
	deviceName string,
) error {
	listKey := m.genDeviceListKey(tenantCode, productID, deviceName)
	tokens, err := m.store.SmembersCtx(ctx, listKey)
	if err != nil {
		return stores.ErrFmt(err)
	}
	for _, token := range tokens {
		info, getErr := m.GetData(ctx, token)
		if getErr != nil {
			if errors.Cmp(getErr, errors.NotFind) {
				_, _ = m.store.SremCtx(ctx, listKey, token)
				continue
			}
			return getErr
		}
		if !multiShareContainsDevice(info, productID, deviceName) {
			_, _ = m.store.SremCtx(ctx, listKey, token)
			continue
		}
		if err = m.DeleteToken(ctx, tenantCode, info.UserID, token); err != nil {
			return err
		}
	}
	return nil
}

// DeleteDeviceTokens 删除指定用户创建且包含目标设备的全部未过期分享 Token。
func (m *UserMultiDeviceShareManager) DeleteDeviceTokens(
	ctx context.Context,
	tenantCode string,
	userID int64,
	productID string,
	deviceName string,
) error {
	items, err := m.GetList(ctx, tenantCode, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if !multiShareContainsDevice(item.Info, productID, deviceName) {
			continue
		}
		if err = m.DeleteToken(ctx, tenantCode, userID, item.Token); err != nil {
			return err
		}
	}
	return nil
}
