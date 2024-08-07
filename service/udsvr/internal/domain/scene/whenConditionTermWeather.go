package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type WeatherType = string

const (
	WeatherTypeTemp     = "temp"     //温度
	WeatherTypeHumidity = "humidity" //湿度
)

// TermProperty 物模型类型 属性
type TermWeather struct {
	Type     WeatherType `json:"type"`     //天气的类型
	TermType CmpType     `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	Values   []string    `json:"values"`   //条件值 参数根据动态条件类型会有多个参数
}

func (c *TermWeather) Validate(repo ValidateRepo) error {
	if c == nil {
		return nil
	}
	if !utils.SliceIn(c.Type, WeatherTypeTemp, WeatherTypeHumidity) {
		return errors.Parameter.AddMsg("天气类型只支持温度和湿度")
	}
	return c.TermType.Validate(c.Values)
}

func (c *TermWeather) IsHit(ctx context.Context, columnType TermColumnType, repo TermRepo) bool {
	return true
}
