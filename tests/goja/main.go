package main

import (
	"fmt"
	"github.com/dop251/goja"
)

func main() {
	r := goja.New()
	r.Set("getVal", func(in goja.FunctionCall) goja.Value {
		var ret []map[string]any
		return r.ToValue(ret)
	})
	v, err := r.RunString(`
	let d = getVal();
	function sum(a,b){
	if (d.length==0){
		return d.length;
	}
	return 9999;
}
`)
	fmt.Println(v, err)
	var sum func(int, int) int
	err = r.ExportTo(r.Get("sum"), &sum)
	if err != nil {
		panic(err)
	}
	s := sum(40, 2)
	fmt.Println(s) // note, _this_ value in the function will be undefined.
}
