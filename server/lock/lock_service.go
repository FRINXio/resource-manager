package lock

import (
	"context"
	log "github.com/net-auto/resourceManager/logging"
	"sync"
	"time"
)

// This mutex handle synchronized access to the locks' storage.
// It is used as prevention before race conditions.
var (
	mutex = &sync.Mutex{}
)

// LockingService provides a storage mechanism for the locks.
// It provides methods to acquire and release locks.
type LockingService interface {
	Acquire(lockName string) *sync.Mutex
	Unlock(lockName string)
}

// Item is a structure that represents a lock.
// It contains a testID, a mutex and a timestamp.
// The timestamp is used to indicate when the lock was last used.
// We use timestamp to help us to know which locks are not used for a long time.
type LockItem struct {
	name      string
	mutex     *sync.Mutex
	timestamp int64
}

// LockingServiceImpl is a structure that represents a storage for the locks.
// It contains a map of locks and a timestamp.
// The timestamp is used to indicate when the lock was last used.
// We use timestamp to help us to know which locks are not used for a long time.
type LockingServiceImpl struct {
	locks           *sync.Map
	invalidateAfter time.Duration
}

// LockingService provides a locking mechanism for the server.

// NewLockingService creates a new locking service.
func NewLockingService(invalidateAfter time.Duration, runCleanJobEvery *time.Duration) LockingService {
	locks := &sync.Map{}
	lockingService := &LockingServiceImpl{locks: locks, invalidateAfter: invalidateAfter}
	go runCleanerJob(lockingService, runCleanJobEvery)

	return lockingService
}

// Lock acquires mutex by lockName and locks it.
func (l *LockingServiceImpl) Acquire(lockName string) *sync.Mutex {
	mutex.Lock()
	defer mutex.Unlock()

	timeNow := time.Now().Unix()
	var mtx *sync.Mutex

	if lockItem, ok := l.locks.Load(lockName); !ok {
		mtx = &sync.Mutex{}
		lockItem = &LockItem{name: lockName, mutex: mtx, timestamp: timeNow}
		l.locks.Store(lockName, lockItem)
	} else {
		lockItem.(*LockItem).timestamp = timeNow
		mtx = lockItem.(*LockItem).mutex
		l.locks.Store(lockName, lockItem)
	}

	return mtx
}

// Unlock unlocks a lockName.
func (l *LockingServiceImpl) Unlock(lockName string) {
	mutex.Lock()
	defer mutex.Unlock()

	timeNow := time.Now().Unix()
	var mtx *sync.Mutex

	if lockItem, ok := l.locks.Load(lockName); !ok {
		log.GetLogger().Warn("Unlocking a lock that does not exist")
		return
	} else {
		lockItem.(*LockItem).timestamp = timeNow
		mtx = lockItem.(*LockItem).mutex
		l.locks.Store(lockName, lockItem)
	}

	mtx.Unlock()
}

// Clean removes locks that are not used for a long time.
// It is based on the timestamp.
// It is called in the background by the locks storage (executed when creating new locks storage).
// It is called every 5 seconds.
// invalidateAfter is set when creating new locks storage. It is set by developer.
func clean(l *LockingServiceImpl) {
	ctx := context.Background()
	mutex.Lock()
	defer mutex.Unlock()

	timeNow := time.Now().Unix()

	l.locks.Range(func(key, value interface{}) bool {
		if timeNow-value.(*LockItem).timestamp > int64(l.invalidateAfter.Seconds()) {
			load, _ := l.locks.LoadAndDelete(key)
			if load.(*LockItem).mutex.TryLock() {
				// mutex was unlocked and invalid, we can delete
				load.(*LockItem).mutex.Unlock()
			} else {
				// mutex is locked and invalid, what now ?
				// keep it in locks map
				log.Warn(ctx, "Lock taking too long")
				l.locks.Store(key, load)
			}

			// mutex is unlocked but used / referenced from somewhere else -> this is a known rare race condition we will live with
		}

		return true
	})
}

// runCleanerJob is used to start a cleaner job that will remove locks that are not used longer then their invalidateAfter.
func runCleanerJob(l *LockingServiceImpl, runEvery *time.Duration) {
	if runEvery == nil {
		defaultRunEvery := time.Second * 10
		runEvery = &defaultRunEvery
	}

	for {
		time.Sleep(*runEvery)
		clean(l)
	}
}
