package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/valyala/fasthttp"
)

const (
	// CacheRecentRatio is the fraction of recently added entries in the cache
	//
	// See more at https://godoc.org/github.com/hashicorp/golang-lru#pkg-constants
	CacheRecentRatio = 0.15

	// CacheGhostRatio is the fraction of ghost entries in the cache
	//
	// See more at https://godoc.org/github.com/hashicorp/golang-lru#pkg-constants
	CacheGhostRatio = 0.5
)

// Cache wraps the `handler` returns a new handler which will cache responses
// (and return responses from the cache if there's an actual record) for a time
// interval `maxDuration`.
//
// To limit memory usage there're parameters `maxEntries` and `maxEntrySize`.
// Total consumed memory (by this middleware) could be estimated as `maxEntries * maxEntrySize`.
//
// Also this middleware controls HTTP header `Cache-Control` to ask a browser to cache
// the response on it's side, too.
//
// If there were received two similar (with the same URI) concurrent requests then
// it will wait for the execution of one of them then will just return the cached
// value to the second one.
func Cache(
	maxEntries uint64,
	maxEntrySize uint64,
	maxDuration time.Duration,
	handler fasthttp.RequestHandler,
) fasthttp.RequestHandler {
	if maxEntries <= 0 {
		return handler
	}
	if maxEntrySize <= 0 {
		return handler
	}
	if maxDuration <= 0 {
		return handler
	}

	return newCacheFilter(
		maxEntries,
		maxEntrySize,
		maxDuration,
		handler,
	).Handler
}

type cacheFilter struct {
	NextHandler        fasthttp.RequestHandler
	HeaderCacheControl string
	MaxEntrySize       uint64
	MaxDuration        time.Duration
	Storage            *lru.TwoQueueCache
	URILocker          *cacheURILocker
}

func newCacheFilter(
	maxEntries uint64,
	maxEntrySize uint64,
	maxDuration time.Duration,
	handler fasthttp.RequestHandler,
) *cacheFilter {
	if maxDuration.Seconds() <= 1 {
		// TODO: report error
	}

	cacheStorage, err := lru.New2QParams(
		int(maxEntries),
		CacheRecentRatio,
		CacheGhostRatio,
	)
	if err != nil {
		panic(err)
	}

	return &cacheFilter{
		NextHandler:        handler,
		MaxEntrySize:       maxEntrySize,
		MaxDuration:        maxDuration,
		HeaderCacheControl: fmt.Sprintf(`max-age=%d`, uint64(maxDuration.Seconds())),
		Storage:            cacheStorage,
		URILocker:          newCacheURILocker(),
	}
}

func (c *cacheFilter) Handler(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	switch method {
	case `GET`, `OPTIONS`:
	default:
		panic(`middleware "Cache" doesn't support method ` + method)
	}

	uri := ctx.Request.URI()
	uri.QueryArgs().Sort(bytes.Compare)
	uriString := string(uri.RequestURI())

	// Wait until requests to the same URI will finish (to get a cached value
	// instead of real processing).
	c.URILocker.Lock(uriString)
	defer c.URILocker.Unlock(uriString)

	// Try from cache
	cacheItem, filledFromCache := c.fillResponseFromCache(ctx, uriString)
	if filledFromCache {
		// OK, the reply is already placed into `ctx`.
		//
		// But we also want to save it to the client-side cache:
		secondsLeft := cacheItem.ExpireTS.Sub(time.Now()).Seconds()
		if secondsLeft > 0 {
			ctx.Response.Header.Set(`Cache-Control`, fmt.Sprintf(`max-age=%d`, uint64(secondsLeft)))
		}
		return
	}

	// No cache, just process it
	c.NextHandler(ctx)

	// Cache only if the status code is "200" (OK) and "400" (BadRequest)
	status := ctx.Response.StatusCode()
	switch status {
	case http.StatusOK, http.StatusBadRequest:
	default:
		return
	}

	// Save to the server-side cache
	c.saveToCacheFromResponse(cacheItem, ctx, uriString)

	// Save to the client-side cache
	ctx.Response.Header.Set(`Cache-Control`, c.HeaderCacheControl)
}

func (c *cacheFilter) fillResponseFromCache(ctx *fasthttp.RequestCtx, uriString string) (item *cacheItem, applied bool) {
	cacheItemI, found := c.Storage.Get(uriString)

	if !found {
		// nothing in cache
		return
	}
	item = cacheItemI.(*cacheItem)

	// Set the response from the cache
	err := item.Apply(&ctx.Response)
	if err == ErrCacheItemExpired {
		// The item is expired and cannot be applied
		return
	}
	// Check if a real error occured
	if err != nil {
		ctx.Logger().Printf(`[cache] unexpected error: %v`, err)
		ctx.Response.ResetBody()
		return
	}
	applied = true

	// end
	return
}

func (c *cacheFilter) saveToCacheFromResponse(cacheItem *cacheItem, ctx *fasthttp.RequestCtx, uriString string) {
	if estimateCacheItemSizeFor(ctx) > c.MaxEntrySize {
		if cacheItem != nil {
			c.Storage.Remove(uriString)
		}
		return
	}

	alreadyAdded := cacheItem != nil
	if !alreadyAdded {
		cacheItem = newCacheItem()
	}

	cacheItem.FillFrom(ctx)
	cacheItem.ExpireTS = time.Now().Add(c.MaxDuration)

	if !alreadyAdded {
		c.Storage.Add(uriString, cacheItem)
	}
}
