package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareJobIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewOtaFirmwareJobIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareJobIndexLogic {
	return &OtaFirmwareJobIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// //创建动态升级批次
func (l *OtaFirmwareJobIndexLogic) OtaFirmwareJobIndex(in *dm.OtaFirmwareJobIndexReq) (*dm.OtaFirmwareJobIndexResp, error) {
	jobFilter := relationDB.OtaJobFilter{
		FirmwareId: in.FirmwareID,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	total, err := l.OjDB.CountByFilter(l.ctx, jobFilter)
	if err != nil {
		return nil, err
	}
	otaJobList, err := l.OjDB.FindByFilter(l.ctx, jobFilter, logic.ToPageInfo(in.PageInfo))
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	var list []*dm.OtaFirmwareJobInfo
	for _, v := range otaJobList {
		var otaJobInfo = dm.OtaFirmwareJobInfo{Dynamic: &dm.OtaJobDynamicInfo{}, Static: &dm.OtaJobStaticInfo{}}
		copier.Copy(&otaJobInfo, v)
		list = append(list, &otaJobInfo)
	}
	return &dm.OtaFirmwareJobIndexResp{List: list, Total: total}, nil
}
