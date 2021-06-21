package model

type (
	DeviceReq struct {
		Method      string                 `json:"method"`             //操作方法
		ClientToken string                 `json:"clientToken"`        //方便排查随机数
		Params      map[string]interface{} `json:"params,omitempty"`   //参数列表
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		EventID     string                 `json:"eventId,omitempty"`  //事件的 Id，在数据模板事件中定义。
		ActionID    string                 `json:"actionId,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
		Timestamp   int64                  `json:"timestamp,omitempty"`
		Showmeta    int64                  `json:"showmeta,omitempty"` //标识回复消息是否带 metadata，缺省为0表示不返回 metadata
		Type        string                 `json:"type,omitempty"`     //	表示获取什么类型的信息。report:表示设备上报的信息 info:信息 alert:告警 fault:故障
	}
	DeviceResp struct {
		Method      string                 `json:"method"`             //操作方法
		ClientToken string                 `json:"clientToken"`        //方便排查随机数
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		Code        int64                  `json:"code"`               //状态码
		Status      string                 `json:"status,omitempty"`   //返回信息
		Type        string                 `json:"type,omitempty"`     //	表示什么类型的信息。report:表示设备上报的信息
		Data        map[string]interface{} `json:"data,omitempty"`     //返回具体设备上报的最新数据内容
		Response    map[string]interface{} `json:"response,omitempty"` //设备行为中定义的返回参数，设备行为执行成功后，向云端返回执行结果
	}
)
