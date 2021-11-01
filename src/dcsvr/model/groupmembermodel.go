package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	groupMemberFieldNames          = builderx.RawFieldNames(&GroupMember{})
	groupMemberRows                = strings.Join(groupMemberFieldNames, ",")
	groupMemberRowsExpectAutoSet   = strings.Join(stringx.Remove(groupMemberFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	groupMemberRowsWithPlaceHolder = strings.Join(stringx.Remove(groupMemberFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDcGroupMemberIdPrefix                        = "cache:dc:groupMember:id:"
	cacheDcGroupMemberGroupIDMemberIDMemberTypePrefix = "cache:dc:groupMember:groupID:memberID:memberType:"
)

type (
	GroupMemberModel interface {
		Insert(data GroupMember) (sql.Result, error)
		FindOne(id int64) (*GroupMember, error)
		FindOneByGroupIDMemberIDMemberType(groupID int64, memberID string, memberType int64) (*GroupMember, error)
		Update(data GroupMember) error
		Delete(id int64) error
	}

	defaultGroupMemberModel struct {
		sqlc.CachedConn
		table string
	}

	GroupMember struct {
		Id          int64        `db:"id"`
		GroupID     int64        `db:"groupID"`    // 组id
		MemberID    string       `db:"memberID"`   // 成员id
		MemberType  int64        `db:"memberType"` // 成员类型:1:设备 2:用户
		CreatedTime time.Time    `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
	}
)

func NewGroupMemberModel(conn sqlx.SqlConn, c cache.CacheConf) GroupMemberModel {
	return &defaultGroupMemberModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`group_member`",
	}
}

func (m *defaultGroupMemberModel) Insert(data GroupMember) (sql.Result, error) {
	dcGroupMemberGroupIDMemberIDMemberTypeKey := fmt.Sprintf("%s%v:%v:%v", cacheDcGroupMemberGroupIDMemberIDMemberTypePrefix, data.GroupID, data.MemberID, data.MemberType)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, groupMemberRowsExpectAutoSet)
		return conn.Exec(query, data.GroupID, data.MemberID, data.MemberType, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, dcGroupMemberGroupIDMemberIDMemberTypeKey)
	return ret, err
}

func (m *defaultGroupMemberModel) FindOne(id int64) (*GroupMember, error) {
	dcGroupMemberIdKey := fmt.Sprintf("%s%v", cacheDcGroupMemberIdPrefix, id)
	var resp GroupMember
	err := m.QueryRow(&resp, dcGroupMemberIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMemberRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupMemberModel) FindOneByGroupIDMemberIDMemberType(groupID int64, memberID string, memberType int64) (*GroupMember, error) {
	dcGroupMemberGroupIDMemberIDMemberTypeKey := fmt.Sprintf("%s%v:%v:%v", cacheDcGroupMemberGroupIDMemberIDMemberTypePrefix, groupID, memberID, memberType)
	var resp GroupMember
	err := m.QueryRowIndex(&resp, dcGroupMemberGroupIDMemberIDMemberTypeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `groupID` = ? and `memberID` = ? and `memberType` = ? limit 1", groupMemberRows, m.table)
		if err := conn.QueryRow(&resp, query, groupID, memberID, memberType); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupMemberModel) Update(data GroupMember) error {
	dcGroupMemberIdKey := fmt.Sprintf("%s%v", cacheDcGroupMemberIdPrefix, data.Id)
	dcGroupMemberGroupIDMemberIDMemberTypeKey := fmt.Sprintf("%s%v:%v:%v", cacheDcGroupMemberGroupIDMemberIDMemberTypePrefix, data.GroupID, data.MemberID, data.MemberType)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupMemberRowsWithPlaceHolder)
		return conn.Exec(query, data.GroupID, data.MemberID, data.MemberType, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
	}, dcGroupMemberIdKey, dcGroupMemberGroupIDMemberIDMemberTypeKey)
	return err
}

func (m *defaultGroupMemberModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	dcGroupMemberIdKey := fmt.Sprintf("%s%v", cacheDcGroupMemberIdPrefix, id)
	dcGroupMemberGroupIDMemberIDMemberTypeKey := fmt.Sprintf("%s%v:%v:%v", cacheDcGroupMemberGroupIDMemberIDMemberTypePrefix, data.GroupID, data.MemberID, data.MemberType)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, dcGroupMemberIdKey, dcGroupMemberGroupIDMemberIDMemberTypeKey)
	return err
}

func (m *defaultGroupMemberModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDcGroupMemberIdPrefix, primary)
}

func (m *defaultGroupMemberModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMemberRows, m.table)
	return conn.QueryRow(v, query, primary)
}
