package main

import (
	"bytes"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sort"
	"strings"
)

type AccessInfo struct {
	Name       string // 请求名称
	Code       string // 请求名称
	Group      string // 接口组
	IsNeedAuth int64  // 是否需要认证（1是 2否）
	Desc       string // 备注
}

type AccessInfos []AccessInfo

type ApiInfo struct {
	AccessCode   string // 范围编码
	Method       string // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Route        string // 路由
	Name         string // 请求名称
	BusinessType int64  // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Desc         string // 备注
	IsAuthTenant int64  // 是否可以授权给普通租户
}

type ApiInfos []ApiInfo

func (us ApiInfos) Len() int {
	return len(us)
}

func (us ApiInfos) Less(i, j int) bool {
	if us[i].AccessCode == us[j].AccessCode {
		return us[i].Route < us[j].Route
	}
	return us[i].AccessCode < us[j].AccessCode
}

func (us ApiInfos) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

func (us AccessInfos) Len() int {
	return len(us)
}

func (us AccessInfos) Less(i, j int) bool {
	if us[i].Group == us[j].Group {
		return us[i].Name < us[j].Name
	}
	return us[i].Group < us[j].Group
}

func (us AccessInfos) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

var swagger spec.Swagger
var fileName = "./iThings.openapi.json"
var apiFile = "./api.md"
var accessFile = "./access.md"
var apiTemplate = `{AccessCode: "%v",IsAuthTenant: %v, Route: "%v", Method: %v, Name: "%v", BusinessType: %v, Desc: %v},
`

var accessTemplate = `{Name:"%v", Code:"%v", Group:"%v", IsNeedAuth:%v, Desc:"%v"},
`
var (
	groupPrefix = map[string]string{
		"物联网相关接口": "物联网",
		"系统管理new": "系统管理",
		"大屏管理":    "大屏管理",
		"大数据":     "大数据",
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
	AccessSuffix = map[int64]string{
		1: "write",
		2: "write",
		3: "write",
		4: "read",
		5: "write",
	}
	AccessName = map[int64]string{
		1: "操作权限",
		2: "操作权限",
		3: "操作权限",
		4: "读权限",
		5: "操作权限",
	}
)

func main() {
	file, err := os.ReadFile(fileName)
	logx.Must(err)
	err = swagger.UnmarshalJSON(file)
	logx.Must(err)
	var (
		apis          ApiInfos
		access        AccessInfos
		accessCodeSet = map[string]struct{}{}
	)
	for p, path := range swagger.Paths.Paths {
		var (
			Route        = p
			Method       string
			Name         string
			BusinessType int64 = 5
			AccessCode   string
			IsNeedAuth   int64 = def.True
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

		for k, v := range BusinessTypeSuffix {
			if strings.HasSuffix(p, k) {
				BusinessType = v
				break
			}
		}

		//AccessCode = tags[len(tags)-1]
		var routes = strings.Split(Route, "/")
		routes = routes[3:6]
		routes = append(routes, AccessSuffix[BusinessType])
		AccessCode = utils.UderscoreToLowerCamelCase(strings.ReplaceAll(strings.Join(routes, "_"), "-", "_"))
		{
			apis = append(apis, ApiInfo{
				AccessCode:   AccessCode,
				Method:       Method,
				Route:        Route,
				Name:         Name,
				BusinessType: BusinessType,
				Desc:         Desc,
				IsAuthTenant: IsNeedAuth,
			})
			//line := fmt.Sprintf(apiTemplate, AccessCode, IsNeedAuth, Route, Method, Name, BusinessType, Desc)
			//_, err := apis.WriteString(line)
			//logx.Must(err)
		}
		if _, ok := accessCodeSet[AccessCode]; !ok { //只处理第一次
			accessCodeSet[AccessCode] = struct{}{}
			tag := opt.Tags[0]
			tags := strings.Split(tag, "/")
			tags = tags[2:]
			group := tags[0]
			//if len(tags) > 1 {
			//	group = strings.Join(tags[0:len(tags)-1], "-")
			//}
			if len(tags) > 1 {
				tags = tags[1:]
			}
			name := strings.Join(tags, "") + AccessName[BusinessType]

			access = append(access, AccessInfo{
				Name:       name,
				Code:       AccessCode,
				Group:      group,
				IsNeedAuth: def.False,
				Desc:       "",
			})
			//line := fmt.Sprintf(accessTemplate, name, AccessCode, group, 1, "")
			//_, err := access.WriteString(line)
			//logx.Must(err)
		}
	}
	{
		sort.Sort(apis)
		f, err := os.Create(apiFile)
		logx.Must(err)
		defer f.Close()
		var bs bytes.Buffer
		for _, v := range apis {
			line := fmt.Sprintf(apiTemplate, v.AccessCode, v.IsAuthTenant, v.Route, v.Method, v.Name, v.BusinessType, v.Desc)
			_, err := bs.WriteString(line)
			logx.Must(err)
		}
		_, err = f.Write(bs.Bytes())
		logx.Must(err)
	}
	{
		sort.Sort(access)
		f, err := os.Create(accessFile)
		logx.Must(err)
		defer f.Close()
		var bs bytes.Buffer
		for _, v := range access {
			line := fmt.Sprintf(accessTemplate, v.Name, v.Code, v.Group, v.IsNeedAuth, v.Desc)
			_, err := bs.WriteString(line)
			logx.Must(err)
		}

		_, err = f.Write(bs.Bytes())
		logx.Must(err)
	}

}
