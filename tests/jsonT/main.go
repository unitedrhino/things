package main

import (
	"encoding/base64"
	"fmt"
)

//	func main() {
//		str := "{\n    \"bool\": true,\n    \"int64\": ${groupId},\n    \"float64\": 12.34234,\n    \"string\": \"stringBody\",\n    \"position\": {\n      \"longitude\": 0.012334992177784441,\n      \"latitude\": 50.045040130615234\n    }\n  }"
//		fmt.Println(str)
//	}
func main() {
	str := base64.StdEncoding.EncodeToString([]byte{123, 42, 52, 23, 53, 64})
	v, err := base64.StdEncoding.DecodeString(str)
	fmt.Println(v, err)
}
