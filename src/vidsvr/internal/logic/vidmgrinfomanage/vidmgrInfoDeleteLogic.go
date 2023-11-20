package vidmgrinfomanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoDeleteLogic {
	return &VidmgrInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

// 删除服务
func (l *VidmgrInfoDeleteLogic) VidmgrInfoDelete(in *vid.VidmgrInfoDeleteReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	fmt.Printf("Vidsvr VidmgrInfoDelete \n")
	err := l.PiDB.DeleteByFilter(l.ctx, relationDB.VidmgrFilter{VidmgrIDs: []string{in.VidmgrtID}})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	//更新删除事件
	return &vid.Response{}, nil
}

func (l *VidmgrInfoDeleteLogic) DropVidmgr(in *vid.VidmgrInfoDeleteReq) error {
	//需要删除与该流媒体服务器绑定的其它服务
	return nil
}

func (l *VidmgrInfoDeleteLogic) Check(in *vid.VidmgrInfoDeleteReq) error {
	//需要判断该流媒体服务器下是否有绑定过流设备  如果有流设备，则不能删除
	return nil
}
