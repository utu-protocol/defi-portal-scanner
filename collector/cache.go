package collector

import (
	"strings"
	"sync"
)

var (
	addressCache map[string]string
	addressType  map[string]string
	addressM     sync.RWMutex
)

func init() {
	addressCache = make(map[string]string)
	addressType = make(map[string]string)
}

func cachePush(key, value, typ string) {
	addressM.Lock()
	defer addressM.Unlock()

	k := strings.ToLower(strings.TrimSpace(key))
	addressCache[k] = value
	// only cache the type defiProtocol
	if typ == TypeDefiProtocol {
		addressType[k] = typ
	}
}

func cacheGet(key string) (v string, t string, found bool) {
	addressM.RLock()
	defer addressM.RUnlock()

	k := strings.ToLower(strings.TrimSpace(key))
	v, found = addressCache[k]
	if !found {
		return
	}
	t = TypeAddress
	if typ, hasT := addressType[k]; hasT {
		t = typ
	}
	return
}
