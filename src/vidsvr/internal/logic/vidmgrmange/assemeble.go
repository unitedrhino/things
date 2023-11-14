package vidmgrmangelogic

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
)

func ToVidmgrInfo(ctx context.Context, pi *relationDB.VidmgrInfo, svcCtx *svc.ServiceContext) *vid.VidmgrInfo {

	if pi.VidmgrType == def.Unknown {
		pi.VidmgrType = def.VidmgrTypeZLMedia //当前默认仅支持zlmediakit
	}
	dpi := &vid.VidmgrInfo{
		VidmgrID:     pi.VidmgrID,   //服务id
		VidmgrName:   pi.VidmgrName, //服务名
		VidmgrIpV4:   utils.InetNtoA(pi.VidmgrIpV4),
		VidmgrPort:   pi.VidmgrPort,
		VidmgrType:   pi.VidmgrType,                         //流服务器类型:1:zlmediakit,2:srs,3:monibuca
		VidmgrStatus: pi.VidmgrStatus,                       //服务状态: 1：离线 2：在线  0：未激活
		VidmgrSecret: pi.VidmgrSecret,                       //流服务器注秘钥 只读
		Desc:         &wrappers.StringValue{Value: pi.Desc}, //描述
		CreatedTime:  pi.CreatedTime.Unix(),                 //创建时间
		Tags:         pi.Tags,                               //产品tags
	}

	return dpi
}
