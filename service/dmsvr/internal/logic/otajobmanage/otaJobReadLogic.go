package otajobmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaJobReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewOtaJobReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaJobReadLogic {
	return &OtaJobReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 查询指定升级批次的详情
func (l *OtaJobReadLogic) OtaJobRead(in *dm.JobReq) (*dm.OtaJobInfo, error) {
	otaJob, err := l.OjDB.FindOne(l.ctx, in.JobId)
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	var otaJobInfo *dm.OtaJobInfo
	err = copier.Copy(&otaJobInfo, &otaJob)
	if err != nil {
		l.Errorf("%s.JobInfo.CopyJobInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return otaJobInfo, nil
}
