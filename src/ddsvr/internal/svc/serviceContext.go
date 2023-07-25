package svc

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/domain/custom"
	"github.com/i-Things/things/src/ddsvr/internal/repo/cache"

	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubInner"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config   config.Config
	PubDev   pubDev.PubDev
	PubInner pubInner.PubInner
	ProductM productmanage.ProductManage
	Script   custom.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		productM productmanage.ProductManage
	)
	dl, err := pubDev.NewPubDev(c.DevLink)
	if err != nil {
		logx.Error("NewDevClient err", err)
		os.Exit(-1)
	}

	il, err := pubInner.NewPubInner(c.Event)
	if err != nil {
		logx.Error("NewInnerDevPub err", err)
		os.Exit(-1)
	}
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
	}
	scriptCache := cache.NewScriptRepo(func(ctx context.Context, productID string) (info *custom.Info, err error) {
		ret, err := productM.ProductCustomRead(ctx, &dm.ProductCustomReadReq{ProductID: productID})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) { //如果是没找到
				return nil, nil
			}
			return nil, err
		}
		return &custom.Info{
			ProductID:       ret.ProductID,
			TransformScript: utils.ToNullString(ret.TransformScript),
			ScriptLang:      ret.ScriptLang,
		}, nil
	})
	return &ServiceContext{
		Config:   c,
		PubDev:   dl,
		PubInner: il,
		ProductM: productM,
		Script:   scriptCache,
	}
}
