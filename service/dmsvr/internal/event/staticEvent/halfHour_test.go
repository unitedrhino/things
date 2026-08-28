package staticEvent

import (
	"errors"
	"reflect"
	"testing"
)

func TestRunDeviceStatusStagesOrder(t *testing.T) {
	var calls []string
	stage := func(name string) func() error {
		return func() error {
			calls = append(calls, name)
			return nil
		}
	}
	if err := runDeviceStatusStages(stage("expire"), stage("recover"), stage("mark")); err != nil {
		t.Fatal(err)
	}
	if want := []string{"expire", "recover", "mark"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestRunDeviceStatusStagesStopsAfterExpireFailure(t *testing.T) {
	var calls []string
	errExpected := errors.New("expire failed")
	err := runDeviceStatusStages(func() error {
		calls = append(calls, "expire")
		return errExpected
	}, func() error {
		calls = append(calls, "recover")
		return nil
	}, func() error {
		calls = append(calls, "mark")
		return nil
	})
	if !errors.Is(err, errExpected) {
		t.Fatalf("error = %v, want %v", err, errExpected)
	}
	if want := []string{"expire"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}
