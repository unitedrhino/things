package utils

import "reflect"


/*
@in src 赋值的数据源
@in dst 赋值对象的结构体
@out dst类型的结构体
*/
func Convert(src interface{},dst interface{}) interface{}{


	srcType := reflect.TypeOf(src)  //获取type
	dstType := reflect.TypeOf(dst)

	srcEl := reflect.ValueOf(src).Elem() //获取value
	dstEl := reflect.ValueOf(dst).Elem()
	//双循环，对相同名字对字段进行赋值
	for i := 0; i < srcType.NumField(); i++ {
		for j:=0;j<dstType.NumField();j++ {
			if srcType.Field(i).Name == dstType.Field(j).Name {
				dstEl.Field(i).Set(srcEl.Field(j))
			}
		}
	}
	return dst
}