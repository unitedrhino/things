package utils

import (
	"database/sql"
	"reflect"
)

// 返回修改的数量
func SetVal(key any, src, dst any) (num int) {
	switch key.(type) {
	case string:
		setVal(key.(string), src, dst)
	case []string:
		for _, v := range key.([]string) {
			num += setVal(v, src, dst)
		}
	default:
		panic("SetVal not support key type")
	}
	return num
}

// 返回是否修改
func setVal(key string, src, dst any) (num int) {
	srcKey := reflect.ValueOf(src).FieldByName(key)
	if srcKey.IsValid() == true {
		if dstV := reflect.ValueOf(dst).Elem().FieldByName(key); dstV.CanSet() == true && dstV.Type().Kind() == srcKey.Kind() {
			dstV.Set(srcKey)
			return 1
		}
	}
	return 0
}

func SetVals(src any, dst any) (num int) {
	srcV := reflect.ValueOf(src).Elem()
	dstV := reflect.ValueOf(dst).Elem()
	srcT := srcV.Type()
	dstT := dstV.Type()
	for i := 0; i < srcT.NumField(); i++ {
		for j := 0; j < dstT.NumField(); j++ {
			if srcT.Field(i).Name == dstT.Field(j).Name && srcT.Field(i).Type.AssignableTo(dstT.Field(j).Type) {
				dstV.Field(j).Set(srcV.Field(i))
				num++
			}
		}
	}
	return num
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if SliceIn(vi.Kind(), reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan) {
		return vi.IsNil()
	}
	return false
}

// 解析数据库 表结构字段， ptr：数据表结构， excludeFields：排除的字段名， 返回值：解析出的表结构字段值列表
func ReflectFields(ptr any, excludeFields []string) []any {
	exclude := make(map[string]bool, 0)
	for _, v := range excludeFields {
		exclude[v] = true
	}

	fieldValues := make([]any, 0)
	elem := reflect.ValueOf(ptr).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tag := elem.Type().Field(i).Tag.Get("db")
		if _, ok := exclude[tag]; ok {
			continue
		}
		value := elem.Field(i).Interface()

		// 特殊处理的类型
		switch value.(type) {
		case sql.NullTime:
			value = value.(sql.NullTime)
		default:
		}
		fieldValues = append(fieldValues, value)
	}
	return fieldValues
}
