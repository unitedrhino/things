package protocol

import "gitee.com/i-Things/share/errors"

const (
	CodeIThings = "iThings" //默认协议
)

type ConfigFields []*ConfigField

type ConfigField struct {
	ID         int64  `json:"id"`
	Group      string `json:"group"`      //分组(不传为默认)
	Key        string `json:"key"`        //配置文件里的关键字
	Label      string `json:"label"`      //展示名称
	IsRequired bool   `json:"isRequired"` //是否必填
	Sort       int64  `json:"sort"`       //排序
}

type ConfigInfos []*ConfigInfo

// 自定义协议的配置表
type ConfigInfo struct {
	ID     int64             `json:"id"`
	Config map[string]string `json:"config"` //协议配置内容,key要从protocolInfo的ConfigKeys里拿
	Desc   string            `json:"desc"`   // 描述
}

func Check(fields ConfigFields, infos ConfigInfos) error {
	if len(fields) == 0 {
		return nil
	}
	for _, info := range infos {
		for _, field := range fields {
			v := info.Config[field.Key]
			if len(v) == 0 && field.IsRequired { //如果是必传的但是没有传,则需要报错
				return errors.Parameter.AddMsgf("%s 必填", field.Label)
			}
		}
	}
	return nil
}
func (c ConfigInfos) ToPubStu() (ret []map[string]string) {
	for _, info := range c {
		ret = append(ret, info.Config)
	}
	return ret
}
