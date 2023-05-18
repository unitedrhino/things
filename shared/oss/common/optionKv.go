package common

import (
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/url"
)

const (
	ForbidWrite                = "X-oss-Forbid-Overwrite"
	Process                    = "x-process"
	ResponseContentDisposition = "response-content-disposition"
)

//var (
//	imageProcessParser = NewImageProcessParser()
//)

type OptionKv struct {
	httpHeader map[string]interface{}
	httpParams map[string]interface{}
}

func (o *OptionKv) ToAliYunOptions() []oss.Option {
	res := make([]oss.Option, 0)
	if forbidWrite, ok := o.httpHeader[ForbidWrite]; ok {
		res = append(res, oss.ForbidOverWrite(forbidWrite.(bool)))
		delete(o.httpHeader, ForbidWrite)
	}
	if value, ok := o.httpParams[Process]; ok {
		res = append(res, oss.Process(value.(string)))
		delete(o.httpParams, Process)
	}
	if value, ok := o.httpParams[ResponseContentDisposition]; ok {
		res = append(res, oss.ResponseContentDisposition("attachment;filename="+value.(string)))
		delete(o.httpParams, ResponseContentDisposition)
	}
	for k, v := range o.httpHeader {
		res = append(res, oss.SetHeader(k, v))
	}

	for k, v := range o.httpParams {
		res = append(res, oss.AddParam(k, v))
	}
	return res
}

func (o *OptionKv) SetHeader(k string, v interface{}) {
	if o.httpHeader == nil {
		o.httpHeader = make(map[string]interface{})
	}
	o.httpHeader[k] = v
}

func (o *OptionKv) SetHttpParams(k string, v interface{}) {
	if o.httpParams == nil {
		o.httpParams = make(map[string]interface{})
	}
	o.httpParams[k] = v
}

// 不指定x-oss-forbid-overwrite时，默认覆盖同名Object。
// 指定x-oss-forbid-overwrite为false时，表示允许覆盖同名Object。
// 指定x-oss-forbid-overwrite为true时，表示禁止覆盖同名Object，如果同名Object已存在，程序将报错。
func (o *OptionKv) IsForbidOverwrite() bool {
	if forbidWrite, ok := o.httpHeader[ForbidWrite]; ok {
		return forbidWrite.(bool)
	}
	return false
}

func (o *OptionKv) ToMinioReqParams() url.Values {
	values := url.Values{}
	if value, ok := o.httpParams[ResponseContentDisposition]; ok {
		values.Add("response-content-disposition", "attachment;filename="+value.(string))
		delete(o.httpParams, ResponseContentDisposition)
	}
	return values
}

func (o *OptionKv) CheckAndGetMinioProcess() (interface{}, bool) {
	value, ok := o.httpParams[Process]
	return value, ok
}

func (o *OptionKv) ToMinioFilePath(filePath string) string {
	if value, ok := o.httpParams[Process]; ok {
		delete(o.httpParams, Process)
		return filePath + "?x-oss-process=" + value.(string)
	}
	return filePath
}
