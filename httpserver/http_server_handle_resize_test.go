package httpserver

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/unsafetools"
)

func TestHTTPServerHandleResize(t *testing.T) {
	srv := testDummyServer(t)

	checkRequest := func(expectedStatusCode int, uri string) {
		ctx := &fasthttp.RequestCtx{}
		*unsafetools.FieldByName(ctx, `s`).(**fasthttp.Server) = srv.httpBackend
		ctx.Request.SetRequestURI(uri)
		srv.httpRouter.Handler(ctx)

		assert.Equal(t, expectedStatusCode, ctx.Response.StatusCode(), string(ctx.Response.Body()))
	}

	checkRequest(200, `/resize?width=100&height=100&url=test_picture`)
	checkRequest(400, `/resize?height=100&url=test_picture`)
	checkRequest(400, `/resize?width=100&url=test_picture`)
	checkRequest(400, `/resize?width=100&height=100`)
	checkRequest(400, `/resize?width=0&height=100&url=test_picture`)
	checkRequest(400, `/resize?width=-100&height=100&url=test_picture`)
	checkRequest(400, `/resize?width=100000&height=100&url=test_picture`)
	checkRequest(400, `/resize?width=100&height=0&url=test_picture`)
	checkRequest(400, `/resize?width=100&height=-100&url=test_picture`)
	checkRequest(502, `/resize?width=100&height=100&url=htttp://invalid_url/`)
}

func BenchmarkHTTPServerResize(b *testing.B) {
	srv := testDummyServer(nil)

	ctxPool := sync.Pool{
		New: func() interface{} {
			ctx := &fasthttp.RequestCtx{}
			*unsafetools.FieldByName(ctx, `s`).(**fasthttp.Server) = srv.httpBackend
			ctx.Request.SetHost(`localhost`)
			ctx.Request.SetRequestURI(`/resize?width=100&height=100&url=test_picture`)
			return ctx
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := ctxPool.Get().(*fasthttp.RequestCtx)
			srv.httpRouter.Handler(ctx)
			if ctx.Response.StatusCode() != 200 {
				panic(fmt.Sprintf("unexpected status code: %v.\nHeaders:\n%v\nBody:\n%v",
					ctx.Response.StatusCode(),
					ctx.Response.Header.String(),
					string(ctx.Response.Body()),
				))
			}
			ctx.ResetBody()
			ctxPool.Put(ctx)
		}
	})
}
