package utils

import (
	"reflect"
)

//返回修改的数量
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

//返回是否修改
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
