package protocolmanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolScriptDebugLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolScriptDebugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolScriptDebugLogic {
	return &ProtocolScriptDebugLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var PublishMsgStr = `
{ //发布消息结构体
	"topic":string,  //可选,只用于日志记录
	"handle":string, //对应 mqtt topic的第一个 thing ota config 等等
	"type":string,   //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
	"payload":bytes, //base64格式的消息内容
	"timestamp":int64, //可选,毫秒时间戳
	"productID":string,//必填
	"deviceName":string,//必填
	"explain":string,      //可选,内部使用的拓展字段
	"protocolCode":string, //可选,如果有该字段则回复的时候也会带上该字段
}`

var str = "{\"DeviceName\":\"867713079344958\",\"Handle\":\"thing\",\"Payload\":\"e1wibWV0aG9kXCI6XCJyZXBvcnRcIixcIm1zZ1Rva2VuXCI6XCI2MzVcIixcInRpbWVzdGFtcFwiOjE3Mzk2MzU0MzcwMjUsXCJzeXNcIjp7XCJub0Fza1wiOnRydWV9LFwicGFyYW1zXCI6e1wiZGV2aWNlX3RpbWVcIjpcIjE4MDI3MDc0MzUwMDBcIn19\",\"ProductID\":\"120\",\"Timestamp\":1739635437025,\"Type\":\"property\",\"protocolCode\":\"wumei\"}"

var script = "import \"log\"\nimport \"context\"\n  import \"dm\"\nfunc Handle(ctx context.Context,in *dm.PublishMsg) *dm.PublishMsg{\nlog.Print(in)\nreturn in\n  }"

func (l *ProtocolScriptDebugLogic) ProtocolScriptDebug(in *dm.ProtocolScriptDebugReq) (*dm.ProtocolScriptDebugResp, error) {
	switch in.TriggerTimer {
	case protocol.TriggerTimerBefore:
		var req deviceMsg.PublishMsg
		err := json.Unmarshal([]byte(in.Req), &req)
		if err != nil {
			return nil, errors.Parameter.AddMsgf("入参错误,需要的结构体格式为:%s", PublishMsgStr)
		}
		ret, logs, err := l.svcCtx.ScriptTrans.PublishMsgRun(l.ctx, &req, in.Script)
		if err != nil {
			return nil, err
		}
		return &dm.ProtocolScriptDebugResp{
			Out:  utils.MarshalNoErr(ret),
			Logs: logs,
		}, nil
	case protocol.TriggerTimerAfter:
		var req deviceMsg.PublishMsg
		err := json.Unmarshal([]byte(in.Req), &req)
		if err != nil {
			return nil, errors.Parameter.AddMsgf("入参错误,需要的结构体格式为:%s", PublishMsgStr)
		}
		var resp deviceMsg.PublishMsg
		err = json.Unmarshal([]byte(in.Resp), &resp)
		if err != nil {
			return nil, errors.Parameter.AddMsgf("入参错误,需要的结构体格式为:%s", PublishMsgStr)
		}
		logs, err := l.svcCtx.ScriptTrans.RespMsgRun(l.ctx, &req, &resp, in.Script)
		if err != nil {
			return nil, err
		}
		return &dm.ProtocolScriptDebugResp{
			Logs: logs,
		}, nil
	}

	return &dm.ProtocolScriptDebugResp{}, nil
}
