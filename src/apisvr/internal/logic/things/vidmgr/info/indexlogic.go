package info

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

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

func (l *IndexLogic) Index(req *types.VidmgrInfoIndexReq) (resp *types.VidmgrInfoIndexResp, err error) {
	// todo: add your logic here and delete this line
	fmt.Println("[--airgens-]", req)
	vidReq := &vid.VidmgrInfoIndexReq{
		VidmgrType:  req.VidmgrType,
		VidmgrIDs:   req.VidmgrIDs,
		VidmgrtName: req.VidmgrName,
		Tags:        logic.ToTagsMap(req.Tags),
		Page:        logic.ToVidPageRpc(req.Page),
	}
	vidResp, err := l.svcCtx.VidmgrM.VidmgrInfoIndex(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrInfoIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if len(vidResp.List) > 0 {
		pis := make([]*types.VidmgrInfo, 0, len(vidResp.List))
		for _, v := range vidResp.List {
			pi := VidmgrInfoToApi(v)
			pis = append(pis, pi)
		}
		return &types.VidmgrInfoIndexResp{
			Total: vidResp.Total,
			List:  pis,
			Num:   int64(len(pis)),
		}, nil
	} else {
		return &types.VidmgrInfoIndexResp{
			Total: 0,
			List:  nil,
			Num:   0,
		}, nil
	}
}
