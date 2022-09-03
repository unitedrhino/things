package mysql

import (
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	MenuModel interface {
		Index(in *sys.MenuIndexReq) ([]*MenuInfo, int64, error)
	}

	menuModel struct {
		sqlc.CachedConn
		menuInfo string
	}
)

func NewMenuModel(conn sqlx.SqlConn, c cache.CacheConf) MenuModel {
	return &menuModel{
		CachedConn: sqlc.NewConn(conn, c),
		menuInfo:   "`menu_info`",
	}
}

func (m *menuModel) Index(in *sys.MenuIndexReq) ([]*MenuInfo, int64, error) {
	var resp []*MenuInfo
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	//支持账号模糊匹配
	sql_where := ""
	if in.Name != "" || in.Path != "" {
		sql_where += " where "
		if in.Name != "" && in.Path != "" {
			sql_where += "name like '%" + in.Name + "%' and path like '%" + in.Path + "%'"
		} else if in.Name != "" && in.Path == "" {
			sql_where += "name like '%" + in.Name + "%'"
		} else {
			sql_where += "path like '%" + in.Path + "%'"
		}
	}

	query := fmt.Sprintf("select %s from %s %s limit %d offset %d ",
		menuInfoRows, m.menuInfo, sql_where, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, 0, err
	}

	count := fmt.Sprintf("select count(1) from %s %s", m.menuInfo, sql_where)
	var total int64
	err = m.CachedConn.QueryRowNoCache(&total, count)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}