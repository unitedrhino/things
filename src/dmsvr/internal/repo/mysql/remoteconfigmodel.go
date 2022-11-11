package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	RemoteConfigModel interface {
		Index(ctx context.Context, in *RemoteConfigFilter) ([]*RemoteConfigInfo, int64, error)
		GetLastRecord(ctx context.Context, in *RemoteConfigFilter) (*RemoteConfigInfo, error)
	}

	remoteConfigModel struct {
		sqlx.SqlConn
		remoteConfig string
	}

	RemoteConfigFilter struct {
		Page      *def.PageInfo
		ProductID string
	}

	RemoteConfigInfo struct {
		ID          int64
		ProductID   string
		Content     string
		CreatedTime int64
	}
)

func NewRemoteConfigModel(conn sqlx.SqlConn) RemoteConfigModel {
	return &remoteConfigModel{
		SqlConn:      conn,
		remoteConfig: "`product_remote_config`",
	}
}

func (g *RemoteConfigFilter) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if g.ProductID != "" {
		sql = sql.Where("`productID` = ?", g.ProductID)
		sql = sql.OrderBy("createdTime desc")
	}
	return sql
}

func (m *remoteConfigModel) GetRemoteConfigCountByFilter(ctx context.Context, f RemoteConfigFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.remoteConfig)
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

func (m *remoteConfigModel) FindRemoteConfigByFilter(ctx context.Context, f RemoteConfigFilter, page def.PageInfo) ([]*ProductRemoteConfig, error) {
	var resp []*ProductRemoteConfig
	sql := sq.Select(productRemoteConfigRows).From(m.remoteConfig).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (m *remoteConfigModel) FindLastRemoteConfigByFilter(ctx context.Context, f RemoteConfigFilter) ([]*ProductRemoteConfig, error) {
	var resp []*ProductRemoteConfig
	sql := sq.Select(productRemoteConfigRows).From(m.remoteConfig).Limit(1)
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

func (m *remoteConfigModel) Index(ctx context.Context, in *RemoteConfigFilter) ([]*RemoteConfigInfo, int64, error) {
	filter := RemoteConfigFilter{
		ProductID: in.ProductID,
	}
	size, err := m.GetRemoteConfigCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	resp, err := m.FindRemoteConfigByFilter(ctx, filter, logic.ToPageInfo(&dm.PageInfo{Page: in.Page.Page, Size: in.Page.Size}))
	if err != nil {
		return nil, 0, err
	}
	info := make([]*RemoteConfigInfo, 0, len(resp))
	for i, v := range resp {
		if in.Page.Page == 1 && i == 0 {
			continue //最新1条数据不返回
		}
		info = append(info, &RemoteConfigInfo{
			ID:          v.Id,
			ProductID:   v.ProductID,
			Content:     v.Content,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	return info, size - 1, nil
}

func (m *remoteConfigModel) GetLastRecord(ctx context.Context, in *RemoteConfigFilter) (*RemoteConfigInfo, error) {
	filter := RemoteConfigFilter{
		ProductID: in.ProductID,
	}
	resp, err := m.FindLastRemoteConfigByFilter(ctx, filter)
	if err != nil || resp == nil {
		return nil, err
	}
	return &RemoteConfigInfo{
		ID:          resp[0].Id,
		ProductID:   resp[0].ProductID,
		Content:     resp[0].Content,
		CreatedTime: resp[0].CreatedTime.Unix(),
	}, nil
}
