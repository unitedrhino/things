package sceneLinkage

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/ruledirect"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
	"time"
)

const (
	productID  = "254pwnKQsvK"
	deviceName = "test5"
)

var (
	svcCtx *svc.ServiceContext
	ctx    = context.TODO()
	Logger logx.Logger
)

func TestMain(t *testing.M) {
	fmt.Println("start test appDeviceEvent")
	ruledirect.ConfigFile = "../../etc/rule.yaml"
	svcCtx = ruledirect.GetSvcCtx()
	Logger = logx.WithContext(ctx)
	Init()
	t.Run()
	fmt.Println("end test appDeviceEvent")
}

func Init() {
	err := svcCtx.SceneRepo.Insert(ctx, &scene.Info{
		Name:        "上线测试1",
		Desc:        "上线测试1",
		TriggerType: scene.TriggerTypeDevice,
		Trigger: scene.Trigger{
			Device: scene.TriggerDevices{&scene.TriggerDevice{
				ProductID: productID,
				Selector:  scene.TriggerDeviceSelectorAll,
				Operator:  scene.DeviceOperationOperatorConnected,
				OperationSchema: &scene.OperationSchema{
					DataID:   []string{"GPS_Info", "longtitude"},
					TermType: scene.TermTypeGt,
					Values:   []string{"0.001"},
				},
			}},
		},
		When: nil,
		Then: scene.Actions{&scene.Action{
			Executor: scene.ActionExecutorAlarm,
			Alarm:    &scene.ActionAlarm{Mode: scene.ActionAlarmModeTrigger},
		}},
		State: 1,
	})
	fmt.Println(err)
}

func TestAppDeviceHandle_DeviceEventReport(t *testing.T) {
	type args struct {
		in *application.EventReport
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
			if err := a.DeviceEventReport(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("DeviceEventReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppDeviceHandle_DevicePropertyReport(t *testing.T) {
	type args struct {
		in *application.PropertyReport
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
			if err := a.DevicePropertyReport(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("DevicePropertyReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppDeviceHandle_DeviceStatusConnected(t *testing.T) {
	type fields struct {
		svcCtx *svc.ServiceContext
		ctx    context.Context
		Logger logx.Logger
	}
	type args struct {
		in *application.ConnectMsg
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{in: &application.ConnectMsg{
			Timestamp: time.Now().UnixMilli(), Device: devices.Core{ProductID: productID, DeviceName: deviceName}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
			if err := a.DeviceStatusConnected(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("DeviceStatusConnected() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppDeviceHandle_DeviceStatusDisConnected(t *testing.T) {
	type fields struct {
		svcCtx *svc.ServiceContext
		ctx    context.Context
		Logger logx.Logger
	}
	type args struct {
		in *application.ConnectMsg
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
			if err := a.DeviceStatusDisConnected(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("DeviceStatusDisConnected() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAppDeviceHandle(t *testing.T) {
	var tests []struct {
		name string
		want *appDeviceEvent.AppDeviceHandle
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppDeviceHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
