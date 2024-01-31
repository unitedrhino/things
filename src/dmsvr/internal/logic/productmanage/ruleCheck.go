package productmanagelogic

import (
	"gitee.com/i-Things/core/shared/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
)

func CheckAffordance(po *relationDB.DmSchemaCore, cs *relationDB.DmCommonSchema) error {
	var affordance interface {
		ValidateWithFmt() error
	}
	switch schema.AffordanceType(po.Type) {
	case schema.AffordanceTypeEvent:
		var (
			do    *schema.Event
			tmpDo *schema.Event
		)
		do = relationDB.ToEventDo(po)
		if cs != nil {
			tmpDo = relationDB.ToEventDo(&cs.DmSchemaCore)
		}
		affordance = schema.EventFromCommonSchema(do, tmpDo)
	case schema.AffordanceTypeProperty:
		var (
			do    *schema.Property
			tmpDo *schema.Property
		)
		do = relationDB.ToPropertyDo(po)
		if cs != nil {
			tmpDo = relationDB.ToPropertyDo(&cs.DmSchemaCore)
		}
		affordance = schema.PropertyFromCommonSchema(do, tmpDo)
	case schema.AffordanceTypeAction:
		var (
			do    *schema.Action
			tmpDo *schema.Action
		)
		do = relationDB.ToActionDo(po)
		if cs != nil {
			tmpDo = relationDB.ToActionDo(&cs.DmSchemaCore)
		}
		affordance = schema.ActionFromCommonSchema(do, tmpDo)
	}
	if err := affordance.ValidateWithFmt(); err != nil {
		return err
	}
	//todo 规范化了之后前端会有问题
	po.Affordance = relationDB.ToAffordancePo(affordance)
	return nil
}
