//本文件是提供设备模型数据存储的信息
package repo

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/def"
	"gitee.com/godLei6/things/src/dmsvr/internal/repo/model"
	"gitee.com/godLei6/things/src/dmsvr/internal/vars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DeviceData struct {
	ctx context.Context
	DBName string
}

func NewDeviceData(ctx context.Context,dbName string)*DeviceData{
	return &DeviceData{ctx,dbName}
}

func (d *DeviceData)InsertEventData(event *model.Event)error{
	dd := model.DeviceData{
		Property: nil,
		Event: struct {
			ID     string                 `json:"id"`
			Type   string                 `json:"type"`
			Params map[string]interface{} `json:"params"`
		}{
			event.ID,
			event.Type,
			event.Params,
		},
		TimeStamp: event.TimeStamp,
	}
	_,err := vars.Svrctx.Mongo.Collection(d.DBName).InsertOne(d.ctx, dd)
	return err
}
func (d *DeviceData) InsertPropertyData(property *model.Properties)error{
	dd := model.DeviceData{
		Property: property.Params,
		TimeStamp: property.TimeStamp,
	}
	_,err := vars.Svrctx.Mongo.Collection(d.DBName).InsertOne(d.ctx, dd)
	return err
}


func getFileName(method,dataID string) bson.D {
	//property 属性 event事件 action 操作 log 所有日志
	switch method {
	case def.PROPERTY_METHOD:
		return bson.D{
			//{"isp", isp},
			{fmt.Sprintf("%s.%s", method, dataID), bson.M{"$ne": primitive.Null{}}},
		}
	case def.EVENT_METHOD:
		return bson.D{
			{fmt.Sprintf("%s.id", method), bson.M{"$eq": dataID}},
		}
	}
	return nil
}


//通过属性的id及方法获取一段时间或最新时间的记录
func (d *DeviceData)GetEventDataWithID(dataID string,timeStart,timeEnd int64,limit int64)(dds []*model.Event,err error) {
	filter := getFileName(def.EVENT_METHOD,dataID)
	if timeStart != 0 {
		filter = append(filter, bson.E{"timestamp", bson.M{"$gte": time.Unix(timeStart, 0)}})
	}
	if timeEnd != 0 {
		filter = append(filter, bson.E{"timestamp", bson.M{"$lte": time.Unix(timeEnd, 0)}})
	}
	opts := options.Find().SetProjection(bson.D{{"timestamp", 1}, {def.EVENT_METHOD, 1}}).
		SetLimit(limit).SetSort(bson.D{{"timestamp", -1}})
	cursor, err := vars.Svrctx.Mongo.Collection(d.DBName).Find(d.ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(d.ctx)
	for cursor.Next(d.ctx) {
		err = cursor.Err()
		if err != nil {
			return nil, err
		}
		var dd model.DeviceData
		err = cursor.Decode(&dd)
		if err != nil {
			return nil, err
		}
		dds = append(dds,&model.Event{
			ID:dd.Event.ID,
			Type:dd.Event.Type,
			Params:dd.Event.Params,
			TimeStamp:dd.TimeStamp,
		})
	}
	return dds,nil
}


//通过属性的id及方法获取一段时间或最新时间的记录
func (d *DeviceData)GetPropertyDataWithID(dataID string,timeStart,timeEnd int64,limit int64)(dds []*model.Property,err error) {
	filter := getFileName(def.PROPERTY_METHOD,dataID)
	if timeStart != 0 {
		filter = append(filter, bson.E{"timestamp", bson.M{"$gte": time.Unix(timeStart, 0)}})
	}
	if timeEnd != 0 {
		filter = append(filter, bson.E{"timestamp", bson.M{"$lte": time.Unix(timeEnd, 0)}})
	}
	opts := options.Find().SetProjection(bson.D{{"timestamp", 1}, {def.PROPERTY_METHOD, 1}}).
		SetLimit(limit).SetSort(bson.D{{"timestamp", -1}})
	cursor, err := vars.Svrctx.Mongo.Collection(d.DBName).Find(d.ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(d.ctx)
	for cursor.Next(d.ctx) {
		err = cursor.Err()
		if err != nil {
			return nil, err
		}
		var dd model.DeviceData
		err = cursor.Decode(&dd)
		if err != nil {
			return nil, err
		}
		dds = append(dds,&model.Property{
			Param:     dd.Property[dataID],
			TimeStamp: dd.TimeStamp,
		})
	}
	return dds,nil
}

