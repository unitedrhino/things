package store

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
	"reflect"
	"strings"
	"time"
)

func prepareValues(values []any, columnTypes []*sql.ColumnType, columns []string) {
	if len(columnTypes) > 0 {
		for idx, columnType := range columnTypes {
			if columnType.ScanType() != nil {
				values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
			} else {
				values[idx] = new(any)
			}
		}
	} else {
		for idx := range columns {
			values[idx] = new(any)
		}
	}
}

func scanIntoMap(mapValue map[string]any, values []any, columns []string) {
	for idx, column := range columns {
		if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(values[idx]))); reflectValue.IsValid() {
			mapValue[column] = reflectValue.Interface()
			if valuer, ok := mapValue[column].(driver.Valuer); ok {
				mapValue[column], _ = valuer.Value()
			} else if b, ok := mapValue[column].(sql.RawBytes); ok {
				mapValue[column] = string(b)
			}
		} else {
			mapValue[column] = nil
		}
	}
}

func Scan(rows *sql.Rows, Dest any) error {
	columns, _ := rows.Columns()
	values := make([]any, len(columns))

	switch dest := Dest.(type) {
	case map[string]any, *map[string]any:
		if rows.Next() {
			columnTypes, _ := rows.ColumnTypes()
			prepareValues(values, columnTypes, columns)
			if err := rows.Scan(values...); err != nil {
				return err
			}

			mapValue, ok := dest.(map[string]any)
			if !ok {
				if v, ok := dest.(*map[string]any); ok {
					mapValue = *v
				}
			}
			scanIntoMap(mapValue, values, columns)
		}
	case *[]map[string]any:
		columnTypes, _ := rows.ColumnTypes()
		for rows.Next() {
			prepareValues(values, columnTypes, columns)
			if err := rows.Scan(values...); err != nil {
				return err
			}

			mapValue := map[string]any{}
			scanIntoMap(mapValue, values, columns)
			*dest = append(*dest, mapValue)
		}
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
		*float32, *float64,
		*bool, *string, *time.Time,
		*sql.NullInt32, *sql.NullInt64, *sql.NullFloat64,
		*sql.NullBool, *sql.NullString, *sql.NullTime:
		for rows.Next() {
			if err := rows.Scan(dest); err != nil {
				return err
			}
		}
	default:
		return errors.Database.AddMsgf("not support type:%v", reflect.TypeOf(dest))
	}
	return nil
}

func GetTdType(define schema.Define) string {
	switch define.Type {
	case schema.DataTypeBool:
		return "BOOL"
	case schema.DataTypeInt:
		return "BIGINT"
	case schema.DataTypeString:
		return "BINARY(5000)"
	case schema.DataTypeStruct:
		return "BINARY(5000)"
	case schema.DataTypeFloat:
		return "DOUBLE"
	case schema.DataTypeTimestamp:
		return "TIMESTAMP"
	case schema.DataTypeArray:
		return "BINARY(5000)"
	case schema.DataTypeEnum:
		return "SMALLINT"
	default: //走到这里说明前面没有进行校验需要检查是否是前面有问题
		panic(fmt.Sprintf("%v not support", define.Type))
	}
}

// GenParams 返回占位符?,?,?,? 参数id名:aa,bbb,ccc 参数值列表
func GenParams(params map[string]any) (string, string, []any, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		return "", "", nil, errors.Parameter.AddMsgf("SchemaDataRepo|GenParams|params num == 0")
	}
	var (
		paramPlaceholder = strings.Repeat("?,", len(params))
		paramValList     []any //参数值列表
		paramIds         []string
	)
	//将最后一个?去除
	paramPlaceholder = paramPlaceholder[:len(paramPlaceholder)-1]
	for k, v := range params {
		paramIds = append(paramIds, "`"+k+"`")
		if _, ok := v.([]any); !ok {
			paramValList = append(paramValList, v)
		} else { //如果是数组类型,需要序列化为json
			param, err := json.Marshal(v)
			if err != nil {
				return "", "", nil, errors.System.AddDetail("param json parse failure")
			}
			paramValList = append(paramValList, param)
		}
	}
	return paramPlaceholder, strings.Join(paramIds, ","), paramValList, nil
}

func ArrayToSql[arrType any](arr []arrType) (sql string) {
	if len(arr) == 0 {
		return ""
	}
	for _, v := range arr {
		sql += "\"" + cast.ToString(v) + "\","
	}
	sql = sql[:len(sql)-1]
	return
}
