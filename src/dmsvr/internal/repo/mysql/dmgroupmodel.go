package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/userHeader"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DmGroupModel interface {
		Index(ctx context.Context, in *GroupFilter) ([]*GroupInformation, int64, error)
		IndexAll(ctx context.Context, in *GroupFilter) ([]*GroupInformation, error)
		IndexGD(ctx context.Context, in *GroupDeviceFilter) ([]*DmGroupDevice, int64, error)
		Delete(ctx context.Context, groupID int64) error
		GroupDeviceCreate(ctx context.Context, groupID int64, list []*devices.Core) error
		GroupDeviceDelete(ctx context.Context, groupID int64, list []*devices.Core) error
		DeleteDevice(ctx context.Context, device *devices.Core) error
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

func NewDmGroupModel(conn sqlx.SqlConn) DmGroupModel {
	return &groupModel{
		SqlConn:     conn,
		groupInfo:   "`dm_group_info`",
		groupDevice: "`dm_group_device`",
		deviceInfo:  "`dm_device_info`",
	}
}

func (m *groupModel) FmtGroupSql(ctx context.Context, f GroupFilter, sql sq.SelectBuilder, parentFlag bool) sq.SelectBuilder {
	if parentFlag == true && f.ParentID != 0 {
		sql = sql.Where("`parentID`=?", f.ParentID)
	}
	if f.GroupName != "" {
		sql = sql.Where("`groupName` like ?", "%"+f.GroupName+"%")
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			sql = sql.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return sql
}

func (m *groupModel) FmtGroupDeviceSql(ctx context.Context, f GroupDeviceFilter, sql sq.SelectBuilder) sq.SelectBuilder {
	sql = sql.LeftJoin(fmt.Sprintf("%s as di on di.productID=gd.productID and di.deviceName=gd.deviceName", m.deviceInfo))

	//数据权限条件（企业版功能）
	if uc := userHeader.GetUserCtxOrNil(ctx); uc != nil && !uc.IsAllData { //存在用户态&&无所有数据权限
		mdProjectID := userHeader.GetMetaProjectID(ctx)
		if mdProjectID != 0 {
			sql = sql.Where("di.`ProjectID` = ?", mdProjectID)
		}
	}
	//业务过滤条件
	if f.GroupID != 0 {
		sql = sql.Where("gd.`groupID`=?", f.GroupID)
	}
	if f.ProductID != "" {
		sql = sql.Where("gd.`productID`=?", f.ProductID)
	}
	if f.DeviceName != "" {
		sql = sql.Where("gd.`deviceName`=?", f.DeviceName)
	}

	return sql
}

func (m *groupModel) CountGroupsCountByFilter(ctx context.Context, f GroupFilter, parentFlag bool) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.groupInfo)
	sql = m.FmtGroupSql(ctx, f, sql, parentFlag)
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
func (m *groupModel) FindGroupInfoByFilter(ctx context.Context, f GroupFilter, page *def.PageInfo, parentFlag bool) ([]*DmGroupInfo, error) {
	var resp []*DmGroupInfo
	sql := sq.Select(dmGroupInfoRows).From(m.groupInfo).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = m.FmtGroupSql(ctx, f, sql, parentFlag)

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

func (m *groupModel) CountGroupDeviceCountByFilter(ctx context.Context, f GroupDeviceFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.groupDevice + " as gd")
	sql = m.FmtGroupDeviceSql(ctx, f, sql)

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
func (m *groupModel) FindGroupDeviceByFilter(ctx context.Context, f GroupDeviceFilter, page *def.PageInfo) ([]*DmGroupDevice, error) {
	var resp []*DmGroupDevice
	sql := sq.Select("gd.*").From(m.groupDevice + " as gd")
	sql = sql.Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = m.FmtGroupDeviceSql(ctx, f, sql)

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
	size, err := m.CountGroupsCountByFilter(ctx, filter, true)
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

	dg, err := m.FindGroupInfoByFilter(ctx, filter, nil, false)
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

func (m *groupModel) IndexGD(ctx context.Context, in *GroupDeviceFilter) ([]*DmGroupDevice, int64, error) {
	filter := GroupDeviceFilter{
		GroupID:    in.GroupID,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	size, err := m.CountGroupDeviceCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	dg, err := m.FindGroupDeviceByFilter(ctx, filter, in.Page)
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

func (m *groupModel) GroupDeviceCreate(ctx context.Context, groupID int64, list []*devices.Core) error {
	return m.Transact(func(session sqlx.Session) error {
		for _, v := range list {
			if v == nil {
				continue
			}
			var count int64
			query := fmt.Sprintf("select count(1) from %s where productID = %q and deviceName = %q", m.deviceInfo, v.ProductID, v.DeviceName)
			err := session.QueryRow(&count, query)
			if err != nil {
				logx.WithContext(ctx).Errorf("groupModel.deviceInfoTable.QueryRowCtx data:%v err:%v", v, err)
				continue
			}
			if count == 0 {
				return errors.Parameter.WithMsgf("设备不存在:产品ID:%v,设备名:%", v.ProductID, v.DeviceName)
			}
			query = fmt.Sprintf("insert into %s (groupID,productID,deviceName) values (%d, '%s', '%s') ON duplicate KEY UPDATE id = id",
				m.groupDevice, groupID, v.ProductID, v.DeviceName)
			_, err = session.Exec(query)
			if err != nil {
				logx.WithContext(ctx).Errorf("groupModel.GroupDeviceCreate data:%v err:%v", v, err)
				continue
			}
		}
		return nil
	})
}

func (m *groupModel) GroupDeviceDelete(ctx context.Context, groupID int64, list []*devices.Core) error {
	return m.Transact(func(session sqlx.Session) error {
		for _, v := range list {
			if v == nil {
				continue
			}

			query := fmt.Sprintf("delete from %s where groupID = %d and productID = '%s' and deviceName = '%s' ",
				m.groupDevice, groupID, v.ProductID, v.DeviceName)
			_, err := session.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (m *groupModel) DeleteDevice(ctx context.Context, device *devices.Core) error {
	query := fmt.Sprintf("delete from %s where  productID = '%s' and deviceName = '%s' ",
		m.groupDevice, device.ProductID, device.DeviceName)
	_, err := m.SqlConn.ExecCtx(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
