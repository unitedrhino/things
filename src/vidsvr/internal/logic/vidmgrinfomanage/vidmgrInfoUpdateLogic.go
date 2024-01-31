package vidmgrinfomanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoUpdateLogic {
	return &VidmgrInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

// 更新服务
func (l *VidmgrInfoUpdateLogic) VidmgrInfoUpdate(in *vid.VidmgrInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	fmt.Printf("Vidsvr VidmgrInfoUpdate \n")
	po, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.VidmgrFilter{VidmgrIDs: []string{in.VidmgrID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find Product_id id:" + cast.ToString(in.VidmgrID))
		}
		return nil, err
	}
	err = common.UpdatVidmgrInfoDB(po, in)
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
