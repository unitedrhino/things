package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)


func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		data,_:= ioutil.ReadAll(c.Request.Body)
		fmt.Printf("data=%s|url=%s|uri=%s\n",string(data),c.Request.URL.String(),c.Request.RequestURI)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}