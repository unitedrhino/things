package cache

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	deviceStatusMaintenanceLockKey = "dm:device:status:maintenance:lock"
	deviceStatusMaintenanceLockTTL = 10 * time.Minute
)

const releaseDeviceStatusMaintenanceLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end
return 0
`

type deviceStatusMaintenanceLockStore interface {
	SetnxExCtx(ctx context.Context, key, value string, seconds int) (bool, error)
	EvalCtx(ctx context.Context, script, key string, args ...any) (any, error)
}

// DeviceStatusMaintenanceLock prevents concurrent half-hour device status jobs.
type DeviceStatusMaintenanceLock struct {
	store deviceStatusMaintenanceLockStore
}

func NewDeviceStatusMaintenanceLock(store deviceStatusMaintenanceLockStore) *DeviceStatusMaintenanceLock {
	return &DeviceStatusMaintenanceLock{store: store}
}

// TryLock returns an ownership token when the caller acquires the lock.
func (l *DeviceStatusMaintenanceLock) TryLock(ctx context.Context) (string, bool, error) {
	token := uuid.NewString()
	ok, err := l.store.SetnxExCtx(ctx, deviceStatusMaintenanceLockKey, token, int(deviceStatusMaintenanceLockTTL/time.Second))
	if err != nil || !ok {
		return "", ok, err
	}
	return token, true, nil
}

// Unlock removes the lock only when token still belongs to this caller.
func (l *DeviceStatusMaintenanceLock) Unlock(ctx context.Context, token string) error {
	_, err := l.store.EvalCtx(ctx, releaseDeviceStatusMaintenanceLockScript, deviceStatusMaintenanceLockKey, token)
	return err
}
