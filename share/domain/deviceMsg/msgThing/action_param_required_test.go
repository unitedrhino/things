// 文件说明：验证行为调用和回复参数的必填兼容规则。
package msgThing

import (
	"testing"

	"gitee.com/unitedrhino/things/share/domain/schema"
)

func actionParamBoolPtr(value bool) *bool {
	return &value
}

func actionParamModel(required *bool) *schema.Model {
	param := &schema.Param{
		Identifier: "level",
		Name:       "等级",
		Required:   required,
		Define: schema.Define{
			Type:  schema.DataTypeInt,
			Min:   "0",
			Max:   "10",
			Start: "0",
			Step:  "1",
		},
	}
	return &schema.Model{
		Action: schema.ActionMap{
			"setLevel": {
				In:  map[string]*schema.Param{"level": param},
				Out: map[string]*schema.Param{"level": param},
			},
		},
	}
}

func TestVerifyReqParamActionInputRequiredCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		required *bool
		params   map[string]any
		wantErr  bool
	}{
		{name: "历史参数缺失时仍必填", required: nil, params: map[string]any{}, wantErr: true},
		{name: "明确必填参数缺失时报错", required: actionParamBoolPtr(true), params: map[string]any{}, wantErr: true},
		{name: "明确可选参数缺失时跳过", required: actionParamBoolPtr(false), params: map[string]any{}, wantErr: false},
		{name: "明确可选参数提供合法值时校验", required: actionParamBoolPtr(false), params: map[string]any{"level": 3}, wantErr: false},
		{name: "明确可选参数提供非法值时仍报错", required: actionParamBoolPtr(false), params: map[string]any{"level": "bad"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &Req{
				ActionID: "setLevel",
				Params:   test.params,
			}
			_, err := req.VerifyReqParam(actionParamModel(test.required), schema.ParamActionInput)
			if (err != nil) != test.wantErr {
				t.Fatalf("VerifyReqParam() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestVerifyRespParamActionOutputRequiredCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		required *bool
		params   map[string]any
		wantErr  bool
	}{
		{name: "历史参数缺失时仍必填", required: nil, params: map[string]any{}, wantErr: true},
		{name: "明确必填参数缺失时报错", required: actionParamBoolPtr(true), params: map[string]any{}, wantErr: true},
		{name: "明确可选参数缺失时跳过", required: actionParamBoolPtr(false), params: map[string]any{}, wantErr: false},
		{name: "明确可选参数提供合法值时校验", required: actionParamBoolPtr(false), params: map[string]any{"level": 3}, wantErr: false},
		{name: "明确可选参数提供非法值时仍报错", required: actionParamBoolPtr(false), params: map[string]any{"level": "bad"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := &Resp{}
			resp.Data = test.params
			_, err := resp.VerifyRespParam(actionParamModel(test.required), "setLevel", schema.ParamActionOutput)
			if (err != nil) != test.wantErr {
				t.Fatalf("VerifyRespParam() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
