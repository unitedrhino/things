package mysql

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserCoreModel = (*customUserCoreModel)(nil)

type (
	// UserCoreModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCoreModel.
	UserCoreModel interface {
		userCoreModel
		FindOneByPhone(ctx context.Context, phone string) (*UserCore, error)
		FindOneByWechat(ctx context.Context, phone string) (*UserCore, error)
		FindOneByUserName(ctx context.Context, userName string) (*UserCore, error)
	}

	customUserCoreModel struct {
		*defaultUserCoreModel
	}
)

// NewUserCoreModel returns a model for the database table.
func NewUserCoreModel(conn sqlx.SqlConn) UserCoreModel {
	return &customUserCoreModel{
		defaultUserCoreModel: newUserCoreModel(conn),
	}
}
func (m *defaultUserCoreModel) FindOneByPhone(ctx context.Context, phone string) (*UserCore, error) {
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userCoreRows, m.table)
	var resp UserCore
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultUserCoreModel) FindOneByWechat(ctx context.Context, weChat string) (*UserCore, error) {
	query := fmt.Sprintf("select %s from %s where `wechat` = ? limit 1", userCoreRows, m.table)
	var resp UserCore
	err := m.conn.QueryRowCtx(ctx, &resp, query, weChat)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserCoreModel) FindOneByUserName(ctx context.Context, userName string) (*UserCore, error) {
	query := fmt.Sprintf("select %s from %s where `userName` = ? limit 1", userCoreRows, m.table)
	var resp UserCore
	err := m.conn.QueryRowCtx(ctx, &resp, query, userName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
