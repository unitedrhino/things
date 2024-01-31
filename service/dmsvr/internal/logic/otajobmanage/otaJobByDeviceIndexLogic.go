package otajobmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaJobByDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewOtaJobByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaJobByDeviceIndexLogic {
	return &OtaJobByDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 获取设备所在的升级包升级批次列表
func (l *OtaJobByDeviceIndexLogic) OtaJobByDeviceIndex(in *dm.OtaJobByDeviceIndexReq) (*dm.OtaJobInfoIndexResp, error) {
	filter := relationDB.OtaJobFilter{
		ProductId:  in.ProductId,
		DeviceName: in.DeviceName,
	}
	total, err := l.OjDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	var otaJobInfoList []*dm.OtaJobInfo
	otaJobList, err := l.OjDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.PageInfo))
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	err = copier.Copy(&otaJobInfoList, &otaJobList)
	if err != nil {
		l.Errorf("%s.JobInfo.CopyJobInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.OtaJobInfoIndexResp{Total: total, OtaJobInfo: otaJobInfoList}, nil
}
