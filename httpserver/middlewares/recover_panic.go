package middlewares

import (
	"log"
	"os"

	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/picapi/httpserver/helpers"
)

// RecoverPanic recovers panics in the `handler` and
// returns `500 Internal Server Error` if a panic occurs.
func RecoverPanic(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var panicError error

		func() {
			defer helpers.RecoverPanic(&panicError)
			handler(ctx)
		}()

		if panicError != nil {
			ctx.SetStatusCode(500)
			ctx.Response.ResetBody()

			var logger Printfer
			func() {
				defer helpers.RecoverPanic(nil)
				logger = ctx.Logger()
			}()
			if logger == nil {
				logger = log.New(os.Stderr, ``, 0)
			}
			logger.Printf("[panic] %v\n", panicError)
		}
	}
}
