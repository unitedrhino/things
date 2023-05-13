package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysApiInfoModel = (*customSysApiInfoModel)(nil)

type (
	// SysApiInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiInfoModel.
	SysApiInfoModel interface {
		sysApiInfoModel
		Index(ctx context.Context, in *ApiFilter) ([]*SysApiInfo, int64, error)
	}

	customSysApiInfoModel struct {
		*defaultSysApiInfoModel
		sqlx.SqlConn
		api string
	}

	ApiFilter struct {
		Page   *def.PageInfo
		Route  string
		Method int64
		Group  string
		Name   string
	}
)

// NewSysApiInfoModel returns a model for the database table.
func NewSysApiInfoModel(conn sqlx.SqlConn) SysApiInfoModel {
	return &customSysApiInfoModel{
		defaultSysApiInfoModel: newSysApiInfoModel(conn),
	}
}

func (g *ApiFilter) FmtSqlApi(sql sq.SelectBuilder) sq.SelectBuilder {
	if g.Route != "" {
		sql = sql.Where("`route` like ?", "%"+g.Route+"%")
	}
	if g.Method != 0 {
		sql = sql.Where("`method` = ?", g.Method)
	}
	if g.Group != "" {
		sql = sql.Where("`group` like ?", "%"+g.Group+"%")
	}
	if g.Name != "" {
		sql = sql.Where("`name` like ?", "%"+g.Name+"%")
	}

	return sql
}

func (m *customSysApiInfoModel) GetApiCountByFilter(ctx context.Context, f ApiFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.api)
	sql = f.FmtSqlApi(sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = m.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *customSysApiInfoModel) FindApiByFilter(ctx context.Context, f ApiFilter, page *def.PageInfo) ([]*SysApiInfo, error) {
	var resp []*SysApiInfo
	sql := sq.Select(sysApiInfoRows).From(m.api).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset())).OrderBy("createdTime desc")
	sql = f.FmtSqlApi(sql)

	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysApiInfoModel) Index(ctx context.Context, in *ApiFilter) ([]*SysApiInfo, int64, error) {
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	filter := ApiFilter{
		Page:   &page,
		Route:  in.Route,
		Method: in.Method,
		Group:  in.Group,
		Name:   in.Name,
	}

	size, err := m.GetApiCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	resp, err := m.FindApiByFilter(ctx, filter, &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size})
	if err != nil {
		return nil, 0, err
	}

	info := make([]*SysApiInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, &SysApiInfo{
			Id:           v.Id,
			Route:        v.Route,
			Method:       v.Method,
			Name:         v.Name,
			BusinessType: v.BusinessType,
			Group:        v.Group,
		})
	}

	return info, size, nil

}
