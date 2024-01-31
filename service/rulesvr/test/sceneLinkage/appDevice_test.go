package sceneLinkage

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/schema"
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/service/rulesvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/ruledirect"
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
	svcCtx = ruledirect.GetSvcCtx()
	Logger = logx.WithContext(ctx)
	Init()
	t.Run()
	fmt.Println("end test appDeviceEvent")
}

func Init() {
	_, err := svcCtx.SceneRepo.Insert(ctx,
		&scene.Info{
			Name:        "上线测试2",
			Desc:        "上线测试2",
			TriggerType: scene.TriggerTypeDevice,
			Trigger: scene.Trigger{
				Device: scene.TriggerDevices{&scene.TriggerDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Operator:  scene.DeviceOperationOperatorConnected,
					OperationSchema: &scene.OperationSchema{
						DataID:   []string{"GPS_Info", "longtitude"},
						TermType: scene.CmpTypeGt,
						Values:   []string{"0.001"},
					},
				}},
			},
			When: scene.Terms{&scene.Term{
				ColumnType: scene.TermColumnTypeProperty,
				ColumnSchema: &scene.ColumnSchema{
					ProductID:  productID,
					DeviceName: deviceName,
					DataID:     []string{"GPS_Info", "longtitude"},
					TermType:   scene.CmpTypeGt,
					Values:     []string{"0.001"},
				},
				NextCondition: scene.TermConditionTypeOr,
				Terms:         nil,
			}},
			Then: scene.Actions{&scene.Action{
				Executor: scene.ActionExecutorAlarm,
				Alarm:    &scene.ActionAlarm{Mode: scene.ActionAlarmModeTrigger},
			},
				&scene.Action{
					Executor: scene.ActionExecutorDelay,
					Delay: &scene.UnitTime{
						Time: 2,
						Unit: scene.TimeUnitSeconds,
					},
				},
				&scene.Action{
					Executor: scene.ActionExecutorDevice,
					Delay:    nil,
					Alarm:    nil,
					Device: &scene.ActionDevice{
						ProductID: productID,
						Selector:  scene.DeviceSelectorAll,
						Type:      scene.ActionDeviceTypePropertyControl,
						Value: `{"GPS_Info":
			{
				"longtitude":33,
				"latitude":33
			}}`,
					},
				},
				&scene.Action{
					Executor: scene.ActionExecutorDevice,
					Delay:    nil,
					Alarm:    nil,
					Device: &scene.ActionDevice{
						ProductID: productID,
						Selector:  scene.DeviceSelectorAll,
						Type:      scene.ActionDeviceTypePropertyControl,
						Value:     `{"battery_state":14}`,
					},
				},
				&scene.Action{
					Executor: scene.ActionExecutorDevice,
					Delay:    nil,
					Alarm:    nil,
					Device: &scene.ActionDevice{
						ProductID: productID,
						Selector:  scene.DeviceSelectorAll,
						Type:      scene.ActionDeviceTypeAction,
						DataID:    "whistle",
						Value:     `{"time":123,"switch":1}`,
					},
				},
			},
			State: 1,
		})
	fmt.Println(err)
	_, err = svcCtx.SceneRepo.Insert(ctx, &scene.Info{
		Name:        "结构体上报1",
		Desc:        "结构体上报1",
		TriggerType: scene.TriggerTypeDevice,
		Trigger: scene.Trigger{
			Device: scene.TriggerDevices{&scene.TriggerDevice{
				ProductID: productID,
				Selector:  scene.DeviceSelectorAll,
				Operator:  scene.DeviceOperationOperatorReportProperty,
				OperationSchema: &scene.OperationSchema{
					DataID:   []string{"GPS_Info", "longtitude"},
					TermType: scene.CmpTypeGt,
					Values:   []string{"0.001"},
				},
			}},
		},
		When: scene.Terms{&scene.Term{
			ColumnType: scene.TermColumnTypeProperty,
			ColumnSchema: &scene.ColumnSchema{
				ProductID:  productID,
				DeviceName: deviceName,
				DataID:     []string{"GPS_Info", "longtitude"},
				TermType:   scene.CmpTypeGt,
				Values:     []string{"0.001"},
			},
			NextCondition: scene.TermConditionTypeOr,
			Terms:         nil,
		}},
		Then: scene.Actions{&scene.Action{
			Executor: scene.ActionExecutorAlarm,
			Alarm:    &scene.ActionAlarm{Mode: scene.ActionAlarmModeTrigger},
		}, &scene.Action{
			Executor: scene.ActionExecutorDelay,
			Delay: &scene.UnitTime{
				Time: 2,
				Unit: scene.TimeUnitSeconds,
			},
		},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypePropertyControl,
					Value: `{"GPS_Info":
			{
				"longtitude":33,
				"latitude":33
			}}`,
				},
			},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypePropertyControl,
					Value:     `{"battery_state":14}`,
				},
			},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypeAction,
					DataID:    "whistle",
					Value:     `{"time":123,"switch":1}`,
				},
			},
		},
		State: 1,
	})
	fmt.Println(err)
	_, err = svcCtx.SceneRepo.Insert(ctx, &scene.Info{
		Name:        "定时两秒",
		Desc:        "定时2秒",
		TriggerType: scene.TriggerTypeTimer,
		Trigger: scene.Trigger{
			Timer: &scene.Timer{
				Type: "cron",
				Cron: "*/2 * * * * ? ",
			},
		},
		When: scene.Terms{&scene.Term{
			ColumnType: scene.TermColumnTypeProperty,
			ColumnSchema: &scene.ColumnSchema{
				ProductID:  productID,
				DeviceName: deviceName,
				DataID:     []string{"GPS_Info", "longtitude"},
				TermType:   scene.CmpTypeGt,
				Values:     []string{"0.001"},
			},
			NextCondition: scene.TermConditionTypeOr,
			Terms:         nil,
		}},
		Then: scene.Actions{&scene.Action{
			Executor: scene.ActionExecutorAlarm,
			Alarm:    &scene.ActionAlarm{Mode: scene.ActionAlarmModeTrigger},
		}, &scene.Action{
			Executor: scene.ActionExecutorDelay,
			Delay: &scene.UnitTime{
				Time: 2,
				Unit: scene.TimeUnitSeconds,
			},
		},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypePropertyControl,
					Value: `{"GPS_Info":
			{
				"longtitude":33,
				"latitude":33
			}}`,
				},
			},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypePropertyControl,
					Value:     `{"battery_state":14}`,
				},
			},
			&scene.Action{
				Executor: scene.ActionExecutorDevice,
				Delay:    nil,
				Alarm:    nil,
				Device: &scene.ActionDevice{
					ProductID: productID,
					Selector:  scene.DeviceSelectorAll,
					Type:      scene.ActionDeviceTypeAction,
					DataID:    "whistle",
					Value:     `{"time":123,"switch":1}`,
				},
			},
		},
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
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{in: &application.PropertyReport{
			Device:     devices.Core{ProductID: productID, DeviceName: deviceName},
			Timestamp:  time.Now().UnixMilli(),
			Identifier: "GPS_Info",
			Param: application.ParamValue{Type: schema.DataTypeStruct, Value: application.StructValue{
				"longtitude": {Type: schema.DataTypeFloat, Value: 0.12},
				"latitude":   {Type: schema.DataTypeFloat, Value: 0.42},
			}},
		}}},
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
