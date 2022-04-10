package mongorepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	DeviceData struct {
		mongo *mongo.Database
	}
	DeviceDataContext struct {
		DeviceData
		ctx context.Context
	}

	// Event 数据库模型
	Event struct {
		ID         string                 `json:"id" bson:"id"`                 //事件id
		Type       string                 `json:"type" bson:"type"`             //事件类型: 信息:info  告警alert  故障:fault
		Params     map[string]interface{} `json:"params" bson:"params"`         //事件参数
		DeviceName string                 `json:"deviceName" bson:"deviceName"` //设备名称
		TimeStamp  time.Time              `json:"timeStamp" bson:"timestamp"`   //时间戳
	}

	// Properties 数据库模型
	Properties struct {
		Params     map[string]interface{} `json:"properties" bson:"property"`   //一个属性的参数
		DeviceName string                 `json:"deviceName" bson:"deviceName"` //设备名称
		TimeStamp  time.Time              `json:"timeStamp" bson:"timestamp"`   //时间戳
	}
)

const (
	//时序数据库是时间戳key
	TimeStampKey = "timestamp"
	//mongodb 时序数据库 时间戳别名
	PropertyMD = "property"
	EventMD    = "id"
)

const (
	dbSuffixProperty = "_property"
	dbSuffixEvent    = "_event"
)

func NewDeviceDataRepo(Mongo *mongo.Database) deviceTemplate.GetDeviceDataRepo {
	return func(ctx context.Context) deviceTemplate.DeviceDataRepo {
		return &DeviceDataContext{
			DeviceData: DeviceData{
				mongo: Mongo,
			},
			ctx: ctx,
		}
	}
}

func (d *DeviceDataContext) InsertEventData(productID string, deviceName string, event *deviceTemplate.EventData) error {
	model := Event{
		ID:         event.ID,
		Type:       event.Type,
		Params:     event.Params,
		DeviceName: deviceName,
		TimeStamp:  event.TimeStamp,
	}
	_, err := d.mongo.Collection(productID+dbSuffixEvent).InsertOne(d.ctx, model)
	return err
}
func (d *DeviceDataContext) InsertPropertyData(productID string, deviceName string, property *deviceTemplate.PropertyData) error {
	dd := Properties{
		DeviceName: deviceName,
		TimeStamp:  time.Time{},
	}
	dd.Params[property.ID] = property.Param
	_, err := d.mongo.Collection(productID+dbSuffixProperty).InsertOne(d.ctx, dd)
	return err
}

func (d *DeviceDataContext) InsertPropertiesData(productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error {
	dd := Properties{
		DeviceName: deviceName,
		TimeStamp:  timestamp,
		Params:     params,
	}
	_, err := d.mongo.Collection(productID+dbSuffixProperty).InsertOne(d.ctx, dd)
	return err
}

func (d *DeviceDataContext) CreatePropertyDB(productID string) error {
	opt := options.CreateCollection()
	opt.SetTimeSeriesOptions(options.TimeSeries().SetTimeField(TimeStampKey).SetMetaField("properties")).
		SetExpireAfterSeconds(int64(time.Hour * 24 * 30 * 12 * 2 / time.Second))
	return d.mongo.CreateCollection(d.ctx, productID+dbSuffixProperty, opt)
}
func (d *DeviceDataContext) CreateEventDB(productID string) error {
	opt := options.CreateCollection()
	opt.SetTimeSeriesOptions(options.TimeSeries().SetTimeField(TimeStampKey).SetMetaField(EventMD)).
		SetExpireAfterSeconds(int64(time.Hour * 24 * 30 * 12 * 2 / time.Second))
	return d.mongo.CreateCollection(d.ctx, productID+dbSuffixEvent, opt)
}

//todo 暂时使用mysql存储日志
func (d *DeviceDataContext) CreateLogDB(productID string) error {
	opt := options.CreateCollection()
	opt.SetTimeSeriesOptions(options.TimeSeries().SetTimeField(TimeStampKey).SetMetaField(PropertyMD)).
		SetExpireAfterSeconds(int64(time.Hour * 24 * 30 * 3 / time.Second))
	return d.mongo.CreateCollection(d.ctx, productID+dbSuffixEvent, opt)
}

//通过属性的id及方法获取一段时间或最新时间的记录
func (d *DeviceDataContext) GetEventDataWithID(productID string, deviceName string, dataID string, page def.PageInfo2) (dds []*deviceTemplate.EventData, err error) {
	filter := bson.D{
		{"deviceName", bson.M{"$eq": deviceName}},
	}
	if dataID != "" {
		filter = append(filter, bson.E{Key: "id", Value: bson.M{"$eq": dataID}})
	}
	if page.TimeStart != 0 {
		filter = append(filter, bson.E{TimeStampKey, bson.M{"$gte": time.UnixMilli(page.TimeStart)}})
	}
	if page.TimeEnd != 0 {
		filter = append(filter, bson.E{TimeStampKey, bson.M{"$lte": time.UnixMilli(page.TimeEnd)}})
	}
	opts := options.Find().SetProjection(bson.D{{TimeStampKey, 1}}).
		SetLimit(page.Limit).SetSort(bson.D{{TimeStampKey, -1}})
	cursor, err := d.mongo.Collection(productID+dbSuffixEvent).Find(d.ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(d.ctx)
	for cursor.Next(d.ctx) {
		err = cursor.Err()
		if err != nil {
			return nil, err
		}
		var dd Event
		err = cursor.Decode(&dd)
		if err != nil {
			return nil, err
		}
		dds = append(dds, &deviceTemplate.EventData{
			ID:        dd.ID,
			Type:      dd.Type,
			Params:    dd.Params,
			TimeStamp: dd.TimeStamp,
		})
	}
	return dds, nil
}

//通过属性的id及方法获取一段时间或最新时间的记录
func (d *DeviceDataContext) GetPropertyDataWithID(productID string, deviceName string, dataID string, page def.PageInfo2) (dds []*deviceTemplate.PropertyData, err error) {
	filter := bson.D{
		//{"isp", isp},
		{fmt.Sprintf("%s.%s", PropertyMD, dataID), bson.M{"$ne": primitive.Null{}}},
		{"deviceName", bson.M{"$eq": deviceName}},
	}
	if page.TimeStart != 0 {
		filter = append(filter, bson.E{TimeStampKey, bson.M{"$gte": time.UnixMilli(page.TimeStart)}})
	}
	if page.TimeEnd != 0 {
		filter = append(filter, bson.E{TimeStampKey, bson.M{"$lte": time.UnixMilli(page.TimeEnd)}})
	}
	opts := options.Find().SetLimit(page.GetLimit()).SetSort(bson.D{{TimeStampKey, -1}})
	cursor, err := d.mongo.Collection(productID+dbSuffixProperty).Find(d.ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(d.ctx)
	for cursor.Next(d.ctx) {
		err = cursor.Err()
		if err != nil {
			return nil, err
		}
		var dd Properties
		err = cursor.Decode(&dd)
		if err != nil {
			return nil, err
		}
		dds = append(dds, &deviceTemplate.PropertyData{
			ID:        dataID,
			Param:     dd.Params[dataID],
			TimeStamp: dd.TimeStamp,
		})
	}
	return dds, nil
}
