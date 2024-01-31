package relationDB

import (
	"encoding/json"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/domain/shadow"
)

func ToPropertyPo(productID string, in *schema.Property) *DmProductSchema {
	define := PropertyDef{
		Define:      in.Define,
		Mode:        in.Mode,
		IsUseShadow: in.IsUseShadow,
		IsNoRecord:  in.IsNoRecord,
	}
	defineStr, _ := json.Marshal(define)
	return &DmProductSchema{
		ProductID: productID,
		DmSchemaCore: DmSchemaCore{
			Type:         int64(schema.AffordanceTypeProperty),
			Identifier:   in.Identifier,
			Name:         in.Name,
			ExtendConfig: in.ExtendConfig,
			Desc:         in.Desc,
			Required:     def.ToIntBool[int64](in.Required),
			Affordance:   string(defineStr),
		},
		Tag: int64(schema.TagCustom),
	}
}

func ToCommonParam(in *DmSchemaCore) schema.CommonParam {
	return schema.CommonParam{
		Identifier:   in.Identifier,
		Name:         in.Name,
		Desc:         in.Desc,
		ExtendConfig: in.ExtendConfig,
		Required:     def.ToBool(in.Required),
	}
}

func ToPropertyDo(in *DmSchemaCore) *schema.Property {
	affordance := PropertyDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Property{
		CommonParam: ToCommonParam(in),
		Define:      affordance.Define,
		Mode:        affordance.Mode,
		IsUseShadow: affordance.IsUseShadow,
		IsNoRecord:  affordance.IsNoRecord,
	}
	newAffordance, _ := json.Marshal(affordance)
	in.Affordance = string(newAffordance)
	do.ValidateWithFmt()
	return do
}

func ToEventPo(productID string, in *schema.Event) *DmProductSchema {
	define := EventDef{
		Type:   in.Type,
		Params: in.Params,
	}
	defineStr, _ := json.Marshal(define)
	return &DmProductSchema{
		ProductID: productID,
		DmSchemaCore: DmSchemaCore{
			Type:         int64(schema.AffordanceTypeEvent),
			Identifier:   in.Identifier,
			Name:         in.Name,
			Desc:         in.Desc,
			ExtendConfig: in.ExtendConfig,
			Required:     def.ToIntBool[int64](in.Required),
			Affordance:   string(defineStr),
		},
		Tag: int64(schema.TagCustom),
	}
}

func ToEventDo(in *DmSchemaCore) *schema.Event {
	affordance := EventDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Event{
		CommonParam: ToCommonParam(in),
		Type:        affordance.Type,
		Params:      affordance.Params,
	}
	newAffordance, _ := json.Marshal(affordance)
	in.Affordance = string(newAffordance)
	do.ValidateWithFmt()
	return do
}

func ToActionPo(productID string, in *schema.Action) *DmProductSchema {
	define := ActionDef{
		Input:  in.Input,
		Output: in.Output,
	}
	defineStr, _ := json.Marshal(define)
	return &DmProductSchema{
		ProductID: productID,
		Tag:       int64(schema.TagCustom),
		DmSchemaCore: DmSchemaCore{
			Identifier:   in.Identifier,
			Type:         int64(schema.AffordanceTypeAction),
			Name:         in.Name,
			ExtendConfig: in.ExtendConfig,
			Desc:         in.Desc,
			Required:     def.ToIntBool[int64](in.Required),
			Affordance:   string(defineStr),
		},
	}
}

func ToAffordancePo(in any) string {
	var defineStr []byte
	switch in.(type) {
	case *schema.Event:
		af := in.(*schema.Event)
		define := EventDef{
			Type:   af.Type,
			Params: af.Params,
		}
		defineStr, _ = json.Marshal(define)
	case *schema.Action:
		af := in.(*schema.Action)
		define := ActionDef{
			Input:  af.Input,
			Output: af.Output,
		}
		defineStr, _ = json.Marshal(define)
	case *schema.Property:
		af := in.(*schema.Property)
		define := PropertyDef{
			Define:      af.Define,
			Mode:        af.Mode,
			IsUseShadow: af.IsUseShadow,
			IsNoRecord:  af.IsNoRecord,
		}
		defineStr, _ = json.Marshal(define)
	}
	return string(defineStr)
}

func ToActionDo(in *DmSchemaCore) *schema.Action {
	affordance := ActionDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Action{
		CommonParam: ToCommonParam(in),
		Input:       affordance.Input,
		Output:      affordance.Output,
		Dir:         affordance.Dir,
	}
	newAffordance, _ := json.Marshal(affordance)
	in.Affordance = string(newAffordance)
	do.ValidateWithFmt()
	return do
}

func ToSchemaDo(productID string, in []*DmProductSchema) *schema.Model {
	model := schema.Model{
		Profile: schema.Profile{ProductID: productID},
	}
	if len(in) == 0 {
		return &model
	}
	for _, v := range in {
		switch schema.AffordanceType(v.Type) {
		case schema.AffordanceTypeProperty:
			model.Properties = append(model.Properties, *ToPropertyDo(&v.DmSchemaCore))
		case schema.AffordanceTypeEvent:
			model.Events = append(model.Events, *ToEventDo(&v.DmSchemaCore))
		case schema.AffordanceTypeAction:
			model.Actions = append(model.Actions, *ToActionDo(&v.DmSchemaCore))
		}
	}
	model.ValidateWithFmt()
	return &model
}

func ToShadowPo(info *shadow.Info) *DmDeviceShadow {
	return &DmDeviceShadow{
		ID:                info.ID,
		ProductID:         info.ProductID,
		DeviceName:        info.DeviceName,
		DataID:            info.DataID,
		UpdatedDeviceTime: info.UpdatedDeviceTime,
		Value:             info.Value,
	}
}
func ToShadowDo(in *DmDeviceShadow) *shadow.Info {
	return &shadow.Info{
		ID:                in.ID,
		ProductID:         in.ProductID,
		DeviceName:        in.DeviceName,
		DataID:            in.DataID,
		Value:             in.Value,
		UpdatedDeviceTime: in.UpdatedDeviceTime,
		CreatedTime:       in.CreatedTime,
		UpdatedTime:       in.UpdatedTime,
	}
}
func ToShadowsDo(in []*DmDeviceShadow) []*shadow.Info {
	if in == nil {
		return nil
	}
	var ret []*shadow.Info
	for _, v := range in {
		ret = append(ret, ToShadowDo(v))
	}
	return ret
}
