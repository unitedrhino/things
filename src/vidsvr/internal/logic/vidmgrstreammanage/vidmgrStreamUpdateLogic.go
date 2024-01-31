package vidmgrstreammanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamUpdateLogic {
	return &VidmgrStreamUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 流更新
func (l *VidmgrStreamUpdateLogic) VidmgrStreamUpdate(in *vid.VidmgrStream) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	po, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.VidmgrStreamFilter{
		StreamIDs: []int64{in.StreamID},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find stream id:" + cast.ToString(in.StreamID))
		}
		return nil, err
	}
	err = common.UpdateVidmgrStreamDB(po, in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
		//if errors.Cmp(err, errors.Duplicate) {
		//	return nil, errors.Duplicate.WithMsgf("产品名称重复:%s", in.VidmgrName)
		//}
		return nil, err
	}
	return &vid.Response{}, nil
}
