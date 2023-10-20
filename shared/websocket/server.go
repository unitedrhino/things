package websocket

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/router"
	"log"
	"net/http"
	"path"
)

type (
	// RunOption defines the method to customize a Server.
	RunOption func(*Server)

	// A Server is a http server.
	Server struct {
		ngin   *engine
		router httpx.Router
	}
)

// MustNewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
// The process will exit if error occurs.
func MustNewServer(c rest.RestConf, opts ...RunOption) *Server {
	server, err := NewServer(c, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

// NewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
func NewServer(c rest.RestConf, opts ...RunOption) (*Server, error) {
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	server := &Server{
		ngin:   newEngine(c),
		router: router.NewRouter(),
	}

	opts = append([]RunOption{WithNotFoundHandler(nil)}, opts...)
	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

// WithMiddlewares adds given middlewares to given routes.
func WithMiddlewares(ms []rest.Middleware, rs ...rest.Route) []rest.Route {
	for i := len(ms) - 1; i >= 0; i-- {
		rs = WithMiddleware(ms[i], rs...)
	}
	return rs
}

// WithMiddleware adds given middleware to given route.
func WithMiddleware(middleware rest.Middleware, rs ...rest.Route) []rest.Route {
	routes := make([]rest.Route, len(rs))

	for i := range rs {
		route := rs[i]
		routes[i] = rest.Route{
			Method:  route.Method,
			Path:    route.Path,
			Handler: middleware(route.Handler),
		}
	}

	return routes
}

// WithNotFoundHandler returns a RunOption with not found handler set to given handler.
func WithNotFoundHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		notFoundHandler := server.ngin.notFoundHandler(handler)
		server.router.SetNotFoundHandler(notFoundHandler)
	}
}

// WithNotAllowedHandler returns a RunOption with not allowed handler set to given handler.
func WithNotAllowedHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(handler)
	}
}

// WithPrefix adds group as a prefix to the route paths.
func WithPrefix(group string) RouteOption {
	return func(r *featuredRoutes) {
		var routes []rest.Route
		for _, rt := range r.routes {
			p := path.Join(group, rt.Path)
			routes = append(routes, rest.Route{
				Method:  rt.Method,
				Path:    p,
				Handler: rt.Handler,
			})
		}
		r.routes = routes
	}
}

// WithPriority returns a RunOption with priority.
func WithPriority() RouteOption {
	return func(r *featuredRoutes) {
		r.priority = true
	}
}

// WithRouter returns a RunOption that make server run with given router.
func WithRouter(router httpx.Router) RunOption {
	return func(server *Server) {
		server.router = router
	}
}

// AddRoutes add given routes into the Server.
func (s *Server) AddRoutes(rs []rest.Route, opts ...RouteOption) {
	r := featuredRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	s.ngin.addRoutes(r)
}

// AddRoute adds given route into the Server.
func (s *Server) AddRoute(r rest.Route, opts ...RouteOption) {
	s.AddRoutes([]rest.Route{r}, opts...)
}

// PrintRoutes prints the added routes to stdout.
func (s *Server) PrintRoutes() {
	s.ngin.print()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ngin.bindRoutes(s.router)
	s.router.ServeHTTP(w, r)
}
