package httpgin

import (
	"context"
	"fmt"
	"net/http"
	httpserver "study-corner-common/pkg/http-server"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    httpserver.Config
	Init      httpserver.RouterInitializer
}

// ginContext adapts gin.Context to httpserver.Context.
type ginContext struct {
	c *gin.Context
}

func (g *ginContext) Request() *http.Request                { return g.c.Request }
func (g *ginContext) Param(name string) string              { return g.c.Param(name) }
func (g *ginContext) Query(name string) string              { return g.c.Query(name) }
func (g *ginContext) BindJSON(target any) error             { return g.c.ShouldBindJSON(target) }
func (g *ginContext) JSON(code int, body any)               { g.c.JSON(code, body) }
func (g *ginContext) Status(code int)                       { g.c.Status(code) }

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

// Module wires Gin as an implementation of httpserver.Router + HTTP server
func Module() fx.Option {
    return fx.Options(
        fx.Provide(
            func() *gin.Engine {
                r := gin.New()
                r.Use(gin.Recovery())
                return r
            },
            func(engine *gin.Engine) httpserver.Router {
                return &ginRouter{engine: engine}
            },
        ),
        fx.Invoke(func(p Params, engine *gin.Engine, router httpserver.Router) {
            // let the service register its routes
            if p.Init != nil {
                p.Init(router)
            }

            srv := &http.Server{
                Addr:    fmt.Sprintf(":%d", p.Config.Port),
                Handler: engine,
            }

            p.Lifecycle.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    go srv.ListenAndServe()
                    return nil
                },
                OnStop: func(ctx context.Context) error {
                    return srv.Shutdown(ctx)
                },
            })
        }),
    )
}