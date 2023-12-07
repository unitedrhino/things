package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnServerStartedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnServerStartedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnServerStartedLogic {
	return &OnServerStartedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnServerStartedLogic) OnServerStarted(req *types.HooksApiServerStartedReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	//用于侦听流服务重启，流服务重启后，当获得当前流服务的配置。
	fmt.Println("---------OnServerStarted--------------:", string(reqStr))
	fmt.Println("[--Debug--] HooksApiServerStartedReq struct:", req)
	//当该配置更新到数据库中去。
	//先要判断MediaserverID是不是有存在
	vidTmp, err := relationDB.NewVidmgrInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{req.GeneralMediaServerId},
	})

	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=VidmgrInfoRead err=%v", utils.FuncName(), er)
	}
	if vidTmp.VidmgrID != "" {
		//查到数据，更新服务状态为在线
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
