package middleware

import (
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type RecordMiddleware struct {
}

func NewRecordMiddleware() *RecordMiddleware {
	return &RecordMiddleware{}
}

func (m *RecordMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spanCtx := trace.SpanContextFromContext(r.Context())
		w.Header().Add(userHeader.GUID, spanCtx.TraceID().String())
		next(w, r)
		guid := r.Header.Get(userHeader.GUID)
		logx.WithContext(r.Context()).Infof("guid=%s", guid)
	}
}
