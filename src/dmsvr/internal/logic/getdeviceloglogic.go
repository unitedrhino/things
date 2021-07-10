package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/device"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceLogLogic {
	return &GetDeviceLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func GetFileName(in *dm.GetDeviceLogReq) bson.D{
	//property 属性 event事件 action 操作 log 所有日志
	switch in.Method {
	case "property":
		return bson.D{
			//{"isp", isp},
			{fmt.Sprintf("%s.%s",in.Method,in.FieldName),  bson.M{"$ne": primitive.Null{}}},
		}

	case "event":
		return bson.D{
			{fmt.Sprintf("%s.id",in.Method),  bson.M{"$eq": in.FieldName}},
		}
	}
	return nil
}

func (l *GetDeviceLogLogic) HandleData(in *dm.GetDeviceLogReq) (*dm.GetDeviceLogResp, error) {
	clientID := fmt.Sprintf("%s%s",in.ProductID,in.DeviceName)
	filter := GetFileName(in)
	if in.TimeStart != 0 {
		filter = append(filter,bson.E{"timestamp",bson.M{"$gte":time.Unix(in.TimeStart,0)}})
	}
	if in.TimeEnd != 0 {
		filter = append(filter,bson.E{"timestamp",bson.M{"$lte":time.Unix(in.TimeEnd,0)}})
	}
	opts :=options.Find().SetProjection(bson.D{{"timestamp", 1},{in.Method,1}}).
		SetLimit(in.Limit).SetSort(bson.D{{"timestamp", -1}})
	ctx,_ := context.WithTimeout(l.ctx,5*time.Second)
	cursor,err := l.svcCtx.Mongo.Collection(clientID).Find(ctx,filter,opts)
	if err != nil {
		l.Errorf("Find|err=%v",err)
		return nil, errors.System
	}
	defer cursor.Close(ctx)
	resp := dm.GetDeviceLogResp{}
	for cursor.Next(ctx) {
		err = cursor.Err()
		if err != nil {
			l.Errorf("cursor|err=%v",err)
			return nil, errors.System
		}
		var original device.DeviceData
		cursor.Decode(&original)
		dd := dm.DeviceData{
			Timestamp: original.TimeStamp.Unix(),
			Method: in.Method,
			FieldName: in.FieldName,
		}
		var payload []byte
		switch in.Method {
		case "property":
			payload,_ =json.Marshal(original.Property[in.FieldName])
		case "event":
			payload,_ =json.Marshal(original.Event)
		}
		dd.Payload = string(payload)
		resp.Data = append(resp.Data,&dd)
		l.Infof("coursor=%+v",original)
	}
	resp.Total = int64(len(resp.Data))
	return &resp, nil
}

func (l *GetDeviceLogLogic) GetDeviceLog(in *dm.GetDeviceLogReq) (*dm.GetDeviceLogResp, error) {
	switch in.Method {
	case "property","action","event"://获取属性信息,获取操作信息,获取事件信息
		return l.HandleData(in)
	case "status"://获取设备状态信息
	case "logs"://获取设备的调试日志
	default:
		return nil, errors.Method.AddDetail(in.Method)
	}

	return &dm.GetDeviceLogResp{}, nil
}
