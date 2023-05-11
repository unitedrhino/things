package custom

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
		ScriptLang:      dmResp.ScriptLang,
		CustomTopic:     dmResp.CustomTopic,
	}, nil
}
