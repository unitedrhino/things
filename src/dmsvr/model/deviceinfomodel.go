package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	deviceInfoFieldNames          = builderx.RawFieldNames(&DeviceInfo{})
	deviceInfoRows                = strings.Join(deviceInfoFieldNames, ",")
	deviceInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(deviceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	deviceInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(deviceInfoFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDeviceInfoIdPrefix = "cache#DeviceInfo#id#"
)

type (
	DeviceInfoModel interface {
		Insert(data DeviceInfo) (sql.Result, error)
		FindOne(id int64) (*DeviceInfo, error)
		Update(data DeviceInfo) error
		Delete(id int64) error
	}

	defaultDeviceInfoModel struct {
		sqlc.CachedConn
		table string
	}

	DeviceInfo struct {
		Id          int64        `db:"id"`
		ProductID   string       `db:"productID"`  // 产品id
		DeviceName  string       `db:"deviceName"` // 设备名称
		Secret      string       `db:"secret"`     // 设备秘钥
		FirstLogin  sql.NullTime `db:"firstLogin"` // 激活时间
		LastLogin   sql.NullTime `db:"lastLogin"`  // 最后上线时间
		CreatedTime time.Time    `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
		Version     string       `db:"version"`  // 固件版本
		LogLevel    int64        `db:"logLevel"` // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试
		Cert        string       `db:"cert"`     // 设备证书
		Template    string       `db:"template"` // 数据模板
	}
)

func NewDeviceInfoModel(conn sqlx.SqlConn, c cache.CacheConf) DeviceInfoModel {
	return &defaultDeviceInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`device_info`",
	}
}

func (m *defaultDeviceInfoModel) Insert(data DeviceInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, deviceInfoRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.ProductID, data.DeviceName, data.Secret, data.FirstLogin, data.LastLogin, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Version, data.LogLevel, data.Cert, data.Template)

	return ret, err
}

func (m *defaultDeviceInfoModel) FindOne(id int64) (*DeviceInfo, error) {
	deviceInfoIdKey := fmt.Sprintf("%s%v", cacheDeviceInfoIdPrefix, id)
	var resp DeviceInfo
	err := m.QueryRow(&resp, deviceInfoIdKey, func(conn sqlx.SqlConn, v interface{}) error {
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

func (m *defaultDeviceInfoModel) Update(data DeviceInfo) error {
	deviceInfoIdKey := fmt.Sprintf("%s%v", cacheDeviceInfoIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, deviceInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.ProductID, data.DeviceName, data.Secret, data.FirstLogin, data.LastLogin, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Version, data.LogLevel, data.Cert, data.Template, data.Id)
	}, deviceInfoIdKey)
	return err
}

func (m *defaultDeviceInfoModel) Delete(id int64) error {

	deviceInfoIdKey := fmt.Sprintf("%s%v", cacheDeviceInfoIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, deviceInfoIdKey)
	return err
}

func (m *defaultDeviceInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDeviceInfoIdPrefix, primary)
}

func (m *defaultDeviceInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
