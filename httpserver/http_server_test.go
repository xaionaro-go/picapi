package httpserver

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/picapi/config"
	"github.com/xaionaro-go/picapi/httpserver/helpers"
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
	"github.com/xaionaro-go/unsafetools"
)

type dummyImageProcessor struct{}

func (proc *dummyImageProcessor) Resize(
	in io.Reader,
	out io.Writer,
	toWidth, toHeight uint,
) (imageFormat imageprocessorcommon.ImageFormat, err error) {
	imageFormat = imageprocessorcommon.ImageFormatJPEG
	return
}

func testDummyServer(t *testing.T) *HTTPServer {
	logger := log.New(ioutil.Discard, ``, 0)

	var cfg config.Config
	err := envconfig.Process("picapi", &cfg)
	if t != nil {
		assert.NoError(t, err)
	}

	srv, err := NewHTTPServer(
		&dummyImageProcessor{},
		logger,
		logger,
		cfg.CacheDuration,
		cfg.CacheMaxEntries,
		cfg.CacheMaxEntrySize,
	)
	if t != nil {
		assert.NoError(t, err)
	}

	return srv
}

func TestHTTPServerStartStop(t *testing.T) {
	srv := testDummyServer(t)

	err := srv.Start(context.Background(), "")
	assert.NoError(t, err)

	err = srv.Stop()
	assert.NoError(t, err)

	err = srv.Start(context.Background(), "")
	assert.NoError(t, err)

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI(`/resize?width=100&height=100&url=test_picture`)
	*unsafetools.FieldByName(ctx, `s`).(**fasthttp.Server) = srv.httpBackend

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err = srv.Stop()
		assert.NoError(t, err)
		wg.Done()
	}()

	srv.httpBackend.Handler(ctx)

	wg.Wait()

	err = srv.Start(context.Background(), "")
	assert.NoError(t, err)

	func() {
		defer helpers.RecoverPanic(&err)
		err = srv.Start(context.Background(), "")
	}()
	assert.Error(t, err)
}

func TestHTTPServerStartStopWithNilLogger(t *testing.T) {
	srv := testDummyServer(t)

	err := srv.Start(context.Background(), "")
	assert.NoError(t, err)

	err = srv.Stop()
	assert.NoError(t, err)
}

func TestHTTPServerStartCancel(t *testing.T) {
	srv := testDummyServer(t)

	ctx, cancelFunc := context.WithCancel(context.Background())

	srv.Start(ctx, ``)

	reqCtx := &fasthttp.RequestCtx{}
	reqCtx.Request.SetRequestURI(`/resize?width=100&height=100&url=test_picture`)
	*unsafetools.FieldByName(reqCtx, `s`).(**fasthttp.Server) = srv.httpBackend

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		srv.httpBackend.Handler(reqCtx)
		wg.Done()
	}()

	cancelFunc()

	wg.Wait()

	srv.Wait()

	srv.Start(context.Background(), ``)
	srv.Stop()
}
