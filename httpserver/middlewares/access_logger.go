package middlewares

import (
	"github.com/valyala/fasthttp"
)

// Printfer is a simple logger interface for `NewHTTPServer`
type Printfer interface {
	// Printf is a pretty-standart `Printf` (`fmt` is the format and `args` are the arguments for the format).
	//
	// For example you can use `log.New` as `Printfer`.
	Printf(fmt string, args ...interface{})
}

// AccessLogger wraps the `handler` to log requests via `logger` (like `access.log`)
func AccessLogger(logger Printfer, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		handler(ctx)
		logger.Printf("[%v] %v %v\n", ctx.Response.StatusCode(), string(ctx.Method()), string(ctx.RequestURI()))
	}
}
