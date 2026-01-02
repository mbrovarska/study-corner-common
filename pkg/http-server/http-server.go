package httpserver

import "net/http"


type Config struct {
	Port int `yaml:"port"`
}

//handlers see context
type Context interface {
	Request() *http.Request
	Param(name string) string
	Query(name string) string
	BindJSON(target any) error
	JSON(statusCode int, body any)
	Status(StatusCode int)
}

//generic handler
type HandlerFunc func(Context)

//minimal API to register routes
type Router interface {
	GET(path string, h HandlerFunc)
	POST(path string, h HandlerFunc)
	PUT(path string, h HandlerFunc)
	DELETE(path string, h HandlerFunc)
}

//plug routes in
type RouterInitializer func(r Router)