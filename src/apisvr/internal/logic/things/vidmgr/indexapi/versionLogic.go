package indexapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VersionLogic {
	return &VersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VersionLogic) Version(req *types.IndexApiReq) (resp *types.IndexApiVersionResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, VERSION, req.VidmgrID)
	if err != nil {
		fmt.Println("***proxyMediaServer Error ****")
		er := errors.Fmt(err)
		fmt.Print("%s proxyMediaServer  err=%+v", utils.FuncName(), er)
		return nil, er
	}
	dataRecv := new(types.IndexApiVersionResp)
	fmt.Println(string(data))
	fmt.Println(dataRecv)
	json.Unmarshal(data, dataRecv)
	return dataRecv, nil
}
