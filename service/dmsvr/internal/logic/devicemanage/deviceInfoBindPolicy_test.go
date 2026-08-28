// deviceInfoBindPolicy_test.go 验证设备跨项目直接重绑策略的作用范围。
package devicemanagelogic

import (
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/product"
)

func TestCanDirectlyRebindOwnedDevice(t *testing.T) {
	tests := []struct {
		name       string
		ownership  deviceBindOwnership
		netType    int64
		bindLevel  int64
		wantDirect bool
	}{
		{
			name:       "本人其它项目的老广播蓝牙弱绑定设备",
			ownership:  deviceBindOwnershipCurrentUserOwned,
			netType:    legacyBroadcastBleNetType,
			bindLevel:  product.BindLeveWeak3,
			wantDirect: true,
		},
		{
			name:       "新蓝牙网络类型保持原行为",
			ownership:  deviceBindOwnershipCurrentUserOwned,
			netType:    5,
			bindLevel:  product.BindLeveWeak3,
			wantDirect: false,
		},
		{
			name:       "普通项目成员不能迁移设备",
			ownership:  deviceBindOwnershipCurrentUserBound,
			netType:    legacyBroadcastBleNetType,
			bindLevel:  product.BindLeveWeak3,
			wantDirect: false,
		},
		{
			name:       "非弱绑定设备保持原行为",
			ownership:  deviceBindOwnershipCurrentUserOwned,
			netType:    legacyBroadcastBleNetType,
			bindLevel:  2,
			wantDirect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := canDirectlyRebindOwnedDevice(test.ownership, test.netType, test.bindLevel)
			if got != test.wantDirect {
				t.Fatalf("canDirectlyRebindOwnedDevice() = %v, want %v", got, test.wantDirect)
			}
		})
	}
}
