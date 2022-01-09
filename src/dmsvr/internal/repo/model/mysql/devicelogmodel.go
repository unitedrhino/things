package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	deviceLogFieldNames          = builder.RawFieldNames(&DeviceLog{})
	deviceLogRows                = strings.Join(deviceLogFieldNames, ",")
	deviceLogRowsExpectAutoSet   = strings.Join(stringx.Remove(deviceLogFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	deviceLogRowsWithPlaceHolder = strings.Join(stringx.Remove(deviceLogFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	DeviceLogModel interface {
		Insert(data *DeviceLog) (sql.Result, error)
		FindOne(id int64) (*DeviceLog, error)
		Update(data *DeviceLog) error
		Delete(id int64) error
	}

	defaultDeviceLogModel struct {
		conn  sqlx.SqlConn
		table string
	}

	DeviceLog struct {
		Id          int64
		ProductID   string    // 产品id
		DeviceName  string    // 设备名称
		Payload     string    // 具体信息
		Topic       string    // 主题
		Action      string    // 操作类型
		Timestamp   time.Time // 操作时间
		CreatedTime time.Time
	}
)

func NewDeviceLogModel(conn sqlx.SqlConn) DeviceLogModel {
	return &defaultDeviceLogModel{
		conn:  conn,
		table: "`device_log`",
	}
}

func (m *defaultDeviceLogModel) Insert(data *DeviceLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, deviceLogRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.ProductID, data.DeviceName, data.Payload, data.Topic, data.Action, data.Timestamp, data.CreatedTime)
	return ret, err
}

func (m *defaultDeviceLogModel) FindOne(id int64) (*DeviceLog, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceLogRows, m.table)
	var resp DeviceLog
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDeviceLogModel) Update(data *DeviceLog) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, deviceLogRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.ProductID, data.DeviceName, data.Payload, data.Topic, data.Action, data.Timestamp, data.CreatedTime, data.Id)
	return err
}

func (m *defaultDeviceLogModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
