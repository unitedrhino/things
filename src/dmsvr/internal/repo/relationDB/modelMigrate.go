package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/stores"
)

func Migrate() error {
	db := stores.GetCommonConn(context.Background())
	return db.AutoMigrate(
		DmProductInfo{},
		DmProductCustom{},
		DmProductSchema{},
		DmDeviceInfo{},
		DmGroupInfo{},
		DmGroupDevice{},
		DmGatewayDevice{},
		DmProductRemoteConfig{},
	)
}
