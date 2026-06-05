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
// 在 caches.Cache 基础上增加 Redis Set 索引，支持按用户列出分享 Token
type UserMultiDeviceShareManager struct {
	dataCache *caches.Cache[dm.UserDeviceShareMultiInfo, string]
	store     kv.Store
}

// MultiShareItem 批量分享列表项，包含 Token 及其对应数据
type MultiShareItem struct {
	Token string                       // 分享 Token
	Info  *dm.UserDeviceShareMultiInfo // 分享数据
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

// SetData 写入分享数据，同时把 Token 加入用户列表索引
func (m *UserMultiDeviceShareManager) SetData(ctx context.Context, tenantCode, token string, data *dm.UserDeviceShareMultiInfo) error {
	err := m.dataCache.SetData(ctx, token, data)
	if err != nil {
		return err
	}
	if data != nil {
		listKey := m.genListKey(tenantCode, data.UserID)
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

// DeleteToken 删除指定 Token，并从用户列表索引中移除
func (m *UserMultiDeviceShareManager) DeleteToken(ctx context.Context, tenantCode string, userID int64, token string) error {
	// 先删除数据缓存
	err := m.dataCache.SetData(ctx, token, nil)
	if err != nil {
		return err
	}
	// 再从用户列表中移除
	_, err = m.store.SremCtx(ctx, m.genListKey(tenantCode, userID), token)
	if err != nil {
		return stores.ErrFmt(err)
	}
	return nil
}
