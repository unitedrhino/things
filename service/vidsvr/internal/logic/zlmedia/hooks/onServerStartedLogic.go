package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/vidsvr/internal/common"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
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

	if vidTmp != nil {
		//流服务注册过，就将配置信息记录更新到数据库
		//1 查表config配置表
		confRepo := relationDB.NewVidmgrConfigRepo(l.ctx)
		confRepo.FindOneByFilter(l.ctx, relationDB.VidmgrConfigFilter{
			Secret: vidTmp.VidmgrSecret,
		})
		if err != nil {
			l.Errorf("%s.Can find vidmgr config err=%v", utils.FuncName(), utils.Fmt(err))
			confRepo.Insert(l.ctx, common.ToVidmgrConfigDB2(req))
		} else {
			//update
			confRepo.Update(l.ctx, common.ToVidmgrConfigDB2(req))
		}
	} else {
		l.Errorf("%s.rpc.ManageVidmgr req=VidmgrInfoRead err=%v", utils.FuncName(), err)
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
