// 本文件验证创建单设备分享后会失效对应的权限缓存。
package userdevicelogic

import (
	"context"
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func TestInvalidateCreatedUserDeviceShareClearsExactCacheKey(t *testing.T) {
	share := &relationDB.DmUserDeviceShare{
		ProductID:    "008",
		DeviceName:   "C8586A69A7D4",
		SharedUserID: 334185784940448,
	}

	var (
		calls   int
		gotKey  userShared.UserShareKey
		gotData *dm.UserDeviceShareInfo
	)
	setData := func(_ context.Context, key userShared.UserShareKey, data *dm.UserDeviceShareInfo) error {
		calls++
		gotKey = key
		gotData = data
		return nil
	}

	invalidateCreatedUserDeviceShare(context.Background(), setData, share)

	if calls != 1 {
		t.Fatalf("SetData calls = %d, want 1", calls)
	}
	wantKey := (userShared.UserShareKey{
		ProductID:    share.ProductID,
		DeviceName:   share.DeviceName,
		SharedUserID: share.SharedUserID,
	})
	if gotKey != wantKey {
		t.Fatalf("cache key = %#v, want %#v", gotKey, wantKey)
	}
	if gotData != nil {
		t.Fatalf("cache data = %#v, want nil", gotData)
	}
}
