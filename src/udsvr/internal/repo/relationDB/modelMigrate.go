package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	return db.AutoMigrate(
	//&RuleSceneInfo{},
	//&RuleAlarmScene{},
	//&RuleAlarmInfo{},
	//&RuleAlarmRecord{},
	//&RuleAlarmLog{},
	//&RuleAlarmDealRecord{},
	)
}
