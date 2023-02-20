package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleSceneInfoModel = (*customRuleSceneInfoModel)(nil)

type (
	// RuleSceneInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleSceneInfoModel.
	RuleSceneInfoModel interface {
		scene.Repo
	}

	customRuleSceneInfoModel struct {
		repo *defaultRuleSceneInfoModel
	}
)

// NewRuleSceneInfoModel returns a model for the database table.
func NewRuleSceneInfoModel(conn sqlx.SqlConn) RuleSceneInfoModel {
	return &customRuleSceneInfoModel{
		repo: newRuleSceneInfoModel(conn),
	}
}

func (c customRuleSceneInfoModel) Insert(ctx context.Context, info *scene.Info) error {
	_, err := c.repo.Insert(ctx, ToScenePo(info))
	return err
}

func (c customRuleSceneInfoModel) Delete(ctx context.Context, id int64) error {
	return c.repo.Delete(ctx, id)
}
func (c customRuleSceneInfoModel) Update(ctx context.Context, info *scene.Info) error {
	err := c.repo.Update(ctx, ToScenePo(info))
	return err
}

func (c customRuleSceneInfoModel) FindOne(ctx context.Context, id int64) (*scene.Info, error) {
	info, err := c.repo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToSceneDo(info), nil
}
func (c customRuleSceneInfoModel) FindOneByName(ctx context.Context, name string) (*scene.Info, error) {
	info, err := c.repo.FindOneByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return ToSceneDo(info), nil
}

func (c customRuleSceneInfoModel) FmtFilter(filter scene.InfoFilter, sql sq.SelectBuilder) sq.SelectBuilder {
	if filter.Name != "" {
		sql = sql.Where("name like ?", "%"+filter.Name+"%")
	}
	return sql
}

func (c customRuleSceneInfoModel) FindByFilter(ctx context.Context, filter scene.InfoFilter, page *def.PageInfo) ([]*scene.Info, error) {
	var poList []*RuleSceneInfo
	sql := sq.Select(ruleSceneInfoRows).From(c.repo.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = c.FmtFilter(filter, sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = c.repo.conn.QueryRowsCtx(ctx, &poList, query, arg...)
	switch err {
	case nil:
		var resp []*scene.Info
		for _, v := range poList {
			resp = append(resp, ToSceneDo(v))
		}
		return resp, nil
	default:
		return nil, err
	}
}

func (c customRuleSceneInfoModel) CountByFilter(ctx context.Context, filter scene.InfoFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(c.repo.table)
	sql = c.FmtFilter(filter, sql)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = c.repo.conn.QueryRowCtx(ctx, &size, query, arg...)
	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}
