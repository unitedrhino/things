package stores

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type StringQueryClause struct {
	Field     *schema.Field
	GetValues GetValues
	Key       string
}

func (sd StringQueryClause) Name() string {
	return ""
}

func (sd StringQueryClause) Build(clause.Builder) {
}

func (sd StringQueryClause) MergeClause(*clause.Clause) {
}

func (sd StringQueryClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(AuthModify, sd.Key)
}
func (sd StringQueryClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口

	values, isRoot, isAllData, err := sd.GetValues(stmt)
	if err != nil {
		stmt.Error = err
		return
	}
	if isRoot && isAllData { //root 权限不用管
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
			clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
		}})
		stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
	}
}
