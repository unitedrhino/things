package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexchnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexchnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexchnLogic {
	return &IndexchnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// logic.ToDmPageRpc(req.Page),
func (l *IndexchnLogic) Indexchn(req *types.VidmgrSipIndexChnReq) (resp *types.VidmgrSipIndexChnResp, err error) {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrGbsipChannelIndexReq{
		ChannelIDs: req.ChannelIDs,
		Page:       logic.ToVidPageRpc(req.Page),
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Indexchn:", string(jsonStr))
	vidResp, err := l.svcCtx.VidmgrG.VidmgrGbsipChannelIndex(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.Indexdev req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	if len(vidResp.List) > 0 {
		pis := make([]*types.CommonSipChannel, 0, len(vidResp.List))
		for _, v := range vidResp.List {
			pi := VidmgrGbsipChanneloApi(v)
			pis = append(pis, pi)
		}
		return &types.VidmgrSipIndexChnResp{
			Total: vidResp.Total,
			List:  pis,
			Num:   int64(len(pis)),
		}, nil
	}
	return &types.VidmgrSipIndexChnResp{
		Total: 0,
		List:  nil,
		Num:   0,
	}, nil
	return
}
