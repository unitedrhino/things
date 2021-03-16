package middleware

import "net/http"

type UsercheckMiddleware struct {
}

func NewUsercheckMiddleware() *UsercheckMiddleware {
	return &UsercheckMiddleware{}
}

func (m *UsercheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
