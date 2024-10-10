package custom

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProductCustomReadReq) (resp *types.ProductCustom, err error) {
	dmResp, err := l.svcCtx.ProductM.ProductCustomRead(l.ctx, &dm.ProductCustomReadReq{ProductID: req.ProductID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ProductCustomRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.ProductCustom{
		ProductID:       dmResp.ProductID,
		TransformScript: utils.ToNullString(dmResp.TransformScript),
		LoginAuthScript: utils.ToNullString(dmResp.LoginAuthScript),
		ScriptLang:      dmResp.ScriptLang,
		CustomTopics:    ToCustomTopicsTypes(dmResp.CustomTopics),
	}, nil
}
