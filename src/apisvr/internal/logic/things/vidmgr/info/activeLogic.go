package info

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActiveLogic {
	return &ActiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActiveLogic) Active(req *types.VidmgrInfoActiveReq) error {
	// todo: add your logic here and delete this line
	//read VidmgrInfo table and update table
	_, err := l.svcCtx.VidmgrM.VidmgrInfoActive(l.ctx, &vid.VidmgrInfoActiveReq{
		VidmgrID: req.VidmgrID,
	})
	if err != nil {
		//err
		fmt.Sprintln("[***testActive**]", utils.FuncName(), err)
		l.Errorf("active falied:", utils.FuncName(), err)
		return errors.MediaActiveError.AddMsg(err.Error())
	}
	return nil
}
