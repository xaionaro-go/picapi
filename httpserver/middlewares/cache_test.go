package middlewares

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func getTestHandler() fasthttp.RequestHandler {
	router := fasthttprouter.New()

	var count uint32
	handler := func(ctx *fasthttp.RequestCtx) {
		result := atomic.AddUint32(&count, 1)

		ctx.Response.Header.Set(`Content-Type`, `application/test-`+fmt.Sprint(result))
		ctx.WriteString(fmt.Sprint(result))

		// The line below is for test "TestCache_concurrentRequests".
		//
		// So the execution of the first request will be paused until the second
		// request will reach this line.
		runtime.Gosched()
	}

	router.GET(`/cached`, Cache(1<<10, 1<<10, 365*24*time.Hour /* close enough to an infinite */, handler))
	router.GET(`/negative_expire_time`, Cache(1<<10, 1<<10, -time.Second, handler))

	return router.Handler
}

func TestCache_cachedResponse(t *testing.T) {
	handler := getTestHandler()

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI(`/cached`)

	handler(ctx)
	assert.Equal(t, `max-age=31536000`, string(ctx.Response.Header.Peek(`Cache-Control`)))
	assert.Empty(t, ctx.Response.Header.Peek(`X-Cached-Response`))
	assert.Equal(t, 200, ctx.Response.StatusCode())
	firstResponseBody := string(ctx.Response.Body())
	firstResponseContentType := string(ctx.Response.Header.Peek(`Content-Type`))

	ctx.ResetBody()
	handler(ctx)
	assert.Equal(t, `max-age=31535999`, string(ctx.Response.Header.Peek(`Cache-Control`)))
	assert.Equal(t, `true`, string(ctx.Response.Header.Peek(`X-Cached-Response`)))
	assert.Equal(t, 200, ctx.Response.StatusCode())
	secondResponseBody := string(ctx.Response.Body())
	secondResponseContentType := string(ctx.Response.Header.Peek(`Content-Type`))

	assert.Equal(t, firstResponseBody, secondResponseBody)
	assert.Equal(t, firstResponseContentType, secondResponseContentType)
}

func TestCache_expiredCacheResponse(t *testing.T) {
	handler := getTestHandler()

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI(`/negative_expire_time`)

	handler(ctx)
	assert.Equal(t, 200, ctx.Response.StatusCode())
	firstResponseBody := string(ctx.Response.Body())

	ctx.ResetBody()
	handler(ctx)
	assert.Equal(t, 200, ctx.Response.StatusCode())
	secondResponseBody := string(ctx.Response.Body())

	assert.NotEqual(t, firstResponseBody, secondResponseBody)
}

func TestCache_concurrentRequests(t *testing.T) {
	handler := getTestHandler()

	ctx0 := &fasthttp.RequestCtx{}
	ctx0.Request.SetRequestURI(`/cached`)

	ctx1 := &fasthttp.RequestCtx{}
	ctx1.Request.SetRequestURI(`/cached`)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		handler(ctx0)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		handler(ctx1)
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, string(ctx0.Response.Body()), string(ctx1.Response.Body()))
}
