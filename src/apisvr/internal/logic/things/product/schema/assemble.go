package schema

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToSchemaSpecRpc(in *types.SchemaSpec) *dm.SchemaSpec {
	if in == nil {
		return nil
	}
	return &dm.SchemaSpec{
		Identifier: in.Identifier,
		Name:       in.Name,
		DataType:   ToSchemaDefineRpc(in.DataType),
	}
}

func ToSchemaSpecsRpc(in []*types.SchemaSpec) []*dm.SchemaSpec {
	if in == nil {
		return nil
	}
	ret := []*dm.SchemaSpec{}
	for _, v := range in {
		ret = append(ret, ToSchemaSpecRpc(v))
	}
	return ret
}

func ToSchemaSpecTypes(in *dm.SchemaSpec) *types.SchemaSpec {
	if in == nil {
		return nil
	}
	return &types.SchemaSpec{
		Identifier: in.Identifier,
		Name:       in.Name,
		DataType:   ToSchemaDefineTypes(in.DataType),
	}
}

func ToSchemaSpecsTypes(in []*dm.SchemaSpec) []*types.SchemaSpec {
	if in == nil {
		return nil
	}
	ret := []*types.SchemaSpec{}
	for _, v := range in {
		ret = append(ret, ToSchemaSpecTypes(v))
	}
	return ret
}

func ToSchemaDefineRpc(in *types.SchemaDefine) *dm.SchemaDefine {
	if in == nil {
		return nil
	}
	return &dm.SchemaDefine{
		Type:      in.Type,
		Mapping:   in.Mapping,
		Min:       in.Min,
		Max:       in.Max,
		Start:     in.Start,
		Step:      in.Step,
		Unit:      in.Unit,
		Specs:     ToSchemaSpecsRpc(in.Specs),
		ArrayInfo: ToSchemaDefineRpc(in.ArrayInfo),
	}
}

func ToSchemaDefineTypes(in *dm.SchemaDefine) *types.SchemaDefine {
	if in == nil {
		return nil
	}
	return &types.SchemaDefine{
		Type:      in.Type,
		Mapping:   in.Mapping,
		Min:       in.Min,
		Max:       in.Max,
		Start:     in.Start,
		Step:      in.Step,
		Unit:      in.Unit,
		Specs:     ToSchemaSpecsTypes(in.Specs),
		ArrayInfo: ToSchemaDefineTypes(in.ArrayInfo),
	}
}

func ToSchemaDefinesTypes(in []*dm.SchemaDefine) []*types.SchemaDefine {
	if in == nil {
		return nil
	}
	ret := []*types.SchemaDefine{}
	for _, v := range in {
		ret = append(ret, ToSchemaDefineTypes(v))
	}
	return ret
}

func ToSchemaParamRpc(in *types.SchemaParam) *dm.SchemaParam {
	if in == nil {
		return nil
	}
	return &dm.SchemaParam{
		Identifier: in.Identifier,
		Name:       in.Name,
		Define:     ToSchemaDefineRpc(in.Define),
	}
}
func ToSchemaParamsRpc(in []*types.SchemaParam) []*dm.SchemaParam {
	if in == nil {
		return nil
	}
	ret := []*dm.SchemaParam{}
	for _, v := range in {
		ret = append(ret, ToSchemaParamRpc(v))
	}
	return ret
}

func ToSchemaParamTypes(in *dm.SchemaParam) *types.SchemaParam {
	if in == nil {
		return nil
	}
	return &types.SchemaParam{
		Identifier: in.Identifier,
		Name:       in.Name,
		Define:     ToSchemaDefineTypes(in.Define),
	}
}
func ToSchemaParamsTypes(in []*dm.SchemaParam) []*types.SchemaParam {
	if in == nil {
		return nil
	}
	ret := []*types.SchemaParam{}
	for _, v := range in {
		ret = append(ret, ToSchemaParamTypes(v))
	}
	return ret
}

func ToActionRpc(in *types.SchemaAction) *dm.SchemaAction {
	if in == nil {
		return nil
	}
	return &dm.SchemaAction{
		Input:  ToSchemaParamsRpc(in.Input),
		Output: ToSchemaParamsRpc(in.Output),
	}
}
func ToActionTypes(in *dm.SchemaAction) *types.SchemaAction {
	if in == nil {
		return nil
	}
	return &types.SchemaAction{
		Input:  ToSchemaParamsTypes(in.Input),
		Output: ToSchemaParamsTypes(in.Output),
	}
}

func ToEventRpc(in *types.SchemaEvent) *dm.SchemaEvent {
	if in == nil {
		return nil
	}
	return &dm.SchemaEvent{
		Type:   in.Type,
		Params: ToSchemaParamsRpc(in.Params),
	}
}
func ToEventTypes(in *dm.SchemaEvent) *types.SchemaEvent {
	if in == nil {
		return nil
	}
	return &types.SchemaEvent{
		Type:   in.Type,
		Params: ToSchemaParamsTypes(in.Params),
	}
}

func ToPropertyRpc(in *types.SchemaProperty) *dm.SchemaProperty {
	if in == nil {
		return nil
	}
	return &dm.SchemaProperty{
		Mode:   in.Mode,
		Define: ToSchemaDefineRpc(in.Define),
	}
}
func ToPropertyTypes(in *dm.SchemaProperty) *types.SchemaProperty {
	if in == nil {
		return nil
	}
	return &types.SchemaProperty{
		Mode:   in.Mode,
		Define: ToSchemaDefineTypes(in.Define),
	}
}

func ToSchemaInfoRpc(in *types.ProductSchemaInfo) *dm.ProductSchemaInfo {
	if in == nil {
		return nil
	}
	rpc := &dm.ProductSchemaInfo{
		ProductID:  in.ProductID,
		Type:       in.Type,
		Tag:        in.Tag,
		Identifier: in.Identifier,
		Name:       utils.ToRpcNullString(in.Name),
		Desc:       utils.ToRpcNullString(in.Desc),
		Required:   in.Required,
		Property:   ToPropertyRpc(in.Property),
		Event:      ToEventRpc(in.Event),
		Action:     ToActionRpc(in.Action),
	}
	return rpc
}

func ToSchemaInfoTypes(in *dm.ProductSchemaInfo) *types.ProductSchemaInfo {
	if in == nil {
		return nil
	}
	rpc := types.ProductSchemaInfo{
		ProductID:  in.ProductID,
		Type:       in.Type,
		Tag:        in.Tag,
		Identifier: in.Identifier,
		Name:       utils.ToNullString(in.Name),
		Desc:       utils.ToNullString(in.Desc),
		Required:   in.Required,
		Property:   ToPropertyTypes(in.Property),
		Event:      ToEventTypes(in.Event),
		Action:     ToActionTypes(in.Action),
	}
	return &rpc
}
