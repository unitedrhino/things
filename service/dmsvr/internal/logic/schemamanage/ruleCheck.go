package schemamanagelogic

import (
	"gitee.com/i-Things/share/domain/schema"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
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
	po.Affordance = relationDB.ToAffordancePo(affordance)
	return nil
}
