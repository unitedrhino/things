package remoteconfiglogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PrcDB *relationDB.ProductRemoteConfigRepo
}

func NewRemoteConfigIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigIndexLogic {
	return &RemoteConfigIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PrcDB:  relationDB.NewProductRemoteConfigRepo(ctx),
	}
}

func (l *RemoteConfigIndexLogic) RemoteConfigIndex(in *dm.RemoteConfigIndexReq) (*dm.RemoteConfigIndexResp, error) {
	f := relationDB.RemoteConfigFilter{ProductID: in.ProductID}
	rcs, err := l.PrcDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	list := make([]*dm.ProductRemoteConfig, 0, len(rcs))
	for _, v := range rcs {
		list = append(list, &dm.ProductRemoteConfig{
			Id:          v.ID,
			ProductID:   v.ProductID,
			Content:     v.Content,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	total, err := l.PrcDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &dm.RemoteConfigIndexResp{List: list, Total: total}, nil
}
