package middlewares

import (
	"sync"
)

type cacheURILocker struct {
	mapLocker sync.RWMutex
	lockerMap map[string]*sync.Mutex
}

func newCacheURILocker() *cacheURILocker {
	l := &cacheURILocker{}
	l.lockerMap = map[string]*sync.Mutex{}
	return l
}

func (l *cacheURILocker) Lock(uri string) {
	for {
		// At first we try RLock to reduce a bottleneck in case of too many waiting requests
		l.mapLocker.RLock()
		uriLocker := l.lockerMap[uri]
		l.mapLocker.RUnlock()
		if uriLocker != nil {
			// there's a request already, wait until it ends and try again
			uriLocker.Lock()
			uriLocker.Unlock()
			continue
		}

		// OK, URI is not locked (yet) and we can try to lock it

		l.mapLocker.Lock()
		uriLocker = l.lockerMap[uri]
		if uriLocker != nil {
			// somebody else is already locked the URI :(, wait and try again
			l.mapLocker.Unlock()
			uriLocker.Lock()
			uriLocker.Unlock()
			continue
		}

		uriLocker = &sync.Mutex{}
		l.lockerMap[uri] = uriLocker
		uriLocker.Lock()
		l.mapLocker.Unlock()

		break
	}
}

func (l *cacheURILocker) Unlock(uri string) {
	l.mapLocker.Lock()
	uriLocker := l.lockerMap[uri]
	delete(l.lockerMap, uri)
	l.mapLocker.Unlock()

	uriLocker.Unlock()
}
