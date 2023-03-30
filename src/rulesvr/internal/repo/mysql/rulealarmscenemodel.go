package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleAlarmSceneModel = (*customRuleAlarmSceneModel)(nil)

type (
	// RuleAlarmSceneModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmSceneModel.
	RuleAlarmSceneModel interface {
		ruleAlarmSceneModel
		FindByFilter(ctx context.Context, filter alarm.SceneFilter, page *def.PageInfo) ([]*RuleAlarmScene, error)
		CountByFilter(ctx context.Context, filter alarm.SceneFilter) (size int64, err error)
		DeleteByFilter(ctx context.Context, filter alarm.SceneFilter) error
		InsertMulti(ctx context.Context, alarmID int64, sceneIDs []int64) error
	}

	customRuleAlarmSceneModel struct {
		*defaultRuleAlarmSceneModel
	}
)

// NewRuleAlarmSceneModel returns a model for the database table.
func NewRuleAlarmSceneModel(conn sqlx.SqlConn) RuleAlarmSceneModel {
	return &customRuleAlarmSceneModel{
		defaultRuleAlarmSceneModel: newRuleAlarmSceneModel(conn),
	}
}
func (c customRuleAlarmSceneModel) FmtSql(sql sq.SelectBuilder, f alarm.SceneFilter) sq.SelectBuilder {
	if f.SceneID != 0 {
		sql = sql.Where("sceneID=?", f.SceneID)
	}
	if f.AlarmID != 0 {
		sql = sql.Where("alarmID=?", f.AlarmID)
	}
	return sql
}

func (c customRuleAlarmSceneModel) FindByFilter(ctx context.Context, filter alarm.SceneFilter, page *def.PageInfo) ([]*RuleAlarmScene, error) {
	//TODO implement me
	panic("implement me")
}

func (c customRuleAlarmSceneModel) CountByFilter(ctx context.Context, filter alarm.SceneFilter) (size int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (c customRuleAlarmSceneModel) DeleteByFilter(ctx context.Context, f alarm.SceneFilter) error {
	sql := sq.Delete(c.table)
	if f.SceneID != 0 {
		sql = sql.Where("sceneID=?", f.SceneID)
	}
	if f.AlarmID != 0 {
		sql = sql.Where("alarmID=?", f.AlarmID)
	}
	query, arg, err := sql.ToSql()
	if err != nil {
		return err
	}
	_, err = c.conn.ExecCtx(ctx, query, arg...)
	return err
}

func (c customRuleAlarmSceneModel) InsertMulti(ctx context.Context, alarmID int64, sceneIDs []int64) error {
	sql := sq.Insert(c.table).Columns(ruleAlarmSceneRowsExpectAutoSet)
	for _, sceneID := range sceneIDs {
		sql = sql.Values(alarmID, sceneID)
	}
	query, arg, err := sql.ToSql()
	if err != nil {
		return err
	}
	_, err = c.conn.ExecCtx(ctx, query, arg...)
	return err
}
