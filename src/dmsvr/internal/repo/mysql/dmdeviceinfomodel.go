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
		FindByFilter(ctx context.Context, filter DeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error)
		CountByFilter(ctx context.Context, filter DeviceFilter) (size int64, err error)
		CountGroupByField(ctx context.Context, filter DeviceFilter, fieldName string) (map[string]int64, error)
		FindOneByProductIDDeviceName(ctx context.Context, productID string, deviceName string) (*DmDeviceInfo, error)
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
		IsOnline    []int64
		Range       int64
		Position    string
		DeviceAlias string
	}
)

func (d *DeviceFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if d.ProductID != "" {
		sql = sql.Where("`ProductID` = ?", d.ProductID)
	}
	if d.DeviceName != "" {
		sql = sql.Where("`DeviceName` like ?", "%"+d.DeviceName+"%")
	}
	if d.DeviceAlias != "" {
		sql = sql.Where("`DeviceAlias` like ?", "%"+d.DeviceAlias+"%")
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

func (m *customDmDeviceInfoModel) FindByFilter(ctx context.Context, f DeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error) {
	var resp []*DmDeviceInfo

	sSql := strings.Replace(dmDeviceInfoRows, "`position`", "AsText(`position`) as position", 1)
	sql := sq.Select(sSql).From(m.table).
		Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset())).OrderBy(page.GetOrders()...)

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
	if err != nil {
		return nil, err
	}

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
	table := m.table
	fields := dmDeviceInfoRowsExpectAutoSet
	params := []any{ //注意：要和 fields的 字段顺序 对上
		data.ProductID,
		data.DeviceName,
		data.Secret,
		data.Cert,
		data.Imei,
		data.Mac,
		data.Version,
		data.HardInfo,
		data.SoftInfo,
		data.Position, //注意，记住这里的 position pos=10
		data.Address,
		data.Tags,
		data.IsOnline,
		data.FirstLogin,
		data.LastLogin,
		data.LogLevel,
		data.DeviceAlias,
	}
	valsPlace := utils.NewFillPlace(len(params)) //生成 ?,?,... (有len个?)

	//SQL基本结构
	query := fmt.Sprintf("insert into %s (%s) values (%s)", table, fields, valsPlace)
	//SQL特殊处理（position为points类型字段,插入时需用函数ST_GeomFromText转换，而不能使用问号）
	pos := utils.IndexN(query, '?', 10) //注意：这里是上面的 position pos 10，如位置有变值要跟着改（比如加了字段...）
	query = query[0:pos-1] + "ST_GeomFromText(?)," + query[pos+1:]

	_, err := m.conn.ExecCtx(ctx, query, params...)
	return err
}

func (m *customDmDeviceInfoModel) FindOneByProductIDDeviceName(ctx context.Context, productID string, deviceName string) (*DmDeviceInfo, error) {
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
	table := m.table
	fields := dmDeviceInfoRowsWithPlaceHolder
	params := []any{ //注意：要和 fields的 字段顺序 对上
		newData.ProductID,
		newData.DeviceName,
		newData.Secret,
		newData.Cert,
		newData.Imei,
		newData.Mac,
		newData.Version,
		newData.HardInfo,
		newData.SoftInfo,
		newData.Position,
		newData.Address,
		newData.Tags,
		newData.IsOnline,
		newData.FirstLogin,
		newData.LastLogin,
		newData.LogLevel,
		newData.DeviceAlias,
		newData.Id,
	}

	//SQL基本结构
	query := fmt.Sprintf("update %s set %s where `id` = ?", table, fields)
	//SQL特殊处理（position为points类型字段,插入时需用函数ST_GeomFromText转换，而不能使用问号）
	query = strings.Replace(query, "`position`=?", "`position`=ST_GeomFromText(?)", 1)

	_, err := m.conn.ExecCtx(ctx, query, params...)
	return err
}
