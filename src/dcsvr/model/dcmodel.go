package model

import (
	"fmt"
	"gitee.com/godLei6/things/shared/def"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	DmModel interface {
		CheckMemeberWithGoupID(MemberID1 string,MemberType1 int64,MemberID2 string,MemberType2 int64) (bool, error)
		FindByGroupInfo(page def.PageInfo) ([]*GroupInfo, error)
		FindByGroupMemberGroupID(GroupID  int64,page def.PageInfo) ([]*GroupMember, error)
		FindByGroupMemberMemberID(MemberID string, page def.PageInfo) ([]*GroupMember, error)
		GetCountByGroupInfo() (size int64, err error)
		GetCountByGroupMemberGroupID(GroupID  int64) (size int64, err error)
		GetCountByGroupMemberMemberID(MemberID string) (size int64, err error)

	}

	defaultDcModel struct {
		sqlc.CachedConn
		groupInfo string
		groupMember  string
	}
)

func NewDcModel(conn sqlx.SqlConn, c cache.CacheConf) DmModel {
	return &defaultDcModel{
		CachedConn:  sqlc.NewConn(conn, c),
		groupInfo: "`group_info`",
		groupMember:  "group_member",
	}
}


//获取两个成员是否有在同一组
func (m *defaultDcModel) CheckMemeberWithGoupID(MemberID1 string,MemberType1 int64,MemberID2 string,MemberType2 int64) (bool, error) {
	query := fmt.Sprintf("select count(1) from %s  where memberId = ? and memberType = ? and " +
		"groupID in (select groupID from %s  where memberID=? and memberType =?)",
		m.groupMember,m.groupMember)
	var size int64
	err := m.CachedConn.QueryRowNoCache(&size, query, MemberID1,MemberType1,MemberID2,MemberType2)
	if size > 0{
		return true,err
	}
	return false, err
}

func (m *defaultDcModel) FindByGroupInfo(page def.PageInfo) ([]*GroupInfo, error) {
	var resp []*GroupInfo
	query := fmt.Sprintf("select %s from %s  limit %d offset %d ",
		groupInfoRows, m.groupInfo, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDcModel) GetCountByGroupMemberGroupID(GroupID int64) (size int64, err error) {
	query := fmt.Sprintf("select count(1) from %s where `groupID` = ?",
		m.groupMember)
	err = m.CachedConn.QueryRowNoCache(&size, query, GroupID)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *defaultDcModel) GetCountByGroupMemberMemberID(MemberID string) (size int64, err error) {
	query := fmt.Sprintf("select count(1) from %s where `memberID` = ?",
		m.groupMember)
	err = m.CachedConn.QueryRowNoCache(&size, query, MemberID)
	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}

func (m *defaultDcModel) FindByGroupMemberGroupID(GroupID  int64,page def.PageInfo) ([]*GroupMember, error) {
	var resp []*GroupMember
	query := fmt.Sprintf("select %s from %s where `groupID` = ? limit %d offset %d ",
		groupMemberRows, m.groupMember, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query, GroupID)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDcModel) FindByGroupMemberMemberID(MemberID string, page def.PageInfo) ([]*GroupMember, error) {
	var resp []*GroupMember
	query := fmt.Sprintf("select %s from %s where `memberID` = ? limit %d offset %d ",
		groupMemberRows, m.groupMember, page.GetLimit(), page.GetOffset())
	err := m.CachedConn.QueryRowsNoCache(&resp, query, MemberID)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDcModel) GetCountByGroupInfo() (size int64, err error) {
	query := fmt.Sprintf("select count(1)  from %s ",
		m.groupInfo)
	err = m.CachedConn.QueryRowNoCache(&size, query)
	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}
