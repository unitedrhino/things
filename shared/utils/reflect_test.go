package utils

import (
	"fmt"
	"testing"
)

type Test1 struct {
	INT32 int32
	INT64 int64
}

type Test2 struct {
	INT32 int64
	INT64 int64
}

func TestSetVals(t *testing.T) {
	test1 := Test1{INT32: 998}
	test2 := Test2{}
	SetVals(&test1, &test2)
	fmt.Printf("test1:%+v|test2:%+v\n", test1, test2)
}
func TestSetVal(t *testing.T) {
	test1 := Test1{INT32: 998}
	test2 := Test2{}
	SetVal("INT32", test1, &test2)
	fmt.Printf("test1:%+v|test2:%+v\n", test1, test2)
}
