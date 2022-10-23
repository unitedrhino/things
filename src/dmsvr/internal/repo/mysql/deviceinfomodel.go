package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DeviceInfoModel = (*customDeviceInfoModel)(nil)

type (
	// DeviceInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeviceInfoModel.
	DeviceInfoModel interface {
		deviceInfoModel
		FindByFilter(ctx context.Context, filter DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error)
		CountByFilter(ctx context.Context, filter DeviceFilter) (size int64, err error)
	}

	customDeviceInfoModel struct {
		*defaultDeviceInfoModel
	}
	DeviceFilter struct {
		ProductID  string
		DeviceName string
		Tags       map[string]string
	}
)

// NewDeviceInfoModel returns a model for the database table.
func NewDeviceInfoModel(conn sqlx.SqlConn) DeviceInfoModel {
	return &customDeviceInfoModel{
		defaultDeviceInfoModel: newDeviceInfoModel(conn),
	}
}

func (d *DeviceFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if d.ProductID != "" {
		sql = sql.Where("`ProductID` = ?", d.ProductID)
	}
	if d.DeviceName != "" {
		sql = sql.Where("`DeviceName` like ?", "%"+d.DeviceName+"%")
	}
	if d.Tags != nil {
		for k, v := range d.Tags {
			sql = sql.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return sql
}

func (m *customDeviceInfoModel) FindByFilter(ctx context.Context, f DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	sql := sq.Select(deviceInfoRows).From(m.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.conn.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customDeviceInfoModel) CountByFilter(ctx context.Context, f DeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.table)
	sql = f.FmtSql(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = m.conn.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}
