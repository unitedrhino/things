package svc

import (
	"context"
	"github.com/i-Things/things/shared/db/mongodb"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/repo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mongorepo"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
)

type ServiceContext struct {
	Config          config.Config
	DeviceInfo      mysql.DeviceInfoModel
	ProductInfo     mysql.ProductInfoModel
	ProductTemplate mysql.ProductTemplateModel
	DeviceLog       mysql.DeviceLogModel
	DmDB            mysql.DmModel
	DeviceID        *utils.SnowFlake
	ProductID       *utils.SnowFlake
	DeviceData      repo.GetDeviceDataRepo
	InnerLink       innerLink.InnerLink
	Store           kv.Store
	Td              *TDengine.Td
}

//func TestTD(taos *TDengine.Td) {
//	taos.Exec("create database if not exists test")
//	taos.Exec("create table if not exists tb1 (ts timestamp, a int)")
//	_, err := taos.Exec("insert into tb1 values(now, 0)(now+1s,1)(now+2s,2)(now+3s,3)")
//	if err != nil {
//		fmt.Println("failed to insert, err:", err)
//		return
//	}
//	rows, err := taos.Query("select * from tb1")
//	if err != nil {
//		fmt.Println("failed to select from table, err:", err)
//		return
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var r struct {
//			ts time.Time
//			a  int
//		}
//		err := rows.Scan(&r.ts, &r.a)
//		if err != nil {
//			fmt.Println("scan error:\n", err)
//			return
//		}
//		fmt.Println("get data:", r.ts, r.a)
//	}
//}

func NewServiceContext(c config.Config) *ServiceContext {
	td, err := TDengine.NewTDengine(c.TDengine.DataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	//TestTD(td)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := mysql.NewDeviceInfoModel(conn, c.CacheRedis)
	pi := mysql.NewProductInfoModel(conn, c.CacheRedis)
	pt := mysql.NewProductTemplateModel(conn, c.CacheRedis)
	dl := mysql.NewDeviceLogModel(conn)
	DmDB := mysql.NewDmModel(conn, c.CacheRedis)
	store := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	DeviceID := utils.NewSnowFlake(nodeId)
	ProductID := utils.NewSnowFlake(nodeId)
	mongoDB, err := mongodb.NewMongo(c.Mongo.Url, c.Mongo.Database, context.TODO())
	if err != nil {
		logx.Error("NewMongo err", err)
		os.Exit(-1)
	}
	dd := mongorepo.NewDeviceDataRepo(mongoDB)
	il, err := innerLink.NewInnerLink(c.InnerLink)
	if err != nil {
		logx.Error("NewInnerLink err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:          c,
		DeviceInfo:      di,
		ProductInfo:     pi,
		ProductTemplate: pt,
		DmDB:            DmDB,
		DeviceID:        DeviceID,
		ProductID:       ProductID,
		DeviceLog:       dl,
		DeviceData:      dd,
		InnerLink:       il,
		Store:           store,
		Td:              td,
	}
}
