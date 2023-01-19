package mysql

import (
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
)

func ToPropertyPo(productID string, in *schema.Property) *DmProductSchema {
	define := PropertyDef{
		Define: in.Define,
		Mode:   in.Mode,
	}
	defineStr, _ := json.Marshal(define)
	return &DmProductSchema{
		ProductID:  productID,
		Tag:        int64(schema.TagCustom),
		Type:       int64(schema.AffordanceTypeProperty),
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToIntBool[int64](in.Required),
		Affordance: string(defineStr),
	}
}

func ToPropertyDo(in *DmProductSchema) *schema.Property {
	affordance := PropertyDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Property{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Define:     affordance.Define,
		Mode:       affordance.Mode,
	}
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
		ProductID:  productID,
		Tag:        int64(schema.TagCustom),
		Type:       int64(schema.AffordanceTypeEvent),
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToIntBool[int64](in.Required),
		Affordance: string(defineStr),
	}
}

func ToEventDo(in *DmProductSchema) *schema.Event {
	affordance := EventDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Event{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Type:       affordance.Type,
		Params:     affordance.Params,
	}
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
		ProductID:  productID,
		Tag:        int64(schema.TagCustom),
		Type:       int64(schema.AffordanceTypeAction),
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToIntBool[int64](in.Required),
		Affordance: string(defineStr),
	}
}

func ToActionDo(in *DmProductSchema) *schema.Action {
	affordance := ActionDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	do := &schema.Action{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Input:      affordance.Input,
		Output:     affordance.Output,
	}
	do.ValidateWithFmt()
	return do
}

func ToSchemaDo(in []*DmProductSchema) *schema.Model {
	if len(in) == 0 {
		return nil
	}
	model := schema.Model{
		Profile: schema.Profile{ProductID: in[0].ProductID},
	}
	for _, v := range in {
		switch schema.AffordanceType(v.Type) {
		case schema.AffordanceTypeProperty:
			model.Properties = append(model.Properties, *ToPropertyDo(v))
		case schema.AffordanceTypeEvent:
			model.Events = append(model.Events, *ToEventDo(v))
		case schema.AffordanceTypeAction:
			model.Actions = append(model.Actions, *ToActionDo(v))
		}
	}
	model.ValidateWithFmt()
	return &model
}
