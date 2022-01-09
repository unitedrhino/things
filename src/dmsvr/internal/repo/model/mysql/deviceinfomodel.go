package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	deviceInfoFieldNames          = builder.RawFieldNames(&DeviceInfo{})
	deviceInfoRows                = strings.Join(deviceInfoFieldNames, ",")
	deviceInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(deviceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	deviceInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(deviceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDmDeviceInfoIdPrefix                  = "cache:dm:deviceInfo:id:"
	cacheDmDeviceInfoProductIDDeviceNamePrefix = "cache:dm:deviceInfo:productID:deviceName:"
)

type (
	DeviceInfoModel interface {
		Insert(data *DeviceInfo) (sql.Result, error)
		FindOne(id int64) (*DeviceInfo, error)
		FindOneByProductIDDeviceName(productID string, deviceName string) (*DeviceInfo, error)
		Update(data *DeviceInfo) error
		Delete(id int64) error
	}

	defaultDeviceInfoModel struct {
		sqlc.CachedConn
		table string
	}

	DeviceInfo struct {
		Id          int64
		ProductID   string       // 产品id
		DeviceName  string       // 设备名称
		Secret      string       // 设备秘钥
		FirstLogin  sql.NullTime // 激活时间
		LastLogin   sql.NullTime // 最后上线时间
		CreatedTime time.Time
		UpdatedTime sql.NullTime
		DeletedTime sql.NullTime
		Version     string // 固件版本
		LogLevel    int64  // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试
		Cert        string // 设备证书
	}
)

func NewDeviceInfoModel(conn sqlx.SqlConn, c cache.CacheConf) DeviceInfoModel {
	return &defaultDeviceInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`device_info`",
	}
}

func (m *defaultDeviceInfoModel) Insert(data *DeviceInfo) (sql.Result, error) {
	dmDeviceInfoIdKey := fmt.Sprintf("%s%v", cacheDmDeviceInfoIdPrefix, data.Id)
	dmDeviceInfoProductIDDeviceNameKey := fmt.Sprintf("%s%v:%v", cacheDmDeviceInfoProductIDDeviceNamePrefix, data.ProductID, data.DeviceName)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, deviceInfoRowsExpectAutoSet)
		return conn.Exec(query, data.ProductID, data.DeviceName, data.Secret, data.FirstLogin, data.LastLogin, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Version, data.LogLevel, data.Cert)
	}, dmDeviceInfoIdKey, dmDeviceInfoProductIDDeviceNameKey)
	return ret, err
}

func (m *defaultDeviceInfoModel) FindOne(id int64) (*DeviceInfo, error) {
	dmDeviceInfoIdKey := fmt.Sprintf("%s%v", cacheDmDeviceInfoIdPrefix, id)
	var resp DeviceInfo
	err := m.QueryRow(&resp, dmDeviceInfoIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceInfoRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDeviceInfoModel) FindOneByProductIDDeviceName(productID string, deviceName string) (*DeviceInfo, error) {
	dmDeviceInfoProductIDDeviceNameKey := fmt.Sprintf("%s%v:%v", cacheDmDeviceInfoProductIDDeviceNamePrefix, productID, deviceName)
	var resp DeviceInfo
	err := m.QueryRowIndex(&resp, dmDeviceInfoProductIDDeviceNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `productID` = ? and `deviceName` = ? limit 1", deviceInfoRows, m.table)
		if err := conn.QueryRow(&resp, query, productID, deviceName); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDeviceInfoModel) Update(data *DeviceInfo) error {
	dmDeviceInfoIdKey := fmt.Sprintf("%s%v", cacheDmDeviceInfoIdPrefix, data.Id)
	dmDeviceInfoProductIDDeviceNameKey := fmt.Sprintf("%s%v:%v", cacheDmDeviceInfoProductIDDeviceNamePrefix, data.ProductID, data.DeviceName)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, deviceInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.ProductID, data.DeviceName, data.Secret, data.FirstLogin, data.LastLogin, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Version, data.LogLevel, data.Cert, data.Id)
	}, dmDeviceInfoProductIDDeviceNameKey, dmDeviceInfoIdKey)
	return err
}

func (m *defaultDeviceInfoModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	dmDeviceInfoIdKey := fmt.Sprintf("%s%v", cacheDmDeviceInfoIdPrefix, id)
	dmDeviceInfoProductIDDeviceNameKey := fmt.Sprintf("%s%v:%v", cacheDmDeviceInfoProductIDDeviceNamePrefix, data.ProductID, data.DeviceName)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, dmDeviceInfoIdKey, dmDeviceInfoProductIDDeviceNameKey)
	return err
}

func (m *defaultDeviceInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDmDeviceInfoIdPrefix, primary)
}

func (m *defaultDeviceInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
