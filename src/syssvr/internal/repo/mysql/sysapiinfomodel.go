package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	ApiModel interface {
		Index(ctx context.Context, in *ApiFilter) ([]*SysApi, int64, error)
	}
	apiModel struct {
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

func NewApiModel(conn sqlx.SqlConn) ApiModel {
	return &apiModel{
		SqlConn: conn,
		api:     "`sys_api`",
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

func (m *apiModel) GetApiCountByFilter(ctx context.Context, f ApiFilter) (size int64, err error) {
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

func (m *apiModel) FindApiByFilter(ctx context.Context, f ApiFilter, page *def.PageInfo) ([]*SysApi, error) {
	var resp []*SysApi
	sql := sq.Select(sysApiRows).From(m.api).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset())).OrderBy("createdTime desc")
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

func (m *apiModel) Index(ctx context.Context, in *ApiFilter) ([]*SysApi, int64, error) {
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

	info := make([]*SysApi, 0, len(resp))
	for _, v := range resp {
		info = append(info, &SysApi{
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
