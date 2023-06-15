package devicemsglogic

import (
	"context"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/disvr/internal/domain/shadow"
	"github.com/i-Things/things/src/disvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShadowIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShadowIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShadowIndexLogic {
	return &ShadowIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备影子列表
func (l *ShadowIndexLogic) ShadowIndex(in *di.PropertyLatestIndexReq) (*di.ShadowIndexResp, error) {
	sr := relationDB.NewShadowRepo(l.ctx)
	srs, err := sr.FindByFilter(l.ctx, shadow.Filter{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		DataIDs:    in.DataIDs,
	})
	if err != nil {
		return nil, err
	}
	var index []*di.ShadowIndex
	for _, v := range srs {
		index = append(index, &di.ShadowIndex{
			DataID:            v.DataID,
			Value:             v.Value,
			UpdatedDeviceTime: cast.ToInt64(v.UpdatedDeviceTime),
		})
	}
	return &di.ShadowIndexResp{List: index}, nil
}
