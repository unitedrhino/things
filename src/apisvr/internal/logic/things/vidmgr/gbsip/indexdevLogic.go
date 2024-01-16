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

type IndexdevLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexdevLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexdevLogic {
	return &IndexdevLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexdevLogic) Indexdev(req *types.VidmgrSipIndexDevReq) (resp *types.VidmgrSipIndexDevResp, err error) {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrGbsipDeviceIndexReq{
		DeviceIDs: req.DeviceIDs,
		Page:      logic.ToVidPageRpc(req.Page),
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Indexdev:", string(jsonStr))
	vidResp, err := l.svcCtx.VidmgrG.VidmgrGbsipDeviceIndex(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.Indexdev req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	if len(vidResp.List) > 0 {
		pis := make([]*types.CommonSipDevice, 0, len(vidResp.List))
		for _, v := range vidResp.List {
			pi := VidmgrGbsipDeviceToApi(v)
			pis = append(pis, pi)
		}
		return &types.VidmgrSipIndexDevResp{
			Total: vidResp.Total,
			List:  pis,
			Num:   int64(len(pis)),
		}, nil
	}
	return &types.VidmgrSipIndexDevResp{
		Total: 0,
		List:  nil,
		Num:   0,
	}, nil
}
