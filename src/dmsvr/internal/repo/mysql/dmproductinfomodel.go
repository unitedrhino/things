package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DmProductInfoModel = (*customDmProductInfoModel)(nil)

type (
	// DmProductInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmProductInfoModel.
	DmProductInfoModel interface {
		dmProductInfoModel
		FindByFilter(ctx context.Context, filter ProductFilter, page *def.PageInfo) ([]*DmProductInfo, error)
		CountByFilter(ctx context.Context, filter ProductFilter) (size int64, err error)
	}
	ProductFilter struct {
		DeviceType  int64
		ProductName string
		ProductIDs  []string
		Tags        map[string]string
	}
	customDmProductInfoModel struct {
		*defaultDmProductInfoModel
	}
)

func (p *ProductFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.DeviceType != 0 {
		sql = sql.Where("deviceType=?", p.DeviceType)
	}
	if p.ProductName != "" {
		sql = sql.Where("productName like ?", "%"+p.ProductName+"%")
	}
	if len(p.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("productID in (%v)", store.ArrayToSql(p.ProductIDs)))
	}
	if p.Tags != nil {
		for k, v := range p.Tags {
			sql = sql.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return sql
}

// NewDmProductInfoModel returns a model for the database table.
func NewDmProductInfoModel(conn sqlx.SqlConn) DmProductInfoModel {
	return &customDmProductInfoModel{
		defaultDmProductInfoModel: newDmProductInfoModel(conn),
	}
}

func (m *customDmProductInfoModel) FindByFilter(ctx context.Context, f ProductFilter, page *def.PageInfo) ([]*DmProductInfo, error) {
	var resp []*DmProductInfo
	sql := sq.Select(dmProductInfoRows).From(m.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (m *customDmProductInfoModel) CountByFilter(ctx context.Context, f ProductFilter) (size int64, err error) {
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
