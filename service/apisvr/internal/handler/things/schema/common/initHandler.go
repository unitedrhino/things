package common

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"github.com/i-Things/things/service/apisvr/internal/logic/things/schema/common"
	"github.com/i-Things/things/service/apisvr/internal/svc"
)

func InitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewInitLogic(r.Context(), svcCtx)
		err := l.Init()
		result.Http(w, r, nil, err)
	}
}
