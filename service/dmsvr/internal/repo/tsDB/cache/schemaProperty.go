package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

// GenRedisPropertyFirstKey 生成设备属性首次记录的Redis键
func GenRedisPropertyFirstKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:first:%s:%s", productID, deviceName)
}

// GenRedisPropertyLastKey 生成设备属性最后记录的Redis键
func GenRedisPropertyLastKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:last:%s:%s", productID, deviceName)
}

// PropertyCacheManager 属性缓存管理器
type PropertyCacheManager struct {
	kv kv.Store
}

// NewPropertyCacheManager 创建属性缓存管理器
func NewPropertyCacheManager(kv kv.Store) *PropertyCacheManager {
	return &PropertyCacheManager{
		kv: kv,
	}
}

// UpdatePropertyCache 更新属性缓存
func (p *PropertyCacheManager) UpdatePropertyCache(ctx context.Context, productID, deviceName string,
	property *schema.Property, data map[string]any, timestamp time.Time) error {

	log := logx.WithContext(ctx)

	for identifier, value := range data {
		propertyData := msgThing.PropertyData{
			Identifier: identifier,
			Param:      value,
			TimeStamp:  timestamp,
		}
		propertyData.Fmt()

		// 更新最后记录
		err := p.kv.Hset(GenRedisPropertyLastKey(productID, deviceName), identifier, propertyData.String())
		if err != nil {
			log.Errorf("更新属性最后记录失败: %v", err)
			return err
		}

		// 检查是否需要更新首次记录
		err = p.updateFirstRecord(ctx, productID, deviceName, identifier, property, value, propertyData)
		if err != nil {
			log.Errorf("更新属性首次记录失败: %v", err)
			return err
		}
	}

	return nil
}

// updateFirstRecord 更新首次记录
func (p *PropertyCacheManager) updateFirstRecord(ctx context.Context, productID, deviceName, identifier string,
	property *schema.Property, value any, propertyData msgThing.PropertyData) error {

	log := logx.WithContext(ctx)

	// 获取现有的首次记录
	retStr, err := p.kv.Hget(GenRedisPropertyFirstKey(productID, deviceName), identifier)
	if err != nil && !errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
		log.Errorf("获取属性首次记录失败: %v", err)
		return err
	}

	// 如果存在记录，检查值是否相等
	if retStr != "" {
		var ret msgThing.PropertyData
		err = json.Unmarshal([]byte(retStr), &ret)
		if err != nil {
			log.Errorf("解析属性首次记录失败: %v", err)
			return err
		}

		// 如果值相等，不更新首次记录
		if msgThing.IsParamValEq(&property.Define, value, ret.Param) {
			return nil
		}
	}

	// 更新首次记录
	err = p.kv.Hset(GenRedisPropertyFirstKey(productID, deviceName), identifier, propertyData.String())
	if err != nil {
		log.Errorf("设置属性首次记录失败: %v", err)
		return err
	}

	return nil
}

// GetPropertyFirstRecord 获取属性首次记录
func (p *PropertyCacheManager) GetPropertyFirstRecord(ctx context.Context, productID, deviceName, identifier string) (*msgThing.PropertyData, error) {
	retStr, err := p.kv.Hget(GenRedisPropertyFirstKey(productID, deviceName), identifier)
	if err != nil {
		if errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	if retStr == "" {
		return nil, nil
	}

	var data msgThing.PropertyData
	err = json.Unmarshal([]byte(retStr), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetPropertyLastRecord 获取属性最后记录
func (p *PropertyCacheManager) GetPropertyLastRecord(ctx context.Context, productID, deviceName, identifier string) (*msgThing.PropertyData, error) {
	retStr, err := p.kv.Hget(GenRedisPropertyLastKey(productID, deviceName), identifier)
	if err != nil {
		if errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
			return nil, nil // 没有找到记录
		}
		return nil, err
	}

	if retStr == "" {
		return nil, nil
	}

	var data msgThing.PropertyData
	err = json.Unmarshal([]byte(retStr), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetPropertyAllLastRecord 获取设备所有属性的最后记录
func (p *PropertyCacheManager) GetPropertyAllLastRecord(ctx context.Context, productID, deviceName string) ([]*msgThing.PropertyData, error) {
	// 获取所有最后记录
	lastRecords, err := p.kv.Hgetall(GenRedisPropertyLastKey(productID, deviceName))
	if err != nil {
		return nil, err
	}

	var result []*msgThing.PropertyData
	for _, dataStr := range lastRecords {
		var data msgThing.PropertyData
		err = json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			logx.WithContext(ctx).Errorf("解析属性记录失败: %v", err)
			continue
		}
		result = append(result, &data)
	}

	return result, nil
}

// ClearPropertyCache 清除设备属性缓存
func (p *PropertyCacheManager) ClearPropertyCache(ctx context.Context, productID, deviceName string) error {
	_, err := p.kv.DelCtx(ctx, GenRedisPropertyLastKey(productID, deviceName), GenRedisPropertyFirstKey(productID, deviceName))
	if err != nil {
		logx.WithContext(ctx).Errorf("清除设备属性缓存失败: %v", err)
		return err
	}
	return nil
}

// CheckIsChange 检查属性值是否发生变化
func (p *PropertyCacheManager) CheckIsChange(ctx context.Context, dev devices.Core, property *schema.Property, data msgThing.PropertyData) bool {
	if property.RecordMode == schema.RecordModeAll || property.RecordMode == 0 {
		return true
	}
	if property.RecordMode == schema.RecordModeNone {
		return false
	}

	data.Fmt()
	retStr, err := p.kv.Hget(GenRedisPropertyLastKey(dev.ProductID, dev.DeviceName), data.Identifier)
	if err != nil || retStr == "" {
		return true
	}

	var ret msgThing.PropertyData
	err = json.Unmarshal([]byte(retStr), &ret)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return true
	}

	// 如果值相等，不记录
	if msgThing.IsParamValEq(&property.Define, data.Param, ret.Param) {
		return false
	}

	return true
}

// GetAllPropertyRecords 获取设备所有属性记录
func (p *PropertyCacheManager) GetAllPropertyRecords(ctx context.Context, productID, deviceName string) (map[string]*msgThing.PropertyData, error) {
	// 获取最后记录
	lastRecords, err := p.kv.Hgetall(GenRedisPropertyLastKey(productID, deviceName))
	if err != nil {
		return nil, err
	}

	// 获取首次记录
	firstRecords, err := p.kv.Hgetall(GenRedisPropertyFirstKey(productID, deviceName))
	if err != nil {
		return nil, err
	}

	result := make(map[string]*msgThing.PropertyData)

	// 处理最后记录
	for identifier, dataStr := range lastRecords {
		var data msgThing.PropertyData
		err = json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			logx.WithContext(ctx).Errorf("解析属性记录失败: %v", err)
			continue
		}
		result[identifier+"_last"] = &data
	}

	// 处理首次记录
	for identifier, dataStr := range firstRecords {
		var data msgThing.PropertyData
		err = json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			logx.WithContext(ctx).Errorf("解析属性记录失败: %v", err)
			continue
		}
		result[identifier+"_first"] = &data
	}

	return result, nil
}
