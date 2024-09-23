package scene

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/observerly/dusk/pkg/dusk"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TimeType = string

const (
	TimeTypeSys               = "sys"      //系统时间
	TimeTypeSunRises ExecType = "sunRises" //太阳升起
	TimeTypeSunSet   ExecType = "sunSet"   //太阳落下
)

// TermTime 时间执行条件
type TermTime struct {
	Type     TimeType `json:"type"`     //时间的类型
	TermType CmpType  `json:"termType"` //动态条件类型  eq: 相等  not:不相等  btw:在xx之间  gt: 大于  gte:大于等于 lt:小于  lte:小于等于   in:在xx值之间
	/*
			条件值 参数根据动态条件类型会有多个参数
		如果是sys类型 则为触发时间和values中的时间进行比较,比如 7点-8点 ,
		如果是太阳起和落下,如落下后30分钟内则为 now btw  0, 30*60
	*/
	Values []string `json:"values"`
}

func (c *TermTime) Validate(repo CheckRepo) error {
	if c == nil {
		return nil
	}
	if !utils.SliceIn(c.Type, TimeTypeSys, TimeTypeSunRises, TimeTypeSunSet) {
		return errors.Parameter.AddMsg("天气类型只支持温度和湿度")
	}
	switch c.Type {
	case TimeTypeSys:
		for _, v := range c.Values {
			vv := cast.ToInt64(v)
			if vv < 0 || vv > 24*60*60 {
				return errors.Parameter.AddMsg("时间范围只能在0到24小时之间")
			}
		}
	case TimeTypeSunRises, TimeTypeSunSet:
		pi, err := repo.ProjectCache.GetData(repo.Ctx, repo.Info.ProjectID)
		if err != nil {
			return err
		}
		if pi.Position == nil || pi.Position.Latitude == 0 || pi.Position.Longitude == 0 {
			return errors.Parameter.AddMsg("需要填写地理位置才可以使用日出日落执行条件")
		}
		for _, v := range c.Values {
			vv := cast.ToInt64(v)
			if vv > 3*60*60 {
				return errors.Parameter.AddMsg("最晚只能三个小时后")
			}
			if vv < (-3 * 60 * 60) {
				return errors.Parameter.AddMsg("最早只能提前三个小时")
			}
		}
	}
	return c.TermType.Validate(c.Values)
}

func (c *TermTime) IsHit(ctx context.Context, repo CheckRepo) bool {
	now := utils.TimeToDaySec(time.Now())
	switch c.Type {
	case TimeTypeSys:
		return c.TermType.IsHit(schema.DataTypeInt, now, c.Values)
	case TimeTypeSunRises, TimeTypeSunSet:
		pi, err := repo.ProjectCache.GetData(repo.Ctx, repo.Info.ProjectID)
		if err != nil || pi.Position == nil || pi.Position.Latitude == 0 || pi.Position.Longitude == 0 {
			return false
		}
		twilight, _, err := dusk.GetLocalCivilTwilight(time.Now(), pi.Position.Longitude, pi.Position.Latitude, 0)
		if err != nil {
			logx.WithContext(repo.Ctx).Error(err)
			return false
		}
		var sunTime int64
		switch c.Type {
		case TimeTypeSunRises:
			sunTime = utils.TimeToDaySec(twilight.Until)
		case TimeTypeSunSet:
			sunTime = utils.TimeToDaySec(twilight.From)
		}
		var values []string
		for _, v := range c.Values {
			values = append(values, cast.ToString(cast.ToInt64(v)+sunTime))
		}
		return c.TermType.IsHit(schema.DataTypeInt, now, values)
	}
	return true
}
