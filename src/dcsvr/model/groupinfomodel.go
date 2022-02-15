package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	groupInfoFieldNames          = builder.RawFieldNames(&GroupInfo{})
	groupInfoRows                = strings.Join(groupInfoFieldNames, ",")
	groupInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(groupInfoFieldNames, "`create_time`", "`update_time`"), ",")
	groupInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(groupInfoFieldNames, "`groupID`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDcGroupInfoGroupIDPrefix = "cache:dc:groupInfo:groupID:"
)

type (
	GroupInfoModel interface {
		Insert(data GroupInfo) (sql.Result, error)
		FindOne(groupID int64) (*GroupInfo, error)
		Update(data GroupInfo) error
		Delete(groupID int64) error
	}

	defaultGroupInfoModel struct {
		sqlc.CachedConn
		table string
	}

	GroupInfo struct {
		GroupID     int64        `db:"groupID"` // 组id
		Name        string       `db:"name"`    // 组名
		Uid         int64        `db:"uid"`     // 管理员用户id
		CreatedTime time.Time    `db:"createdTime"`
		UpdatedTime sql.NullTime `db:"updatedTime"`
		DeletedTime sql.NullTime `db:"deletedTime"`
	}
)

func NewGroupInfoModel(conn sqlx.SqlConn, c cache.CacheConf) GroupInfoModel {
	return &defaultGroupInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`group_info`",
	}
}

func (m *defaultGroupInfoModel) Insert(data GroupInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, groupInfoRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.GroupID, data.Name, data.Uid, data.CreatedTime, data.UpdatedTime, data.DeletedTime)

	return ret, err
}

func (m *defaultGroupInfoModel) FindOne(groupID int64) (*GroupInfo, error) {
	dcGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcGroupInfoGroupIDPrefix, groupID)
	var resp GroupInfo
	err := m.QueryRow(&resp, dcGroupInfoGroupIDKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `groupID` = ? limit 1", groupInfoRows, m.table)
		return conn.QueryRow(v, query, groupID)
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

func (m *defaultGroupInfoModel) Update(data GroupInfo) error {
	dcGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcGroupInfoGroupIDPrefix, data.GroupID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `groupID` = ?", m.table, groupInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Uid, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.GroupID)
	}, dcGroupInfoGroupIDKey)
	return err
}

func (m *defaultGroupInfoModel) Delete(groupID int64) error {

	dcGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcGroupInfoGroupIDPrefix, groupID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `groupID` = ?", m.table)
		return conn.Exec(query, groupID)
	}, dcGroupInfoGroupIDKey)
	return err
}

func (m *defaultGroupInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDcGroupInfoGroupIDPrefix, primary)
}

func (m *defaultGroupInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `groupID` = ? limit 1", groupInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
