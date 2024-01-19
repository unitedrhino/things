package stores

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type TenantCode string

func (t TenantCode) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil { //系统初始化的时候会掉用这里
		expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
		return
	}
	if uc.TenantCode == "" {
		stmt.Error = errors.Parameter.AddDetail("tenantCode is empty")
		return
	}
	if t != "" && uc.TenantCode == def.TenantCodeDefault && uc.AllTenant {
		expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{uc.TenantCode}}
	return
}
func (t *TenantCode) Scan(value interface{}) error {
	ret := cast.ToString(value)
	p := TenantCode(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t TenantCode) Value() (driver.Value, error) {
	return string(t), nil
}

func (t TenantCode) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeClause{Field: f, T: t, Opt: Select}}
}

func (t TenantCode) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeClause{Field: f, T: t, Opt: Update}}
}

func (t TenantCode) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeClause{Field: f, T: t, Opt: Create}}
}

func (t TenantCode) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeClause{Field: f, T: t, Opt: Delete}}
}

func (t TenantCode) GetAuthIDs(f *schema.Field) GetValues {
	return func(stmt *gorm.Statement) (authIDs []any, isRoot bool, allData bool, err error) {
		uc := ctxs.GetUserCtx(stmt.Context)
		if uc == nil {
			return nil, false, false, nil
		}
		if uc.TenantCode == def.TenantCodeDefault { //只有core租户的可以修改其他租户的租户号
			isRoot = true
		}
		return []any{TenantCode(uc.TenantCode)}, isRoot, uc.AllTenant, nil
	}
}

type TenantCodeClause struct {
	clauseInterface
	Field *schema.Field
	T     TenantCode
	Opt   Opt
}

func (sd TenantCodeClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, "tenantCode")
}

func (sd TenantCodeClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	var (
		isRoot     bool
		tenantCode string
		allTenant  bool
	)

	uc := ctxs.GetUserCtx(stmt.Context)
	if uc != nil {
		allTenant = uc.AllTenant
		if uc.TenantCode == def.TenantCodeDefault { //只有core租户的可以修改其他租户的租户号
			isRoot = true
		}
		if uc.TenantCode != "" {
			tenantCode = uc.TenantCode
		}

	}

	switch sd.Opt {
	case Create:
		if uc != nil {
			destV := reflect.ValueOf(stmt.Dest)
			if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
				for i := 0; i < destV.Len(); i++ {
					dest := destV.Index(i)
					field := dest.Elem().FieldByName(sd.Field.Name)
					if tenantCode != "" && !field.IsZero() { //只有root权限的租户可以设置为其他租户
						continue
					}
					var v TenantCode
					v = TenantCode(tenantCode)
					field.Set(reflect.ValueOf(v))
				}
				return
			}
			field := destV.Elem().FieldByName(sd.Field.Name)
			if tenantCode != "" && !field.IsZero() { //只有root权限的租户可以设置为其他租户
				return
			}
			var v TenantCode
			v = TenantCode(tenantCode)
			field.Set(reflect.ValueOf(v))
		}
	case Update, Delete, Select:
		if isRoot && allTenant {
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
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: []any{tenantCode}},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}
