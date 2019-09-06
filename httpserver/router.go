package httpserver

import (
	"time"

	"github.com/buaazp/fasthttprouter"

	"github.com/xaionaro-go/picapi/httpserver/middlewares"
)

func (srv *HTTPServer) newRouter(
	cacheDuration time.Duration,
	cacheMaxEntries uint64,
	cacheMaxEntrySize uint64,
) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET(`/resize`, middlewares.Cache(
		cacheMaxEntries, cacheMaxEntrySize, cacheDuration,
		srv.handleResize,
	))

	return router
}
