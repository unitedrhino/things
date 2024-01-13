package main

import (
	"bytes"
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

var swagger spec.Swagger
var fileName = "./iThings.openapi.json"
var outFile = "./out.md"
var template = `{ModuleCode: %s,IsNeedAuth: %v, Route: "%v", Method: %v, Name: "%v", BusinessType: %v, Desc: %v, Group: "%v"},
`
var (
	groupPrefix = map[string]string{
		"iThings-apisvr/物联网相关接口": "def.ModuleThings",
		"iThings-apisvr/系统管理new": "def.ModuleTenantSystemManage",
		"iThings-apisvr/大屏管理":    "def.ModuleView",
	}
	BusinessTypeSuffix = map[string]int64{
		"create":       1,
		"mulit-create": 1,
		"update":       2,
		"mulit-update": 2,
		"delete":       3,
		"index":        4,
		"read":         4,
	}
)

func main() {
	file, err := os.ReadFile(fileName)
	logx.Must(err)
	err = swagger.UnmarshalJSON(file)
	logx.Must(err)
	var outPut bytes.Buffer
	for p, path := range swagger.Paths.Paths {
		var (
			ModuleCode   string
			Route        = p
			Method       string
			Name         string
			BusinessType int64 = 5
			Group        string
			IsNeedAuth   int64 = 1
			Desc         string
			opt          *spec.Operation
		)
		if path.Get != nil {
			opt = path.Get
			Method = "http.MethodGet"
		} else if path.Post != nil {
			opt = path.Post
			Method = "http.MethodPost"
		}
		Name = opt.Summary
		Desc = "`" + opt.Description + "`"
		if len(opt.Tags) == 0 {
			logx.Infof("tags is nil")
			continue
		}

		tag := opt.Tags[0]
		for k, v := range groupPrefix {
			if strings.HasPrefix(tag, k) {
				ModuleCode = v
				break
			}
		}
		for k, v := range BusinessTypeSuffix {
			if strings.HasSuffix(p, k) {
				BusinessType = v
				break
			}
		}
		tags := strings.Split(tag, "/")
		Group = tags[len(tags)-1]
		line := fmt.Sprintf(template, ModuleCode, IsNeedAuth, Route, Method, Name, BusinessType, Desc, Group)
		_, err := outPut.WriteString(line)
		logx.Must(err)
	}
	f, err := os.Create(outFile)
	logx.Must(err)
	defer f.Close()
	_, err = f.Write(outPut.Bytes())
	logx.Must(err)
}
