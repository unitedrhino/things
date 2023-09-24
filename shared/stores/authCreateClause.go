package stores

import (
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type CreateClause[keyT ~int64] struct {
	Field        *schema.Field
	GetAuthIDs   GetAuthIDs
	AuthDataType def.AuthDataType
}

func (sd CreateClause[keyT]) Name() string {
	return ""
}

func (sd CreateClause[keyT]) Build(clause.Builder) {
}

func (sd CreateClause[keyT]) MergeClause(*clause.Clause) {
}

func (sd CreateClause[keyT]) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	_, isRoot, err := sd.GetAuthIDs(stmt)
	if err != nil {
		stmt.Error = err
		return
	}
	if isRoot { //root 权限不用管
		return
	}
	if sd.AuthDataType == def.AuthDataTypeProject { //设置默认值
		destV := reflect.ValueOf(stmt.Dest)
		if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
			for i := 0; i < destV.Len(); i++ {
				dest := destV.Index(i)
				field := dest.Elem().FieldByName(sd.Field.Name)
				var v keyT
				v = keyT(ctxs.GetMetaProjectID(stmt.Context))
				field.Set(reflect.ValueOf(v))
			}
			return
		}
		field := destV.Elem().FieldByName(sd.Field.Name)
		var v keyT
		v = keyT(ctxs.GetMetaProjectID(stmt.Context))
		field.Set(reflect.ValueOf(v))
	}
}
