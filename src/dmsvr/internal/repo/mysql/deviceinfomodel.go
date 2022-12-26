package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DeviceInfoModel = (*customDeviceInfoModel)(nil)

type (
	// DeviceInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeviceInfoModel.
	DeviceInfoModel interface {
		deviceInfoModel
		InsertDeviceInfo(ctx context.Context, data *DeviceInfo) error
		FindByFilter(ctx context.Context, filter DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error)
		CountByFilter(ctx context.Context, filter DeviceFilter) (size int64, err error)
		CountGroupByField(ctx context.Context, filter DeviceFilter, fieldName string) (map[string]int64, error)
		FindOneByProductIDAndDeviceName(ctx context.Context, productID string, deviceName string) (*DeviceInfo, error)
		UpdateDeviceInfo(ctx context.Context, data *DeviceInfo) error
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
		Range    int64
		Position string
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

	if d.Range > 0 {
		//d.Position 形如：point(116.393 39.905)
		sql = sql.Where(fmt.Sprintf("round(st_distance_sphere(ST_GeomFromText(%q), ST_GeomFromText(AsText(`position`))),2)>%d", d.Position, d.Range))
	}
	return sql
}

func (m *customDeviceInfoModel) FindByFilter(ctx context.Context, f DeviceFilter, page def.PageInfo) ([]*DeviceInfo, error) {
	var resp []*DeviceInfo
	sSql := "`id`,`productID`,`deviceName`,`secret`,`firstLogin`,`lastLogin`,`createdTime`,`updatedTime`,`deletedTime`,`version`,`logLevel`,`cert`,`isOnline`,`tags`,`address`, AsText(`position`) as position"
	sql := sq.Select(sSql).From(m.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (m *customDeviceInfoModel) InsertDeviceInfo(ctx context.Context, data *DeviceInfo) error {
	sql := fmt.Sprintf("INSERT INTO %s (`productID`,`deviceName`,`secret`,`version`,`logLevel`,`cert`,`isOnline`,`tags`,`address`,`position`) values (%q,%q,%q,%q,%d,%q,%d,%q,%q,%s)",
		m.table, data.ProductID, data.DeviceName, data.Secret, data.Version, data.LogLevel, data.Cert, data.IsOnline, data.Tags, data.Address, data.Position)
	_, err := m.conn.ExecCtx(ctx, sql)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (m *defaultDeviceInfoModel) FindOneByProductIDAndDeviceName(ctx context.Context, productID string, deviceName string) (*DeviceInfo, error) {
	var resp DeviceInfo
	query := fmt.Sprintf("select `id`,`productID`,`deviceName`,`secret`,`firstLogin`,`lastLogin`,`createdTime`,`updatedTime`,`deletedTime`,`version`,`logLevel`,`cert`,`isOnline`,`tags`,`address`, AsText(`position`) as position from %s where `productID` = ? and `deviceName` = ? limit 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, productID, deviceName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customDeviceInfoModel) UpdateDeviceInfo(ctx context.Context, data *DeviceInfo) error {
	values := fmt.Sprintf("productID=%q, deviceName=%q, secret=%q, version=%q, logLevel=%d,cert=%q, isOnline=%d, tags=%q, address=%q, position=%s",
		data.ProductID, data.DeviceName, data.Secret, data.Version, data.LogLevel, data.Cert, data.IsOnline, data.Tags, data.Address, data.Position)
	sql := fmt.Sprintf("update %s set %s where `id` = %d", m.table, values, data.Id)
	_, err := m.conn.ExecCtx(ctx, sql)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}
