package productmanagelogic

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
)

func CheckAffordance(po *mysql.DmProductSchema) error {
	var affordance interface {
		ValidateWithFmt() error
	}
	switch schema.AffordanceType(po.Type) {
	case schema.AffordanceTypeEvent:
		affordance = mysql.ToEventDo(po)
	case schema.AffordanceTypeProperty:
		affordance = mysql.ToPropertyDo(po)
	case schema.AffordanceTypeAction:
		affordance = mysql.ToActionDo(po)
	}
	if err := affordance.ValidateWithFmt(); err != nil {
		return err
	}
	return nil
}
