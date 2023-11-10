package vidmgrmangelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"github.com/spf13/cast"
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
		PiDB:   relationDB.NewVidmgrtInfoRepo(ctx),
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
	err = l.setPoByPb(po, in)
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

func (l *VidmgrInfoUpdateLogic) setPoByPb(old *relationDB.VidmgrInfo, data *vid.VidmgrInfo) error {
	if data.VidmgrName != "" {
		old.VidmgrName = data.VidmgrName
	}
	if data.VidmgrIpV4 != "" {
		old.VidmgrIpV4 = utils.InetAtoN(data.VidmgrIpV4)
	}
	if data.VidmgrPort != 0 {
		old.VidmgrPort = data.VidmgrPort
	}
	if data.VidmgrType != 0 {
		old.VidmgrType = data.VidmgrType
	}
	if data.VidmgrSecret != "" {
		old.VidmgrSecret = data.VidmgrSecret
	}
	if data.Desc != nil {
		old.Desc = data.Desc.GetValue()
	}
	if data.Tags != nil {
		old.Tags = data.Tags
	}
	return nil
}
