package middlewares

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/unsafetools"
)

func TestRecoverPanic(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	*unsafetools.FieldByName(ctx, `s`).(**fasthttp.Server) = &fasthttp.Server{
		Logger: log.New(ioutil.Discard, ``, 0),
	}

	RecoverPanic(func(ctx *fasthttp.RequestCtx) {
		panic(`test`)
	})(ctx)

	assert.Equal(t, 500, ctx.Response.StatusCode())
}
