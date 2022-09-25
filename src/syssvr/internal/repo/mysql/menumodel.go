package mysql

import (
	"fmt"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	MenuModel interface {
		Index(in *sys.MenuIndexReq) ([]*MenuInfo, error)
		DeleteMenu(MenuId int64) error
	}

	menuModel struct {
		sqlc.CachedConn
		menuInfo string
		roleMenu string
	}
)

func NewMenuModel(conn sqlx.SqlConn, c cache.CacheConf) MenuModel {
	return &menuModel{
		CachedConn: sqlc.NewConn(conn, c),
		menuInfo:   "`menu_info`",
		roleMenu:   "`role_menu`",
	}
}

func (m *menuModel) Index(in *sys.MenuIndexReq) ([]*MenuInfo, error) {
	var resp []*MenuInfo

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

	query := fmt.Sprintf("select %s from %s %s",
		menuInfoRows, m.menuInfo, sql_where)
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *menuModel) DeleteMenu(MenuId int64) error {
	m.Transact(func(session sqlx.Session) error {
		//1.从菜单表删除角色
		query := fmt.Sprintf("delete from %s where id = %d",
			m.menuInfo, MenuId)
		_, err := session.Exec(query)
		if err != nil {
			return err
		}
		//2.从角色菜单关系表删除关联菜单项
		query = fmt.Sprintf("delete from %s where  menuID = %d",
			m.roleMenu, MenuId)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
