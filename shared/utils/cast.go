package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"golang.org/x/exp/constraints"
	"reflect"
	"time"
)

/*
@in src 赋值的数据源
@in dst 赋值对象的结构体
@out dst类型的结构体
*/
func Convert(src any, dst any) any {

	srcType := reflect.TypeOf(src) //获取type
	dstType := reflect.TypeOf(dst)

	srcEl := reflect.ValueOf(src).Elem() //获取value
	dstEl := reflect.ValueOf(dst).Elem()
	//双循环，对相同名字对字段进行赋值
	for i := 0; i < srcType.NumField(); i++ {
		for j := 0; j < dstType.NumField(); j++ {
			if srcType.Field(i).Name == dstType.Field(j).Name {
				dstEl.Field(i).Set(srcEl.Field(j))
			}
		}
	}
	return dst
}

func ToNullString(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}
func ToRpcNullString(val *string) *wrappers.StringValue {
	if val != nil {
		return &wrappers.StringValue{
			Value: *val,
		}
	}
	return nil
}

func ToRpcNullDouble(val *float64) *wrappers.DoubleValue {
	if val != nil {
		return &wrappers.DoubleValue{
			Value: *val,
		}
	}
	return nil
}

var empty = time.Time{}

func TimeToInt64(t time.Time) int64 {
	if t == empty {
		return 0
	}
	return t.Unix()
}
func SetToSlice[t constraints.Ordered](in map[t]struct{}) (ret []t) {
	for k, _ := range in {
		ret = append(ret, k)
	}
	return
}

func AnyToNullString(in any) sql.NullString {
	if in == nil || IsNil(in) {
		return sql.NullString{}
	}
	str, err := json.Marshal(in)
	if err != nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(str), Valid: true}
}
func SqlNullStringToAny(in sql.NullString, ret any) error {
	if in.Valid == false {
		return nil
	}
	err := json.Unmarshal([]byte(in.String), ret)
	return err
}
