package common

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config() (resp *types.ConfigResp, err error) {
	rsp, err := l.svcCtx.Common.Config(l.ctx, &sys.Response{})
	if err != nil {
		err = errors.Fmt(err)
		l.Errorf("%s.rpc.SysConfig err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &types.ConfigResp{Map: types.Map{Mode: rsp.Map.Mode, AccessKey: rsp.Map.AccessKey},
		Oss: types.Oss{Host: l.svcCtx.Config.OssConf.CustomHost}}, nil
}
