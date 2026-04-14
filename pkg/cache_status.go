package pkg

import (
	"sync"
	"time"
)

type RemoteCacheInfo struct {
	Mirror   string
	Used     bool      // true if cached data was used
	Created  time.Time // cache creation time (if Used) or last fetch time
	Forced   bool      // true when --refresh was requested
}

var (
	cacheMu sync.RWMutex
	cache   RemoteCacheInfo
)

func setRemoteCacheInfo(info RemoteCacheInfo) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache = info
}

func RemoteCacheInfoSnapshot() RemoteCacheInfo {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return cache
}
