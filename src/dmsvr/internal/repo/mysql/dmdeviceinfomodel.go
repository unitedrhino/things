package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ DmDeviceInfoModel = (*customDmDeviceInfoModel)(nil)

type (
	// DmDeviceInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmDeviceInfoModel.
	DmDeviceInfoModel interface {
		dmDeviceInfoModel
		InsertDeviceInfo(ctx context.Context, data *DmDeviceInfo) error
		FindByFilter(ctx context.Context, filter DeviceFilter, page def.PageInfo) ([]*DmDeviceInfo, error)
		CountByFilter(ctx context.Context, filter DeviceFilter) (size int64, err error)
		CountGroupByField(ctx context.Context, filter DeviceFilter, fieldName string) (map[string]int64, error)
		FindOneByProductIDAndDeviceName(ctx context.Context, productID string, deviceName string) (*DmDeviceInfo, error)
		UpdateDeviceInfo(ctx context.Context, data *DmDeviceInfo) error
	}

	customDmDeviceInfoModel struct {
		*defaultDmDeviceInfoModel
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

// NewDmDeviceInfoModel returns a model for the database table.
func NewDmDeviceInfoModel(conn sqlx.SqlConn) DmDeviceInfoModel {
	return &customDmDeviceInfoModel{
		defaultDmDeviceInfoModel: newDmDeviceInfoModel(conn),
	}
}

func (m *customDmDeviceInfoModel) FindByFilter(ctx context.Context, f DeviceFilter, page def.PageInfo) ([]*DmDeviceInfo, error) {
	var resp []*DmDeviceInfo
	sSql := strings.Replace(dmDeviceInfoRows, "`position`", "AsText(`position`) as position", 1)
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

func (m *customDmDeviceInfoModel) CountByFilter(ctx context.Context, f DeviceFilter) (size int64, err error) {
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

func (m *customDmDeviceInfoModel) CountGroupByField(ctx context.Context, f DeviceFilter, columnName string) (map[string]int64, error) {
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

func (m *customDmDeviceInfoModel) InsertDeviceInfo(ctx context.Context, data *DmDeviceInfo) error {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, dmDeviceInfoRowsExpectAutoSet)
	//position为points类型字段,插入时需用函数ST_GeomFromText转换，而不能使用问号
	i := utils.IndexN(query, '?', 12)
	query = query[0:i-1] + "ST_GeomFromText(?))" + query[i+1:len(query)]
	_, err := m.conn.ExecCtx(ctx, query, data.ProductID, data.DeviceName, data.Secret, data.FirstLogin, data.LastLogin, data.Version, data.LogLevel, data.Cert, data.IsOnline, data.Tags, data.Address, data.Position)
	return err
}

func (m *customDmDeviceInfoModel) FindOneByProductIDAndDeviceName(ctx context.Context, productID string, deviceName string) (*DmDeviceInfo, error) {
	var resp DmDeviceInfo
	query := fmt.Sprintf("select %s from %s where `productID` = ? and `deviceName` = ? limit 1", dmDeviceInfoRows, m.table)
	//position字段为point类型 无法直接读取，需使用函数AsText转换后再读取
	query = strings.Replace(query, "`position`", "AsText(`position`) as position", 1)
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

func (m *customDmDeviceInfoModel) UpdateDeviceInfo(ctx context.Context, newData *DmDeviceInfo) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, dmDeviceInfoRowsWithPlaceHolder)
	query = strings.Replace(query, "`position`=?", "`position`=ST_GeomFromText(?)", 1)
	_, err := m.conn.ExecCtx(ctx, query, newData.ProductID, newData.DeviceName, newData.Secret, newData.FirstLogin, newData.LastLogin, newData.Version, newData.LogLevel, newData.Cert, newData.IsOnline, newData.Tags, newData.Address, newData.Position, newData.Id)

	return err
}
