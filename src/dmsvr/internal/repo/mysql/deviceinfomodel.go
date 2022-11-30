package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
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
		CountGroupByField(ctx context.Context, filter DeviceFilter, fieldName string) (map[string]int64, error)
	}

	customDeviceInfoModel struct {
		*defaultDeviceInfoModel
	}
	DeviceFilter struct {
		ProductID     string
		DeviceName    string
		Tags          map[string]string
		LastLoginTime struct {
			Start int64
			End   int64
		}
		IsOnline []int64
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
	if d.LastLoginTime.Start != 0 {
		sql = sql.Where("lastLogin >= ?", utils.ToYYMMddHHSS(d.LastLoginTime.Start*1000))
	}
	if d.LastLoginTime.End != 0 {
		sql = sql.Where("lastLogin <= ?", utils.ToYYMMddHHSS(d.LastLoginTime.End*1000))
	}
	if len(d.IsOnline) != 0 {
		sql = sql.Where(fmt.Sprintf("isOnline in (%v)", store.ArrayToSql(d.IsOnline)))
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

func (m *customDeviceInfoModel) CountGroupByField(ctx context.Context, f DeviceFilter, columnName string) (map[string]int64, error) {
	sql := sq.Select(fmt.Sprintf("%s as CountKey", columnName), "count(1) as count").From(m.table)
	sql = f.FmtSql(sql)
	sql = sql.GroupBy(columnName)
	query, arg, err := sql.ToSql()
	result := make(map[string]int64, 0)

	type countModel struct {
		CountKey string
		Count    int64
	}
	countModelList := make([]*countModel, 0)

	err = m.conn.QueryRowsCtx(ctx, &countModelList, query, arg...)
	if err != nil {
		return result, err
	}

	for _, v := range countModelList {
		result[v.CountKey] = v.Count
	}

	return result, err
}
