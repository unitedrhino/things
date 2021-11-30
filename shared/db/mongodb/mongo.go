package mongodb

import (
	"context"
	"fmt"
	_ "github.com/tal-tech/go-zero/core/stores/mongo"
	"go.mongodb.org/mongo-driver/bson" //BOSN解析包
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" //MongoDB的Go驱动包
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoZero() {

	//model,err:=mongo.NewModel(
	//	mongoUrl,"test2")
	//fmt.Printf("model=%+v,err=%+v\n",model,err)
	//if err != nil {
	//	log.Fatalf("NewModel fail:%s",err.Error())
	//}
	//doc  :=make(map[string]string,5)
	//
	//doc["name"] = "wefaefag"
	//err = model.Insert(doc)
	//if err != nil {
	//	log.Fatalf("Insert fail:%s",err.Error())
	//}
	//q,err := model.Find("{}")
	//if err != nil {
	//	log.Fatalf("find fail:%s",err.Error())
	//}
	//var rst interface{}
	//err = q.All(&rst)
	//if err != nil {
	//	log.Fatalf("find all fail:%s",err.Error())
	//}
	//fmt.Printf("find:%+v\n",rst)
}

func NewMongo(mongoUrl string, database string, ctx context.Context) (*mongo.Database, error) {
	clientOpt := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		return nil, err
	}
	// 检查连接情况
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client.Database(database), nil
}

type CurlInfo struct {
	DNS float64 `json:"NAMELOOKUP_TIME"` //NAMELOOKUP_TIME
	TCP float64 `json:"CONNECT_TIME"`    //CONNECT_TIME - DNS
	SSL float64 `json:"APPCONNECT_TIME"` //APPCONNECT_TIME - CONNECT_TIME
}

type ConnectData struct {
	Latency  float64  `json:"latency"`
	RespCode int      `json:"respCode"`
	Url      string   `json:"url"`
	Detail   CurlInfo `json:"details"`
}

type Sensor struct {
	ISP       string
	Clientutc int64
	DataByAPP map[string]ConnectData
}

func insertSensor(client *mongo.Client, collection *mongo.Collection) (insertID primitive.ObjectID) {
	apps := make(map[string]ConnectData, 0)
	apps["app1"] = ConnectData{
		Latency:  30.983999967575,
		RespCode: 200,
		Url:      "",
		Detail: CurlInfo{
			DNS: 5.983999967575,
			TCP: 10.983999967575,
			SSL: 15.983999967575,
		},
	}

	record := &Sensor{
		Clientutc: time.Now().UTC().Unix(),
		ISP:       "China Mobile",
		DataByAPP: apps,
	}

	insertRest, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println(err)
		return
	}

	insertID = insertRest.InsertedID.(primitive.ObjectID)
	return insertID
}

func querySensor(collection *mongo.Collection, isp string) {
	//查询一条记录

	//筛选数据
	timestamp := time.Now().UTC().Unix()
	start := timestamp - 1800
	end := timestamp

	filter := bson.D{
		//{"isp", isp},
		{"$and", bson.A{
			bson.D{{"clientutc", bson.M{"$gte": start}}},
			bson.D{{"clientutc", bson.M{"$lte": end}}},
		}},
	}

	var original map[string]interface{}
	err := collection.FindOne(context.TODO(), filter).Decode(&original)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("Found a single document,type=%T: %+v\n", original, original)
}
