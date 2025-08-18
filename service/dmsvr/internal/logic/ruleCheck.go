package logic

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/domain/schema"
)

type CheckOption = func(do any) error

func CheckAffordance(identifier string, po *relationDB.DmSchemaCore, cs *relationDB.DmCommonSchema, checkOptions ...CheckOption) error {
	var affordance interface {
		ValidateWithFmt() error
	}
	switch schema.AffordanceType(po.Type) {
	case schema.AffordanceTypeEvent:
		var (
			do    *schema.Event
			tmpDo *schema.Event
		)
		do = relationDB.ToEventDo(identifier, po)
		if cs != nil {
			tmpDo = relationDB.ToEventDo(identifier, &cs.DmSchemaCore)
		}
		affordance = schema.EventFromCommonSchema(do, tmpDo)
	case schema.AffordanceTypeProperty:
		var (
			do    *schema.Property
			tmpDo *schema.Property
		)
		do = relationDB.ToPropertyDo(identifier, po)
		if cs != nil {
			tmpDo = relationDB.ToPropertyDo(identifier, &cs.DmSchemaCore)
		}
		affordance = schema.PropertyFromCommonSchema(do, tmpDo)
	case schema.AffordanceTypeAction:
		var (
			do    *schema.Action
			tmpDo *schema.Action
		)
		do = relationDB.ToActionDo(identifier, po)
		if cs != nil {
			tmpDo = relationDB.ToActionDo(identifier, &cs.DmSchemaCore)
		}
		affordance = schema.ActionFromCommonSchema(do, tmpDo)
	}
	if err := affordance.ValidateWithFmt(); err != nil {
		return err
	}
	for _, checkOption := range checkOptions {
		if err := checkOption(affordance); err != nil {
			return err
		}
	}
	po.Affordance = relationDB.ToAffordancePo(affordance)
	return nil
}
