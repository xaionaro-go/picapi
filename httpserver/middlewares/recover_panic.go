package middlewares

import (
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/picapi/httpserver/helpers"
)

// RecoverPanic recovers panics in the `handler` and
// returns `500 Internal Server Error` if a panic occurs.
func RecoverPanic(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			var panicError error
			helpers.RecoverPanic(&panicError)

			if panicError != nil {
				ctx.SetStatusCode(500)
				ctx.Response.ResetBody()

				ctx.Logger().Printf("[panic] %v\n", panicError)
			}
		}()

		handler(ctx)
	}
}
