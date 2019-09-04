package httpserver

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/unsafetools"
)

func TestHTTPServerHandleResize(t *testing.T) {
	logger := log.New(os.Stderr, ``, 0)

	srv, err := NewHTTPServer(&dummyImageProcessor{}, logger, logger)
	assert.NoError(t, err)

	checkRequest := func(expectedStatusCode int, uri string) {
		ctx := &fasthttp.RequestCtx{}
		*(**fasthttp.Server)(unsafetools.FieldByName(ctx, `s`, (**fasthttp.Server)(nil))) = srv.httpBackend
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
