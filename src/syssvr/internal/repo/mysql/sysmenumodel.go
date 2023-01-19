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
		Index(in *sys.MenuIndexReq) ([]*SysMenuInfo, error)
		DeleteMenu(MenuId int64) error
		InsertMenuID(data *SysMenuInfo, RoleId int64) error
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
		menuInfo:   "`sys_menu_info`",
		roleMenu:   "`sys_role_menu`",
	}
}

func (m *menuModel) Index(in *sys.MenuIndexReq) ([]*SysMenuInfo, error) {
	var resp []*SysMenuInfo

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
		sysMenuInfoRows, m.menuInfo, sql_where)
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *menuModel) DeleteMenu(MenuId int64) error {
	return m.Transact(func(session sqlx.Session) error {
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
}

func (m *menuModel) InsertMenuID(data *SysMenuInfo, RoleId int64) error {
	return m.Transact(func(session sqlx.Session) error {
		//1.向menu_info表插入菜单项
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.menuInfo, sysMenuInfoRowsExpectAutoSet)
		ret, err := session.Exec(query, data.ParentID, data.Type, data.Order, data.Name, data.Path, data.Component, data.Icon, data.Redirect, data.BackgroundUrl, data.HideInMenu)
		if err != nil {
			return err
		}
		insertID, err := ret.LastInsertId()
		//2.向role_menu表插入关联菜单项
		query = fmt.Sprintf("insert into %s (roleID, menuID) values (%d, %d)",
			m.roleMenu, RoleId, insertID)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}

		return nil
	})
}
