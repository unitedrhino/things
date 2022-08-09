package middleware

import (
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type RecordMiddleware struct {
}

func NewRecordMiddleware() *RecordMiddleware {
	return &RecordMiddleware{}
}

func (m *RecordMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		guid := r.Header.Get(userHeader.GUID)
		logx.WithContext(r.Context()).Infof("guid=%s", guid)
	}
}
