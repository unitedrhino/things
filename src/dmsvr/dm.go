package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/i-Things/things/src/dmsvr/internal/event/eventDevSub"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/dmsvr/internal/server"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "net/http/pprof"
	"time"
)

var configFile = flag.String("f", "etc/dm.yaml", "the config file")

func main() {
	flag.Parse()
	//go device.NewDevice()
	//device.TestMongo()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	Test(svcCtx.DeviceDataRepo)
	svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubHandle {
		return eventDevSub.NewDeviceMsgHandle(ctx, svcCtx)
	})
	//grpc服务初始化

	srv := server.NewDmServer(svcCtx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dm.RegisterDmServer(grpcServer, srv)
		reflection.Register(grpcServer)
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
func Test(DeviceDataRepo deviceTemplate.DeviceData2Repo) {
	var (
		err error
	)

	if false {
		err = DeviceDataRepo.InsertEventData(context.Background(),
			"23FIPSIJPsk", "test5", &deviceTemplate.EventData{
				ID:        "faw",
				Type:      "info",
				Params:    map[string]interface{}{"hello": 123123},
				TimeStamp: time.Now(),
			})
		fmt.Println(err)
	}
	if false {
		err = DeviceDataRepo.InsertPropertyData(context.Background(), "23FIPSIJPsk", "test5", &deviceTemplate.PropertyData{
			ID:        "Wifi_Info",
			Param:     []interface{}{map[string]interface{}{"Mac": "dqwda", "Rssi": 123}},
			TimeStamp: time.Now(),
		})
		fmt.Println(err)
		err = DeviceDataRepo.InsertPropertyData(context.Background(), "23FIPSIJPsk", "test5", &deviceTemplate.PropertyData{
			ID:        "GPS_Info",
			Param:     map[string]interface{}{"longtitude": 12.44, "latitude": 22.987},
			TimeStamp: time.Now(),
		})
		fmt.Println(err)
	}
	if false {
		_, err = DeviceDataRepo.GetPropertyDataByID(context.Background(), "23FIPSIJPsk", "test5", "Wifi_Info", def.PageInfo2{})
		fmt.Println(err)
		_, err = DeviceDataRepo.GetPropertyDataByID(context.Background(), "23FIPSIJPsk", "test5", "GPS_Info", def.PageInfo2{})
		fmt.Println(err)
	}

}
