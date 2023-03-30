package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.AddMsg(err.Error()))
			return
		}

		strIP, _ := utils.GetIP(r)
		c := context.WithValue(r.Context(), userHeader.UserUid, &userHeader.UserCtx{
			IP: strIP,
			Os: r.Header.Get("User-Agent"),
		})
		l := user.NewLoginLogic(c, svcCtx)
		resp, err := l.Login(&req)
		result.Http(w, r, resp, err)
	}
}
