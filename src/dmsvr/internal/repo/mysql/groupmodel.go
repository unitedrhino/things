package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	GroupModel interface {
		Index(ctx context.Context, in *GroupFilter) ([]*GroupInformation, int64, error)
		IndexAll(ctx context.Context, in *GroupFilter) ([]*GroupInformation, error)
		IndexGD(ctx context.Context, in *GroupDeviceFilter) ([]*GroupDevice, int64, error)
		Delete(ctx context.Context, groupID int64) error
		GroupDeviceCreate(ctx context.Context, groupID int64, list []*device.Core) error
		GroupDeviceDelete(ctx context.Context, groupID int64, list []*device.Core) error
	}

	groupModel struct {
		sqlx.SqlConn
		groupInfo   string
		groupDevice string
		deviceInfo  string
	}
	GroupFilter struct {
		Page      *def.PageInfo
		ParentID  int64
		GroupName string
		Tags      map[string]string
	}
	GroupDeviceFilter struct {
		Page       *def.PageInfo
		GroupID    int64
		ProductID  string
		DeviceName string
	}

	GroupInformation struct {
		GroupID     int64
		GroupName   string
		ParentID    int64
		Desc        string
		CreatedTime int64
		Tags        map[string]string
	}
)

func NewGroupModel(conn sqlx.SqlConn) GroupModel {
	return &groupModel{
		SqlConn:     conn,
		groupInfo:   "`group_info`",
		groupDevice: "`group_device`",
		deviceInfo:  "`device_info`",
	}
}

func (g *GroupFilter) FmtSql(sql sq.SelectBuilder, parentFlag bool) sq.SelectBuilder {
	if parentFlag == true && g.ParentID != 0 {
		sql = sql.Where("`parentID`=?", g.ParentID)
	}
	if g.GroupName != "" {
		sql = sql.Where("`groupName` like ?", "%"+g.GroupName+"%")
	}
	if g.Tags != nil {
		for k, v := range g.Tags {
			sql = sql.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return sql
}
func (g *GroupDeviceFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if g.GroupID != 0 {
		sql = sql.Where("`groupID`=?", g.GroupID)
	}
	if g.ProductID != "" {
		sql = sql.Where("`productID`=?", g.ProductID)
	}
	if g.DeviceName != "" {
		sql = sql.Where("`deviceName`=?", g.DeviceName)
	}

	return sql
}

func (m *groupModel) GetGroupsCountByFilter(ctx context.Context, f GroupFilter, parentFlag bool) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.groupInfo)
	sql = f.FmtSql(sql, parentFlag)
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
func (m *groupModel) FindGroupInfoByFilter(ctx context.Context, f GroupFilter, page def.PageInfo, parentFlag bool) ([]*GroupInfo, error) {
	var resp []*GroupInfo
	sql := sq.Select(groupInfoRows).From(m.groupInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql, parentFlag)

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

func (m *groupModel) GetGroupDeviceCountByFilter(ctx context.Context, f GroupDeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.groupDevice)
	sql = f.FmtSql(sql)
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
func (m *groupModel) FindGroupDeviceByFilter(ctx context.Context, f GroupDeviceFilter, page def.PageInfo) ([]*GroupDevice, error) {
	var resp []*GroupDevice
	sql := sq.Select(groupDeviceRows).From(m.groupDevice).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSql(sql)

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

func (m *groupModel) Index(ctx context.Context, in *GroupFilter) ([]*GroupInformation, int64, error) {

	filter := GroupFilter{
		ParentID:  in.ParentID,
		GroupName: in.GroupName,
		Tags:      in.Tags,
	}
	size, err := m.GetGroupsCountByFilter(ctx, filter, true)
	if err != nil {
		return nil, 0, err
	}

	dg, err := m.FindGroupInfoByFilter(ctx, filter, logic.ToPageInfo(&dm.PageInfo{Page: in.Page.Page, Size: in.Page.Size}), true)
	if err != nil {
		return nil, 0, err
	}

	info := make([]*GroupInformation, 0, len(dg))
	for _, v := range dg {
		var tags map[string]string
		if v.Tags != "" {
			_ = json.Unmarshal([]byte(v.Tags), &tags)
		}
		info = append(info, &GroupInformation{
			GroupID:     v.GroupID,
			GroupName:   v.GroupName,
			ParentID:    v.ParentID,
			Desc:        v.Desc,
			CreatedTime: v.CreatedTime.Unix(),
			Tags:        tags,
		})
	}

	return info, size, nil
}

func (m *groupModel) IndexAll(ctx context.Context, in *GroupFilter) ([]*GroupInformation, error) {

	filter := GroupFilter{
		GroupName: in.GroupName,
		Tags:      in.Tags,
	}

	dg, err := m.FindGroupInfoByFilter(ctx, filter, def.PageInfo{Size: 100000, Page: 1}, false)
	if err != nil {
		return nil, err
	}

	info := make([]*GroupInformation, 0, len(dg))
	for _, v := range dg {
		var tags map[string]string
		if v.Tags != "" {
			_ = json.Unmarshal([]byte(v.Tags), &tags)
		}
		info = append(info, &GroupInformation{
			GroupID:     v.GroupID,
			GroupName:   v.GroupName,
			ParentID:    v.ParentID,
			Desc:        v.Desc,
			CreatedTime: v.CreatedTime.Unix(),
			Tags:        tags,
		})
	}

	return info, nil
}

func (m *groupModel) IndexGD(ctx context.Context, in *GroupDeviceFilter) ([]*GroupDevice, int64, error) {

	filter := GroupDeviceFilter{
		GroupID:    in.GroupID,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	size, err := m.GetGroupDeviceCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	dg, err := m.FindGroupDeviceByFilter(ctx, filter, logic.ToPageInfo(&dm.PageInfo{Page: in.Page.Page, Size: in.Page.Size}))
	if err != nil {
		return nil, 0, err
	}

	return dg, size, nil
}

func (m *groupModel) Delete(ctx context.Context, groupID int64) error {
	return m.Transact(func(session sqlx.Session) error {
		//1.查詢是否存在子分组，如果存在则不允许删除该分组
		sql := fmt.Sprintf("select count(1) from %s where parentID = %d", m.groupInfo, groupID)
		var count int64
		err := session.QueryRow(&count, sql)
		if err != sqlc.ErrNotFound && count > 0 {
			return errors.NotEmpty.WithMsg("存在子分组").AddDetailf("the group have sun group can not delete.")
		}

		//2.从group_info表删除角色
		query := fmt.Sprintf("delete from %s where groupID = %d", m.groupInfo, groupID)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}

		//3.从group_device关系表删除关联项
		query = fmt.Sprintf("delete from %s where  groupID = %d",
			m.groupDevice, groupID)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *groupModel) GroupDeviceCreate(ctx context.Context, groupID int64, list []*device.Core) error {
	return m.Transact(func(session sqlx.Session) error {
		for _, v := range list {
			var count int64
			query := fmt.Sprintf("select count(1) from %s where productID = %s and deviceName = %s", m.deviceInfo, v.ProductID, v.DeviceName)
			err := session.QueryRow(&count, query)
			if err != sqlc.ErrNotFound && count > 0 {
				var resp GroupDevice
				query = fmt.Sprintf("select %s from %s where `groupID` = ? and `productID` = ? and `deviceName` = ?  limit 1", groupDeviceRows, m.groupDevice)
				err = session.QueryRow(&resp, query, groupID, v.ProductID, v.DeviceName)
				if err == sqlc.ErrNotFound {
					query := fmt.Sprintf("insert into %s (groupID,productID,deviceName) values (%d, '%s', '%s')",
						m.groupDevice, groupID, v.ProductID, v.DeviceName)
					_, err = session.Exec(query)
					if err != nil {
						return nil
					}
				}
			}
		}
		return nil
	})
}

func (m *groupModel) GroupDeviceDelete(ctx context.Context, groupID int64, list []*device.Core) error {
	return m.Transact(func(session sqlx.Session) error {
		for _, v := range list {
			query := fmt.Sprintf("delete from %s where groupID = %d and productID = '%s' and deviceName = '%s' ",
				m.groupDevice, groupID, v.ProductID, v.DeviceName)
			_, err := session.Exec(query)
			if err != nil {
				return nil
			}
		}
		return nil
	})
}
