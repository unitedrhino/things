package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	LogModel interface {
		LoginLogIndex(ctx context.Context, in *LoginLogFilter) ([]*SysLoginLog, int64, error)
		OperLogIndex(ctx context.Context, in *OperLogFilter) ([]*SysOperLog, int64, error)
	}

	DateRange struct {
		Start string
		End   string
	}

	logModel struct {
		sqlx.SqlConn
		loginLog string
		operLog  string
	}

	LoginLogFilter struct {
		Page          *def.PageInfo
		IpAddr        string
		LoginLocation string
		Data          *DateRange
	}

	OperLogFilter struct {
		Page         *def.PageInfo
		OperName     string
		OperUserName string
		BusinessType int64
	}
)

func NewLogModel(conn sqlx.SqlConn) LogModel {
	return &logModel{
		SqlConn:  conn,
		loginLog: "`sys_login_log`",
		operLog:  "`sys_oper_log`",
	}
}

func (g *OperLogFilter) FmtSqlOperLog(sql sq.SelectBuilder) sq.SelectBuilder {
	if g.OperName != "" {
		sql = sql.Where("`operName` like ?", "%"+g.OperName+"%")
	}
	if g.OperUserName != "" {
		sql = sql.Where("`operUserName` like ?", "%"+g.OperUserName+"%")
	}
	if g.BusinessType > 0 {
		sql = sql.Where("`businessType`= ?", g.OperName)
	}

	return sql
}

func (m *logModel) GetOperLogCountByFilter(ctx context.Context, f OperLogFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.operLog)
	sql = f.FmtSqlOperLog(sql)
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

func (m *logModel) FindOperLogByFilter(ctx context.Context, f OperLogFilter, page *def.PageInfo) ([]*SysOperLog, error) {
	var resp []*SysOperLog
	sql := sq.Select(sysOperLogRows).From(m.operLog).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSqlOperLog(sql)

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

func (g *LoginLogFilter) FmtSqlLoginLog(sql sq.SelectBuilder) sq.SelectBuilder {
	if g.IpAddr != "" {
		sql = sql.Where("`ipAddr`= ?", g.IpAddr)
	}
	if g.LoginLocation != "" {
		sql = sql.Where("`loginLocation` like ?", "%"+g.LoginLocation+"%")
	}
	if g.Data != nil && g.Data.Start != "" && g.Data.End != "" {
		sql = sql.Where("`createdTime` >= ? and `createdTime` <= ?", g.Data.Start, g.Data.End)
	}

	return sql
}

func (m *logModel) GetLoginLogCountByFilter(ctx context.Context, f LoginLogFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(m.loginLog)
	sql = f.FmtSqlLoginLog(sql)
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

func (m *logModel) FindLoginLogByFilter(ctx context.Context, f LoginLogFilter, page *def.PageInfo) ([]*SysLoginLog, error) {
	var resp []*SysLoginLog
	sql := sq.Select(sysLoginLogRows).From(m.loginLog).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = f.FmtSqlLoginLog(sql)

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

func (m *logModel) OperLogIndex(ctx context.Context, in *OperLogFilter) ([]*SysOperLog, int64, error) {
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	filter := OperLogFilter{
		Page:         &page,
		OperName:     in.OperName,
		OperUserName: in.OperUserName,
		BusinessType: in.BusinessType,
	}

	size, err := m.GetOperLogCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	resp, err := m.FindOperLogByFilter(ctx, filter, &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size})
	if err != nil {
		return nil, 0, err
	}

	info := make([]*SysOperLog, 0, len(resp))
	for _, v := range resp {
		info = append(info, &SysOperLog{
			Id:           v.Id,
			OperUid:      v.OperUid,
			OperUserName: v.OperUserName,
			OperName:     v.OperName,
			BusinessType: v.BusinessType,
			Uri:          v.Uri,
			OperIpAddr:   v.OperIpAddr,
			OperLocation: v.OperLocation,
			Req:          v.Req,
			Resp:         v.Resp,
			Code:         v.Code,
			Msg:          v.Msg,
			CreatedTime:  v.CreatedTime,
		})
	}

	return info, size, nil

}

func (m *logModel) LoginLogIndex(ctx context.Context, in *LoginLogFilter) ([]*SysLoginLog, int64, error) {
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	filter := LoginLogFilter{
		Page:          &page,
		IpAddr:        in.IpAddr,
		LoginLocation: in.LoginLocation,
		Data: &DateRange{
			Start: in.Data.Start,
			End:   in.Data.End,
		},
	}

	size, err := m.GetLoginLogCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	resp, err := m.FindLoginLogByFilter(ctx, filter, &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size})
	if err != nil {
		return nil, 0, err
	}

	info := make([]*SysLoginLog, 0, len(resp))
	for _, v := range resp {
		info = append(info, &SysLoginLog{
			Id:            v.Id,
			Uid:           v.Uid,
			UserName:      v.UserName,
			IpAddr:        v.IpAddr,
			LoginLocation: v.LoginLocation,
			Browser:       v.Browser,
			Os:            v.Os,
			Code:          v.Code,
			Msg:           v.Msg,
			CreatedTime:   v.CreatedTime,
		})
	}

	return info, size, nil

}
