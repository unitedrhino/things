package middleware

import (
	"bytes"
	"context"
	"gitee.com/i-Things/core/shared/caches"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/config"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
)

type DataAuthWareMiddleware struct {
	cfg config.Config
}

func NewDataAuthWareMiddleware(cfg config.Config) *DataAuthWareMiddleware {
	caches.InitStore(cfg.CacheRedis)
	return &DataAuthWareMiddleware{cfg: cfg}
}

type DataAuthParam struct {
	ProjectID  string   `json:"projectID,string,optional"` //项目id
	ProjectIDs []string `json:"projectIDs,optional"`       //项目ids
	AreaID     string   `json:"areaID,string,optional"`    //项目区域id
	AreaIDs    []string `json:"areaIDs,optional"`          //项目区域ids
}

func (m *DataAuthWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userCtx := ctxs.GetUserCtxOrNil(ctx)
		//没有用户态or拥有所有数据权限，则不校验数据权限
		if userCtx == nil || userCtx.IsAllData == true || true {
			next(w, r)
			return
		}

		//读出Body 暂存
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//写回Body 给 httpx.Parse 读
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		//httpx.Parse 读
		var param DataAuthParam
		//zeromicro/go-zero@v1.5.0/core/mapping/unmarshaler.go:84 要求接收者是结构体
		err = httpx.Parse(r, &param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//其他人 要校验数据权限
		{ //检查项目权限
			dataType := def.AuthDataTypeProject
			reqIDs := param.ProjectIDs
			//全局项目ID
			mdProjectID := ctxs.GetUserCtx(ctx).ProjectID
			if mdProjectID != 0 {
				reqIDs = append(reqIDs, utils.ToString(mdProjectID))
			}
			//参数项目ID
			if param.ProjectID != "" {
				reqIDs = append(reqIDs, param.ProjectID)
			}
			//校验项目权限
			if code := m.check(ctx, dataType, reqIDs); code != http.StatusOK {
				err := errors.Permissions.AddMsg(def.AuthDataTypeFieldTextMap[dataType] + "不足")
				http.Error(w, err.Error(), code)
				return
			}
		}
		{ //检查区域权限
			dataType := def.AuthDataTypeArea
			reqIDs := param.AreaIDs
			//参数项目ID
			if param.AreaID != "" {
				reqIDs = append(reqIDs, param.AreaID)
			}
			//校验区域权限
			if code := m.check(ctx, dataType, reqIDs); code != http.StatusOK {
				err := errors.Permissions.AddMsg(def.AuthDataTypeFieldTextMap[dataType] + "不足")
				http.Error(w, err.Error(), code)
				return
			}
		}

		//写回Body 恢复原状
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next(w, r)
	}
}

func (m *DataAuthWareMiddleware) check(ctx context.Context, dataType def.AuthDataType, reqIDs []string) int {
	//if len(reqIDs) > 0 {
	//	authIDs, err := caches.GetUserDataAuth(ctx, dataType)
	//	diffIDs := utils.SliceLeftDiff(reqIDs, authIDs)
	//	if err == redis.Nil || (err == nil && len(diffIDs) > 0) { //没有数据权限
	//		logx.WithContext(ctx).Errorf("%s.没有数据权限 dataType=%#v, reqIDs=%#v, authIDs=%#v, diffIDs=%#v", utils.FuncName(), dataType, reqIDs, authIDs, diffIDs)
	//		return http.StatusUnauthorized
	//	} else if err != nil { //校验数据权限异常
	//		logx.WithContext(ctx).Errorf("%s.校验数据权限异常 dataType=%#v, reqIDs=%#v, authIDs=%#v, error=%#v", utils.FuncName(), dataType, reqIDs, authIDs, err)
	//		return http.StatusInternalServerError
	//	}
	//}
	return http.StatusOK
}
