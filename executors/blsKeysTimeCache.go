package executors

import "sync"

type faultyBLSKeyInfo struct {
	numSnoozeEvents uint32
	addedTimestamp  int64
}

type blsKeysTimeCache struct {
	maxSnoozeEvents     uint32
	cacheExpiration     uint64
	getCurrentTimestamp func() int64
	mutFaultyBLSKeys    sync.Mutex
	faultyBLSKeys       map[string]*faultyBLSKeyInfo
}

// ArgsBlsKeysTimeCache is the argument DTO used for the NewBLSKeysTimeCache constructor function
type ArgsBlsKeysTimeCache struct {
	MaxSnoozeEvents     uint32
	CacheExpiration     uint64
	GetCurrentTimestamp func() int64
}

// NewBLSKeysTimeCache creates a new instance of type BLS keys time cache
func NewBLSKeysTimeCache(args ArgsBlsKeysTimeCache) (*blsKeysTimeCache, error) {
	if args.GetCurrentTimestamp == nil {
		return nil, errNilCurrentTimestampHandler
	}

	return &blsKeysTimeCache{
		maxSnoozeEvents:     args.MaxSnoozeEvents,
		cacheExpiration:     args.CacheExpiration,
		getCurrentTimestamp: args.GetCurrentTimestamp,
		faultyBLSKeys:       make(map[string]*faultyBLSKeyInfo),
	}, nil
}

// ShouldNotify returns true if the current BLS key should trigger a notification event
func (timeCache *blsKeysTimeCache) ShouldNotify(blsKey string) bool {
	timeCache.mutFaultyBLSKeys.Lock()
	defer timeCache.mutFaultyBLSKeys.Unlock()

	currentTimestamp := timeCache.getCurrentTimestamp()

	info := timeCache.faultyBLSKeys[blsKey]
	if info == nil {
		info = &faultyBLSKeyInfo{
			addedTimestamp: currentTimestamp,
		}
		timeCache.faultyBLSKeys[blsKey] = info
	}
	if info.addedTimestamp+int64(timeCache.cacheExpiration) <= currentTimestamp {
		info.numSnoozeEvents = 0
		info.addedTimestamp = currentTimestamp
	}

	shouldTrigger := info.numSnoozeEvents < timeCache.maxSnoozeEvents
	info.numSnoozeEvents++

	return shouldTrigger
}

// IsInterfaceNil returns true if there is no value under the interface
func (timeCache *blsKeysTimeCache) IsInterfaceNil() bool {
	return timeCache == nil
}
