package schemamanagelogic

import (
	"gitee.com/i-Things/core/shared/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
)

func CheckAffordance(po *relationDB.DmSchemaCore) error {
	var affordance interface {
		ValidateWithFmt() error
	}
	switch schema.AffordanceType(po.Type) {
	case schema.AffordanceTypeEvent:
		affordance = relationDB.ToEventDo(po)
	case schema.AffordanceTypeProperty:
		affordance = relationDB.ToPropertyDo(po)
	case schema.AffordanceTypeAction:
		affordance = relationDB.ToActionDo(po)
	}
	if err := affordance.ValidateWithFmt(); err != nil {
		return err
	}
	//todo 规范化了之后前端会有问题
	po.Affordance = relationDB.ToAffordancePo(affordance)
	return nil
}
