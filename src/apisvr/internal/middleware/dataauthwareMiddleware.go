package middleware

import (
	"github.com/i-Things/things/src/apisvr/internal/config"
	"net/http"
)

type DataAuthWareMiddleware struct {
	cfg config.Config
}

func NewDataAuthWareMiddleware(cfg config.Config) *DataAuthWareMiddleware {
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
		next(w, r)
	}
}
