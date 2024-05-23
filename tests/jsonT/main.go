package main

import "fmt"

func main() {
	str := "{\n    \"bool\": true,\n    \"int64\": ${groupId},\n    \"float64\": 12.34234,\n    \"string\": \"stringBody\",\n    \"position\": {\n      \"longitude\": 0.012334992177784441,\n      \"latitude\": 50.045040130615234\n    }\n  }"
	fmt.Println(str)
}
