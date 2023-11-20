package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnStreamChangedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnStreamChangedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamChangedLogic {
	return &OnStreamChangedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnStreamChangedLogic) OnStreamChanged(req *types.HooksApiStreamChangedRep) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnStreamChanged--------------:", string(reqStr))
	//获取流注册时分支
	if req.Regist {
		//需要先判断该流服务是否有注册过，未注册过
		//判断流路径是一样

	} else { //注销时分支
		//检查数据库是否存在，
		//存在即删除该流
	}

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
