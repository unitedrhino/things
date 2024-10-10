package common

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/schema/common"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
)

func InitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewInitLogic(r.Context(), svcCtx)
		err := l.Init()
		result.Http(w, r, nil, err)
	}
}
