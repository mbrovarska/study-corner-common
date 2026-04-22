package httpserver

import "net/http"

type HandlerFunc func(Context)

type Context interface {
	Request() *http.Request
	Param(name string) string
	Query(name string) string
	BindJSON(target any) error
	JSON(code int, body any)
	Status(code int)
}

type Router interface {
	GET(path string, h HandlerFunc)
	POST(path string, h HandlerFunc)
	PUT(path string, h HandlerFunc)
	DELETE(path string, h HandlerFunc)
}

type RouteRegistrar func(r Router)

type Config struct {
	Port int
}