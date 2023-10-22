package main

import (
	"fmt"
	"github.com/dop251/goja"
)

func main() {
	r := goja.New()
	v, err := r.RunString(`
	let d = 1+2+3;
	function sum(a,b){
	return a+b+d;
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
