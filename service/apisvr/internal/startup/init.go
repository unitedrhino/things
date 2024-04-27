package startup

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
)

func StartOtaChanWalk(s *svc.ServiceContext) {
	if s.Config.DmRpc.Enable {
		utils.Go(context.Background(), func() {
			FileChanWalk(s)
		})
	}
}
func init() {
	var (
		TagsTypes []*types.Tag
		TagMap    map[string]string
	)
	utils.AddConverter(
		utils.TypeConverter{SrcType: TagsTypes, DstType: TagMap, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsMap(src.([]*types.Tag)), nil
		}},
		utils.TypeConverter{SrcType: TagMap, DstType: TagsTypes, Fn: func(src interface{}) (dst interface{}, err error) {
			return logic.ToTagsType(src.(map[string]string)), nil
		}},
	)

}

func FileChanWalk(s *svc.ServiceContext) {
	//ctx := context.Background()
	////处理因为宕机未执行的file
	//old := &firmwaremanage.OtaFirmwareFileIndexReq{
	//	Size: &wrapperspb.Int64Value{
	//		Value: 0,
	//	},
	//}
	//fileList, err := s.FirmwareM.OtaFirmwareFileIndex(ctx, old)
	//if err != nil {
	//	logx.Errorf("%v.OtaFirmwareFileIndex err:%v", utils.FuncName(), err)
	//	return
	//}
	//for _, f := range fileList.List {
	//	s.FileChan <- f.FirmwareID
	//}
	////chan
	//for {
	//	firmwareID := <-s.FileChan
	//	in := &firmwaremanage.FirmwareInfoReadReq{
	//		FirmwareID: firmwareID,
	//	}
	//	firmwareInfo, err := s.FirmwareM.FirmwareInfoRead(ctx, in)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	for _, f := range firmwareInfo.Files {
	//		//fmt.Println(f.FileUrl, " chain dir")
	//		storageInfo, _ := s.OssClient.PrivateBucket().GetObjectInfo(context.Background(), f.FileUrl)
	//
	//		fileIn := &firmwaremanage.OtaFirmwareFileReq{
	//			FileID:    f.FileID,
	//			Size:      storageInfo.Size,
	//			Signature: storageInfo.Md5,
	//		}
	//		s.FirmwareM.OtaFirmwareFileUpdate(context.Background(), fileIn)
	//	}
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//}
}
