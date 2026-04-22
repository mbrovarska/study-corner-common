package httpgin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	httpserver "study-corner-common/pkg/http-server"
	log "study-corner-common/pkg/logger"
)

type Params struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Config         httpserver.Config
	Logger         log.Logger
	RegisterRoutes httpserver.RouteRegistrar `optional:"true"`
}

type ginContext struct {
	c *gin.Context
}

func (g *ginContext) Request() *http.Request    { return g.c.Request }
func (g *ginContext) Param(name string) string  { return g.c.Param(name) }
func (g *ginContext) Query(name string) string  { return g.c.Query(name) }
func (g *ginContext) BindJSON(target any) error { return g.c.ShouldBindJSON(target) }
func (g *ginContext) JSON(code int, body any)   { g.c.JSON(code, body) }
func (g *ginContext) Status(code int)           { g.c.Status(code) }

type ginRouter struct {
	engine *gin.Engine
}

func (r *ginRouter) wrap(h httpserver.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&ginContext{c: c})
	}
}

func (r *ginRouter) GET(path string, h httpserver.HandlerFunc)    { r.engine.GET(path, r.wrap(h)) }
func (r *ginRouter) POST(path string, h httpserver.HandlerFunc)   { r.engine.POST(path, r.wrap(h)) }
func (r *ginRouter) PUT(path string, h httpserver.HandlerFunc)    { r.engine.PUT(path, r.wrap(h)) }
func (r *ginRouter) DELETE(path string, h httpserver.HandlerFunc) { r.engine.DELETE(path, r.wrap(h)) }

func NewEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	return r
}

func NewRouter(engine *gin.Engine) httpserver.Router {
	return &ginRouter{engine: engine}
}

func RegisterServer(p Params, engine *gin.Engine) {
	if p.RegisterRoutes != nil {
		p.RegisterRoutes(NewRouter(engine))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Config.Port),
		Handler: engine,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					p.Logger.Error("http server stopped unexpectedly")
				}
			}()
			p.Logger.Info("http server started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

var Module = fx.Module(
	"httpgin",
	fx.Provide(
		NewEngine,
		NewRouter,
	),
	fx.Invoke(RegisterServer),
)