package stores

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type AreaID int64

func (t AreaID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	authIDs, err := caches.GatherUserAuthAreaIDs(ctx)
	if err != nil {
		stmt.Error = err
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	return
	if !utils.SliceIn(int64(t), authIDs...) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("区域权限不足")
	}
	return
}
func (t *AreaID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := AreaID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t AreaID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t AreaID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Select}}
}
func (t AreaID) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Update}}
}

func (t AreaID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Create}}
}

func (t AreaID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: Delete}}
}

type AreaClause struct {
	clauseInterface
	Field *schema.Field
	T     AreaID
	Opt   Opt
}

func (sd AreaClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, "areaID")
}

func (sd AreaClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	ids, err := caches.GatherUserAuthAreaIDs(stmt.Context)
	if err != nil {
		stmt.Error = err
		return
	}
	switch sd.Opt {
	case Create:
	case Update, Delete, Select:
		if len(ids) == 0 { //root 权限不用管
			return
		}
		if _, ok := stmt.Clauses[sd.GenAuthKey()]; !ok {
			if c, ok := stmt.Clauses["WHERE"]; ok {
				if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
					for _, expr := range where.Exprs {
						if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
							where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
							c.Expression = where
							stmt.Clauses["WHERE"] = c
							break
						}
					}
				}
			}
			var values []any
			for _, v := range ids {
				values = append(values, v)
			}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}
