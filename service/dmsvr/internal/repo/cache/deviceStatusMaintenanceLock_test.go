package cache

import (
	"context"
	"testing"
)

type fakeMaintenanceLockStore struct {
	value   string
	seconds int
}

func (s *fakeMaintenanceLockStore) SetnxExCtx(_ context.Context, _ string, value string, seconds int) (bool, error) {
	if s.value != "" {
		return false, nil
	}
	s.value = value
	s.seconds = seconds
	return true, nil
}

func (s *fakeMaintenanceLockStore) EvalCtx(_ context.Context, _ string, _ string, args ...any) (any, error) {
	token, _ := args[0].(string)
	if token != s.value {
		return int64(0), nil
	}
	s.value = ""
	return int64(1), nil
}

func TestDeviceStatusMaintenanceLockOwnership(t *testing.T) {
	store := &fakeMaintenanceLockStore{}
	lock := NewDeviceStatusMaintenanceLock(store)
	token, acquired, err := lock.TryLock(context.Background())
	if err != nil || !acquired || token == "" {
		t.Fatalf("TryLock() token=%q acquired=%v err=%v", token, acquired, err)
	}
	if store.seconds != 600 {
		t.Fatalf("lock ttl = %d, want 600", store.seconds)
	}
	if _, acquired, err = lock.TryLock(context.Background()); err != nil || acquired {
		t.Fatalf("second TryLock() acquired=%v err=%v, want locked", acquired, err)
	}
	if err = lock.Unlock(context.Background(), "another-owner"); err != nil {
		t.Fatal(err)
	}
	if store.value != token {
		t.Fatal("wrong owner released the lock")
	}
	if err = lock.Unlock(context.Background(), token); err != nil {
		t.Fatal(err)
	}
	if store.value != "" {
		t.Fatal("owner failed to release the lock")
	}
}
