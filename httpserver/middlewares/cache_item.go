package middlewares

import (
	"time"
	"unsafe"

	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/errors"
)

var (
	// ErrCacheItemExpired is returned when it was called an `Apply()` on an expired cache item.
	ErrCacheItemExpired = errors.New(`cache item has expired`)
)

var (
	sizeOfCacheItemStruct = uint64(unsafe.Sizeof(cacheItem{}))
)

type cacheItem struct {
	// ExpireTS defines the timestamp after which the cache item is considered expired
	ExpireTS time.Time

	// Body is a copy of the HTTP response body
	Body []byte

	// ContentType is the value of HTTP header `Content-Type`
	ContentType string

	// StatusCode is the HTTP status code
	StatusCode int
}

// newCacheItem returns a new cache item with a copy of the response "resp".
//
// The item will not be valid after the moment `expireTS`.
func newCacheItem() *cacheItem {
	return &cacheItem{}
}

// FillFromResponse will the item by the data from response `resp`.
//
// So if `Apply` will be called it will fill the response by this data.
func (c *cacheItem) FillFrom(ctx *fasthttp.RequestCtx) {
	resp := &ctx.Response

	c.StatusCode = resp.StatusCode()

	body := resp.Body()
	c.Body = make([]byte, len(body))
	copy(c.Body, body)

	c.ContentType = string(resp.Header.Peek(`Content-Type`))
}

func estimateCacheItemSizeFor(ctx *fasthttp.RequestCtx) uint64 {
	resp := &ctx.Response
	return uint64(len(resp.Body())) +
		uint64(len(resp.Header.Peek(`Content-Type`))) +
		sizeOfCacheItemStruct
}

// Size returns an estimation of how much memory this item consumes
func (c *cacheItem) Size() uint64 {
	return uint64(len(c.Body)) + uint64(len(c.ContentType)) + sizeOfCacheItemStruct
}

// Apply sends the response from the cache.
//
// This method returns error ErrCacheItemExpired if used on an expired item.
// Also there could be other errors (for example caused by `resp.BodyWriter().Write()).
func (c *cacheItem) Apply(resp *fasthttp.Response) error {
	if time.Now().After(c.ExpireTS) {
		return ErrCacheItemExpired
	}

	_, err := resp.BodyWriter().Write(c.Body)
	if err != nil {
		return errors.Wrap(err)
	}

	resp.Header.Set(`Content-Type`, c.ContentType)
	resp.Header.Set(`X-Cached-Response`, `true`)
	resp.SetStatusCode(c.StatusCode)
	return nil
}
