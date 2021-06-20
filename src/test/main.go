package main

import (
	"fmt"
	"strings"
)

func main()  {
	ip:= "192.168.1.2:803"
	addr := string([]byte(ip)[0:strings.LastIndex(ip,":")])
	fmt.Printf("ip=%s|addr=%s\n",ip,addr)
}
