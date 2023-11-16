package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/client/vidmgrinfomanage"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	vidTmp, err := l.svcCtx.VidmgrM.VidmgrInfoRead(l.ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: req.GeneralMediaServerId,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=VidmgrInfoRead err=%v", utils.FuncName(), er)
	}
	if vidTmp.VidmgrID != "" {
		//根据MediaserverID查询数据，如果数据在的话做更新，不在的话就插入这条数据
		if req.GeneralMediaServerId != "" {
			vidmgrConfig, err1 := l.svcCtx.VidmgrC.VidmgrConfigRead(l.ctx, &vidmgrinfomanage.VidmgrConfigReadReq{
				MediaServerId: req.GeneralMediaServerId,
			})
			if err1 != nil {
				l.Errorf("%s.rpc.VidmgrConfigmange req=VidmgrConfigRead err=%v", utils.FuncName(), err1)
			}
			if vidmgrConfig.GeneralMediaServerId != "" {
				//we inster data.
				//_, err = l.svcCtx.VidmgrC.VidmgrConfigCreate(l.ctx, info.ToVidmgrConfigRpc(req))
				if err != nil {
					er := errors.Fmt(err)
					l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
				}
			} else {
				//we update data
			}
		}
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
