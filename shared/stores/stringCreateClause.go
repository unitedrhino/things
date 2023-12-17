package stores

import (
	"github.com/i-Things/things/shared/def"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type StringCreateClause[keyT ~string] struct {
	Field        *schema.Field
	GetValues    GetValues
	AuthDataType def.AuthDataType
}

func (sd StringCreateClause[keyT]) Name() string {
	return ""
}

func (sd StringCreateClause[keyT]) Build(clause.Builder) {
}

func (sd StringCreateClause[keyT]) MergeClause(*clause.Clause) {
}

func (sd StringCreateClause[keyT]) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	values, isRoot, _, err := sd.GetValues(stmt)
	if err != nil {
		stmt.Error = err
		return
	}
	if len(values) > 0 {
		destV := reflect.ValueOf(stmt.Dest)
		if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
			for i := 0; i < destV.Len(); i++ {
				dest := destV.Index(i)
				field := dest.Elem().FieldByName(sd.Field.Name)
				if isRoot && !field.IsZero() { //只有root权限的租户可以设置为其他租户
					continue
				}
				field.IsZero()
				field.Set(reflect.ValueOf(values[0]))
			}
			return
		}

		field := destV.Elem().FieldByName(sd.Field.Name)
		if isRoot && !field.IsZero() { //只有root权限的租户可以设置为其他租户
			return
		}
		field.Set(reflect.ValueOf(values[0]))
	}
}
