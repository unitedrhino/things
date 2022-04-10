package tdengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"time"
)

type DeviceDataRepo struct {
	t *TDengine.Td
}

func NewDeviceDataRepo(dataSource string) *DeviceDataRepo {
	td, err := TDengine.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td}
}

func (d *DeviceDataRepo) InsertEventData(ctx context.Context, t *deviceTemplate.Template, productID string,
	deviceName string, event *deviceTemplate.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	sql := fmt.Sprintf("insert into %s (ts,event_id,event_type, param) values (?,?,?,?);", getEventTableName(productID, deviceName))
	if _, err := d.t.Exec(sql, event.TimeStamp, event.ID, event.Type, string(param)); err != nil {
		return err
	}
	return nil
}

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string, property *deviceTemplate.PropertyData) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) GetEventDataWithID(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.EventData, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) GetPropertyDataWithID(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.PropertyData, error) {
	//TODO implement me
	panic("implement me")
}

func getTdType(define deviceTemplate.Define) string {
	switch define.Type {
	case deviceTemplate.BOOL:
		return "BOOL"
	case deviceTemplate.INT:
		return "BIGINT"
	case deviceTemplate.STRING:
		return fmt.Sprintf("BINARY(%s)", define.Max)
	case deviceTemplate.STRUCT:
		return "BINARY(5000)"
	case deviceTemplate.FLOAT:
		return "DOUBLE"
	case deviceTemplate.TIMESTAMP:
		return "TIMESTAMP"
	case deviceTemplate.ARRAY:
		return "BINARY(5000)"
	case deviceTemplate.ENUM:
		return "SMALLINT"
	default:
		panic(fmt.Sprintf("%v not support", define.Type))
	}
}

func getPropertyStableName(productID, id string) string {
	return fmt.Sprintf("model_property_%s_%s", productID, id)
}
func getEventStableName(productID string) string {
	return fmt.Sprintf("model_event_%s", productID)
}

func getActionStableName(productID string) string {
	return fmt.Sprintf("model_action_%s", productID)
}

func getPropertyTableName(productID, deviceName, id string) string {
	return fmt.Sprintf("device_property_%s_%s_%s", productID, deviceName, id)
}
func getEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("device_event_%s_%s", productID, deviceName)
}

func getActionTableName(productID, deviceName string) string {
	return fmt.Sprintf("device_action_%s_%s", productID, deviceName)
}
