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
	groupInfoFieldNames          = builderx.RawFieldNames(&GroupInfo{})
	groupInfoRows                = strings.Join(groupInfoFieldNames, ",")
	groupInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(groupInfoFieldNames, "`create_time`", "`update_time`"), ",")
	groupInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(groupInfoFieldNames, "`groupID`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheDcsvrGroupInfoGroupIDPrefix = "cache:dcsvr:groupInfo:groupID:"
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
		table:      "`groupInfo`",
	}
}

func (m *defaultGroupInfoModel) Insert(data GroupInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, groupInfoRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.GroupID, data.Name, data.Uid, data.CreatedTime, data.UpdatedTime, data.DeletedTime)

	return ret, err
}

func (m *defaultGroupInfoModel) FindOne(groupID int64) (*GroupInfo, error) {
	dcsvrGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcsvrGroupInfoGroupIDPrefix, groupID)
	var resp GroupInfo
	err := m.QueryRow(&resp, dcsvrGroupInfoGroupIDKey, func(conn sqlx.SqlConn, v interface{}) error {
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
	dcsvrGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcsvrGroupInfoGroupIDPrefix, data.GroupID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `groupID` = ?", m.table, groupInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Uid, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.GroupID)
	}, dcsvrGroupInfoGroupIDKey)
	return err
}

func (m *defaultGroupInfoModel) Delete(groupID int64) error {

	dcsvrGroupInfoGroupIDKey := fmt.Sprintf("%s%v", cacheDcsvrGroupInfoGroupIDPrefix, groupID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `groupID` = ?", m.table)
		return conn.Exec(query, groupID)
	}, dcsvrGroupInfoGroupIDKey)
	return err
}

func (m *defaultGroupInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDcsvrGroupInfoGroupIDPrefix, primary)
}

func (m *defaultGroupInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `groupID` = ? limit 1", groupInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
