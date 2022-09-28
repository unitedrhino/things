package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	GroupModel interface {
		Index(ctx context.Context, in *dm.GroupInfoIndexReq) ([]*dm.GroupInfo, int64, error)
		IndexAll(ctx context.Context, in *dm.GroupInfoIndexReq) ([]*dm.GroupInfo, error)
		IndexGD(ctx context.Context, in *dm.GroupDeviceIndexReq) ([]*GroupDevice, int64, error)
		Delete(ctx context.Context, groupID int64) error
		GroupDeviceCreate(ctx context.Context, groupID int64, list []*dm.DeviceInfoReadReq) error
		GroupDeviceDelete(ctx context.Context, groupID int64, list map[string]string) error
	}

	groupModel struct {
		sqlx.SqlConn
		groupInfo   string
		groupDevice string
	}
	GroupFilter struct {
		ParentID  int64
		GroupName string
		Tags      map[string]string
	}
	GroupDeviceFilter struct {
		GroupID    int64
		productID  string
		DeviceName string
	}
)

func NewGroupModel(conn sqlx.SqlConn, c cache.CacheConf) GroupModel {
	return &groupModel{
		SqlConn:     conn,
		groupInfo:   "`group_info`",
		groupDevice: "`group_device`",
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
	if g.productID != "" {
		sql = sql.Where("`productID`=?", g.productID)
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

func (m *groupModel) Index(ctx context.Context, in *dm.GroupInfoIndexReq) ([]*dm.GroupInfo, int64, error) {

	filter := GroupFilter{
		ParentID:  in.ParentID,
		GroupName: in.GroupName,
		Tags:      in.Tags,
	}
	size, err := m.GetGroupsCountByFilter(ctx, filter, true)
	if err != nil {
		return nil, 0, err
	}

	dg, err := m.FindGroupInfoByFilter(ctx, filter, def.PageInfo{Size: in.Page.Size, Page: in.Page.Page}, true)
	if err != nil {
		return nil, 0, err
	}

	info := make([]*dm.GroupInfo, 0, len(dg))
	for _, v := range dg {
		var tags map[string]string
		if v.Tags.String != "" {
			_ = json.Unmarshal([]byte(v.Tags.String), &tags)
		}
		info = append(info, &dm.GroupInfo{
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

func (m *groupModel) IndexAll(ctx context.Context, in *dm.GroupInfoIndexReq) ([]*dm.GroupInfo, error) {

	filter := GroupFilter{
		GroupName: in.GroupName,
		Tags:      in.Tags,
	}

	dg, err := m.FindGroupInfoByFilter(ctx, filter, def.PageInfo{Size: 100000, Page: 1}, false)
	if err != nil {
		return nil, err
	}

	info := make([]*dm.GroupInfo, 0, len(dg))
	for _, v := range dg {
		var tags map[string]string
		if v.Tags.String != "" {
			_ = json.Unmarshal([]byte(v.Tags.String), &tags)
		}
		info = append(info, &dm.GroupInfo{
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

func (m *groupModel) IndexGD(ctx context.Context, in *dm.GroupDeviceIndexReq) ([]*GroupDevice, int64, error) {

	filter := GroupDeviceFilter{
		GroupID:    in.GroupID,
		productID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	size, err := m.GetGroupDeviceCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	dg, err := m.FindGroupDeviceByFilter(ctx, filter, def.PageInfo{Size: in.Page.Size, Page: in.Page.Page})
	if err != nil {
		return nil, 0, err
	}

	return dg, size, nil
}

func (m *groupModel) Delete(ctx context.Context, groupID int64) error {
	m.Transact(func(session sqlx.Session) error {
		//1.从group_info表删除角色
		query := fmt.Sprintf("delete from %s where groupID = %d", m.groupInfo, groupID)
		_, err := session.Exec(query)
		if err != nil {
			return err
		}
		//2.从group_device关系表删除关联项
		query = fmt.Sprintf("delete from %s where  groupID = %d",
			m.groupDevice, groupID)
		_, err = session.Exec(query)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (m *groupModel) GroupDeviceCreate(ctx context.Context, groupID int64, list []*dm.DeviceInfoReadReq) error {
	m.Transact(func(session sqlx.Session) error {
		for _, v := range list {
			query1 := fmt.Sprintf("select id from %s where groupID=%d and productID='%s' and deviceName='%s'",
				m.groupDevice, groupID, v.ProductID, v.DeviceName)
			resp, _ := session.Exec(query1)
			count, err := resp.RowsAffected()
			if count == 0 {
				query := fmt.Sprintf("insert into %s (groupID,productID,deviceName) values (%d, '%s', '%s')",
					m.groupDevice, groupID, v.ProductID, v.DeviceName)
				_, err = session.Exec(query)
				if err != nil {
					return nil
				}
			}

		}
		return nil
	})
	return nil
}

func (m *groupModel) GroupDeviceDelete(ctx context.Context, groupID int64, list map[string]string) error {
	m.Transact(func(session sqlx.Session) error {
		for i, v := range list {
			query := fmt.Sprintf("delete from %s where groupID = %d and productID = '%s' and deviceName = '%s' ",
				m.groupDevice, groupID, i, v)
			_, err := session.Exec(query)
			if err != nil {
				return nil
			}
		}
		return nil
	})
	return nil
}
