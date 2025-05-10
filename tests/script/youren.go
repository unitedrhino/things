package script

import (
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"github.com/tidwall/gjson"
	"strings"
)

type Req struct {
	Params struct {
		Dir   string `json:"dir"`
		Id    string `json:"id"`
		RData []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
			Err   string `json:"err"`
		} `json:"r_data"`
	} `json:"params"`
}

var j = `{
    "params": {
        "dir": "up",
        "id": "02102925031500024611",
        "r_data": [
            {
                "name": "365.lux",
                "value": "792",
                "err": "0"
            }
        ]
    }
}`

func Run() {
	var req msgThing.Req
	rData := gjson.Get(j, "params.r_data")
	arr := rData.Array()
	var property = map[string]map[string]any{} //k为设备ID  k2 为属性id v为值
	for _, v := range arr {
		name := v.Get("name").String()
		Value := v.Get("value").String()
		dev, p, _ := strings.Cut(name, ".")
		if property[name] == nil {
			property[name] = make(map[string]any)
		}
		property[dev][p] = Value
	}
	for k, v := range property {
		req.SubDevices = append(req.SubDevices, &msgThing.SubDevice{
			ProductID:  "01S",
			DeviceName: k,
			Properties: []*deviceMsg.TimeParams{{
				Params: v,
			}},
		})
	}
}
