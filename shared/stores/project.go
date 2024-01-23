package stores

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type ProjectID int64

func (t ProjectID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	authIDs, err := caches.GatherUserAuthProjectIDs(ctx)
	if err != nil {
		stmt.Error = err
		return
	}
	uc := ctxs.GetUserCtx(ctx)
	if t != 0 && (uc == nil || authIDs == nil || uc.AllProject) { //root 权限不用管
		expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	} else {
		t = ProjectID(uc.ProjectID)
		expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}
	}
	if !(uc == nil || authIDs == nil || uc.AllProject) && !utils.SliceIn(int64(t), authIDs...) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("项目权限不足")
	}

	return
}
func (t *ProjectID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := ProjectID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t ProjectID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t ProjectID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Select}}
}
func (t ProjectID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Create}}
}

func (t ProjectID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Delete}}
}
func (t ProjectID) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ProjectClause{Field: f, T: t, Opt: Update}}
}

type ProjectClause struct {
	clauseInterface
	Field *schema.Field
	T     ProjectID
	Opt   Opt
}

func (sd ProjectClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, "projectID")
}

func (sd ProjectClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	ids, err := caches.GatherUserAuthProjectIDs(stmt.Context)
	if err != nil {
		stmt.Error = err
		return
	}
	uc := ctxs.GetUserCtx(stmt.Context)
	switch sd.Opt {
	case Create:
		if uc != nil {
			destV := reflect.ValueOf(stmt.Dest)
			if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
				for i := 0; i < destV.Len(); i++ {
					dest := destV.Index(i)
					field := dest.Elem().FieldByName(sd.Field.Name)
					if len(ids) == 0 && !field.IsZero() { //只有root权限的租户可以设置为其他租户
						continue
					}
					var v ProjectID
					v = ProjectID(uc.ProjectID)
					field.Set(reflect.ValueOf(v))
				}
				return
			}
			field := destV.Elem().FieldByName(sd.Field.Name)
			if len(ids) == 0 && !field.IsZero() { //只有root权限的租户可以设置为其他租户
				return
			}
			var v ProjectID
			v = ProjectID(uc.ProjectID)
			field.Set(reflect.ValueOf(v))
		}
	case Update, Delete, Select:
		if uc == nil || (len(ids) == 0 && uc.AllProject) { //root 权限不用管
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
			var values = []any{uc.ProjectID}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}
