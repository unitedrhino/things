// 区域物模型属性统计工具
// 调用 dmsvr gRPC 服务，统计区域树下指定物模型属性的聚合值（min/max/sum），
// 并将结果写入目标设备的对应物模型属性中
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	dm "gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 解析命令行参数
	productID := flag.String("productID", "", "产品ID（必填）")
	deviceName := flag.String("deviceName", "", "目标设备名（必填，统计结果写入此设备）")
	areaID := flag.Int64("areaID", 0, "区域ID（必填，用于获取区域树）")
	projectID := flag.Int64("projectID", 0, "项目ID（必填）")
	areaIDPath := flag.String("areaIDPath", "", "区域ID路径（必填，如 1/2/3）")
	dataID := flag.String("dataID", "", "物模型属性标识符（必填）")
	statFunc := flag.String("statFunc", "", "统计函数: min / max / sum")
	dmsvrAddr := flag.String("dmsvrAddr", "localhost:7540", "dmsvr gRPC地址")
	timeWindow := flag.Int64("timeWindow", 600, "统计时间窗口（秒），默认600（10分钟）")
	flag.Parse()

	// 参数校验
	if *productID == "" || *deviceName == "" || *dataID == "" || *statFunc == "" {
		fmt.Fprintln(os.Stderr, "错误: --productID, --deviceName, --dataID, --statFunc 为必填参数")
		flag.Usage()
		os.Exit(1)
	}
	if *areaID == 0 && *areaIDPath == "" {
		fmt.Fprintln(os.Stderr, "错误: --areaID 和 --areaIDPath 至少填写一个")
		flag.Usage()
		os.Exit(1)
	}
	if *statFunc != "min" && *statFunc != "max" && *statFunc != "sum" {
		fmt.Fprintln(os.Stderr, "错误: --statFunc 仅支持 min / max / sum")
		os.Exit(1)
	}

	// 建立 gRPC 连接
	conn, err := grpc.NewClient(*dmsvrAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接 dmsvr 失败 (%s): %v\n", *dmsvrAddr, err)
		os.Exit(1)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 确定 areaIDPath（优先使用传入的路径，否则从区域树获取）
	finalAreaIDPath := *areaIDPath
	if *areaID > 0 && *projectID > 0 {
		areaClient := dm.NewAreaManageClient(conn)
		areaResp, err := areaClient.AreaInfoGetOne(ctx, &dm.AreaInfoGetOneReq{
			ProjectID:    *projectID,
			AreaID:       *areaID,
			WithChildren: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取区域信息失败: %v\n", err)
			os.Exit(1)
		}
		if finalAreaIDPath == "" {
			finalAreaIDPath = areaResp.AreaIDPath
		}
		fmt.Printf("区域[%d] %s, 路径: %s, 子区域数: %d\n",
			areaResp.AreaID, areaResp.AreaName, areaResp.AreaIDPath, len(areaResp.ChildrenAreaIDs))
	}

	if finalAreaIDPath == "" {
		fmt.Fprintln(os.Stderr, "错误: 无法确定 areaIDPath")
		os.Exit(1)
	}

	// 聚合查询：对区域树下的设备属性执行统计函数
	now := time.Now().UnixMilli()
	timeStart := now - *timeWindow*1000
	timeEnd := now

	msgClient := dm.NewDeviceMsgClient(conn)
	aggResp, err := msgClient.PropertyLogAggGetList(ctx, &dm.PropertyAggGetListReq{
		ProductID:  *productID,
		AreaIDPath: finalAreaIDPath,
		TimeStart:  timeStart,
		TimeEnd:    timeEnd,
		Aggs: []*dm.PropertyAgg{
			{
				DataID:   *dataID,
				ArgFuncs: []string{*statFunc},
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "聚合查询失败: %v\n", err)
		os.Exit(1)
	}

	// 提取聚合结果值
	value := extractAggValue(aggResp, *dataID, *statFunc)
	if value == "" {
		fmt.Println("统计结果为空（区域下无设备或无数据），不执行写入")
		os.Exit(0)
	}

	deviceCount := len(aggResp.List)
	fmt.Printf("统计结果: dataID=%s statFunc=%s value=%s 参与设备数=%d\n",
		*dataID, *statFunc, value, deviceCount)
	fmt.Printf("时间窗口: %s ~ %s\n",
		time.UnixMilli(timeStart).Format("15:04:05"),
		time.UnixMilli(timeEnd).Format("15:04:05"))

	// 写回目标设备物模型
	interactClient := dm.NewDeviceInteractClient(conn)
	writeData := fmt.Sprintf(`{"%s":%s}`, *dataID, value)
	_, err = interactClient.PropertyControlSend(ctx, &dm.PropertyControlSendReq{
		ProductID:     *productID,
		DeviceName:    *deviceName,
		Data:          writeData,
		ShadowControl: 3, // 仅修改云端值，不下发设备
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "写回设备物模型失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("写入成功: device=%s/%s data=%s\n", *productID, *deviceName, writeData)
}

// extractAggValue 从聚合查询响应中提取指定统计函数的结果值
func extractAggValue(resp *dm.PropertyLogAggGetListResp, dataID, statFunc string) string {
	for _, item := range resp.List {
		for _, detail := range item.Values {
			if detail.DataID == dataID {
				for fn, v := range detail.Values {
					if fn == statFunc && v != nil {
						return v.Value
					}
				}
			}
		}
	}
	return ""
}
