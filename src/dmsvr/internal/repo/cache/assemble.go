package cache

import (
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"time"
)

func ToPropertyDB(productID string, in *schema.Property) *mysql.ProductSchema2 {
	define := mysql.PropertyDef{
		Define: in.Define,
		Mode:   in.Mode,
	}
	defineStr, _ := json.Marshal(define)
	return &mysql.ProductSchema2{
		ProductID:   productID,
		Tag:         int64(schema.TagCustom),
		Type:        int64(schema.AffordanceTypeProperty),
		Identifier:  in.Identifier,
		Name:        in.Name,
		Desc:        in.Desc,
		Required:    def.ToIntBool[int64](in.Required),
		Affordance:  string(defineStr),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}

func ToPropertyDo(in *mysql.ProductSchema2) *schema.Property {
	affordance := mysql.PropertyDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	return &schema.Property{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Define:     affordance.Define,
		Mode:       affordance.Mode,
	}
}

func ToEventDB(productID string, in *schema.Event) *mysql.ProductSchema2 {
	define := mysql.EventDef{
		Type:   in.Type,
		Params: in.Params,
	}
	defineStr, _ := json.Marshal(define)
	return &mysql.ProductSchema2{
		ProductID:   productID,
		Tag:         int64(schema.TagCustom),
		Type:        int64(schema.AffordanceTypeEvent),
		Identifier:  in.Identifier,
		Name:        in.Name,
		Desc:        in.Desc,
		Required:    def.ToIntBool[int64](in.Required),
		Affordance:  string(defineStr),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}

func ToEventDo(in *mysql.ProductSchema2) *schema.Event {
	affordance := mysql.EventDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	return &schema.Event{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Type:       affordance.Type,
		Params:     affordance.Params,
	}
}

func ToActionDB(productID string, in *schema.Action) *mysql.ProductSchema2 {
	define := mysql.ActionDef{
		Input:  in.Input,
		Output: in.Output,
	}
	defineStr, _ := json.Marshal(define)
	return &mysql.ProductSchema2{
		ProductID:   productID,
		Tag:         int64(schema.TagCustom),
		Type:        int64(schema.AffordanceTypeAction),
		Identifier:  in.Identifier,
		Name:        in.Name,
		Desc:        in.Desc,
		Required:    def.ToIntBool[int64](in.Required),
		Affordance:  string(defineStr),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}

func ToActionDo(in *mysql.ProductSchema2) *schema.Action {
	affordance := mysql.ActionDef{}
	_ = json.Unmarshal([]byte(in.Affordance), &affordance)
	return &schema.Action{
		Identifier: in.Identifier,
		Name:       in.Name,
		Desc:       in.Desc,
		Required:   def.ToBool(in.Required),
		Input:      affordance.Input,
		Output:     affordance.Output,
	}
}

func ToSchemaDo(in []*mysql.ProductSchema2) *schema.Model {
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
