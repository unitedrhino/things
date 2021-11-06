package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)
const mongoUrl = "mongodb://admin:password@127.0.0.1:27017"

func TestMongo(t *testing.T) {
	clientOpt := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	// 检查连接情况
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB!")

	//指定要操作的数据集
	collection := client.Database("things").Collection("test4")

	//执行增删改查操作
	//insertSensor(client, collection)

	querySensor(collection, "")
	// 断开客户端连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
