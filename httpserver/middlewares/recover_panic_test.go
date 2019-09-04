package middlewares

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRecoverPanic(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}

	RecoverPanic(func(ctx *fasthttp.RequestCtx) {
		panic(`test`)
	})(ctx)

	assert.Equal(t, 500, ctx.Response.StatusCode())
}
