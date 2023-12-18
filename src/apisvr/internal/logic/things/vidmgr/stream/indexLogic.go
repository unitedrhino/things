package stream

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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

func (l *IndexLogic) Index(req *types.VidmgrStreamIndexReq) (resp *types.VidmgrStreamIndexResp, err error) {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrStreamIndexReq{
		VidmgrID:   req.VidmgrID,
		StreamName: req.StreamName,
		StreamIDs:  req.StreamIDs,
		Tags:       logic.ToTagsMap(req.Tags),
		Page:       logic.ToVidPageRpc(req.Page),
	}
	//tmpByte, _ := json.Marshal(req)
	//fmt.Println("HttpReq_VidmgrStreamIndex:", string(tmpByte))
	//
	//tmpByte1, _ := json.Marshal(vidReq)
	//fmt.Println("VidReq_VidmgrStreamIndex:", string(tmpByte1))

	vidResp, err := l.svcCtx.VidmgrS.VidmgrStreamIndex(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrStreamIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if len(vidResp.List) > 0 {
		pis := make([]*types.VidmgrStream, 0, len(vidResp.List))
		for _, v := range vidResp.List {
			pi := VidmgrStreamToApi(v)
			pis = append(pis, pi)
		}
		return &types.VidmgrStreamIndexResp{
			Total: vidResp.Total,
			List:  pis,
			Num:   int64(len(pis)),
		}, nil
	} else {
		return &types.VidmgrStreamIndexResp{
			Total: 0,
			List:  nil,
			Num:   0,
		}, nil
	}
}
