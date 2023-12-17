package stores

import (
	"fmt"
	"github.com/i-Things/things/shared/def"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type AuthQueryClause struct {
	Field        *schema.Field
	GetAuthIDs   GetAuthIDs
	AuthDataType def.AuthDataType
}

func (sd AuthQueryClause) Name() string {
	return ""
}

func (sd AuthQueryClause) Build(clause.Builder) {
}

func (sd AuthQueryClause) MergeClause(*clause.Clause) {
}

const AuthModify = "authModify:%v"

func (sd AuthQueryClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, sd.AuthDataType)
}
func (sd AuthQueryClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口

	ids, isRoot, err := sd.GetAuthIDs(stmt)
	if err != nil {
		stmt.Error = err
		return
	}
	if isRoot { //root 权限不用管
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
