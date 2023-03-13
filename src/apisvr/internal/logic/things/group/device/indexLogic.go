package device

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func toTagsType(tags map[string]string) (retTag []*types.Tag) {
	for k, v := range tags {
		retTag = append(retTag, &types.Tag{
			Key:   k,
			Value: v,
		})
	}
	return
}
func (l *IndexLogic) Index(req *types.GroupDeviceIndexReq) (resp *types.GroupDeviceIndexResp, err error) {
	var list []*types.DeviceInfo
	var page dm.PageInfo
	copier.Copy(&page, req.Page)
	gd, err := l.svcCtx.DeviceG.GroupDeviceIndex(l.ctx, &dm.GroupDeviceIndexReq{
		Page:       &page,
		GroupID:    req.GroupID,
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceGroup GroupDeviceIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, err
	}
	for _, v := range gd.List {
		list = append(list, &types.DeviceInfo{
			ProductID:   v.ProductID,
			DeviceName:  v.DeviceName,
			CreatedTime: v.CreatedTime,
			Secret:      v.Secret,
			FirstLogin:  v.FirstLogin,
			LastLogin:   v.LastLogin,
			//Version:     v.Version.String(),
			LogLevel: v.LogLevel,
			Cert:     v.Cert,
			Tags:     toTagsType(v.Tags),
			IsOnline: v.IsOnline,
		})
	}

	return &types.GroupDeviceIndexResp{
		List:  list,
		Total: gd.Total,
	}, nil
}
