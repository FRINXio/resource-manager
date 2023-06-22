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

// lockEntry is a structure that represents a lock.
// It contains a name, a mutex and a timestamp.
// The timestamp is used to indicate when the lock was last used.
// We use timestamp to help us to know which locks should be auto cleaned up.
type lockEntry struct {
	name      string
	mutex     *sync.Mutex
	timestamp int64
}

// LockingServiceImpl is a structure that represents a storage for the locks.
// It contains a map of locks and a timestamp.
// The timestamp is used to indicate when the lock was last used.
// We use timestamp to help us to know which locks should be auto cleaned up.
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

// Acquire creates a new mutex and returns it. In case such mutex already exists, it updates its timestamp..
func (l *LockingServiceImpl) Acquire(lockName string) *sync.Mutex {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debug(context.Background(), "Lock: %s locking.", lockName)
	timeNow := time.Now().Unix()
	var mtx *sync.Mutex

	if lockItem, ok := l.locks.Load(lockName); !ok {
		log.Debug(context.Background(), "Lock: %s timestamp refresh.", lockName)
		mtx = &sync.Mutex{}
		lockItem = &lockEntry{name: lockName, mutex: mtx, timestamp: timeNow}
		l.locks.Store(lockName, lockItem)
	} else {
		lockItem.(*lockEntry).timestamp = timeNow
		mtx = lockItem.(*lockEntry).mutex
		l.locks.Store(lockName, lockItem)
	}

	return mtx
}

// Unlock unlocks a lockName.
func (l *LockingServiceImpl) Unlock(lockName string) {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debug(context.Background(), "Lock: %s unlocking.", lockName)
	timeNow := time.Now().Unix()
	var mtx *sync.Mutex

	if lockItem, ok := l.locks.Load(lockName); !ok {
		log.Warn(context.Background(), "Unlocking a lock: %s that does not exist. Probably removed by cleanUp.", lockName)
		return
	} else {
		lockItem.(*lockEntry).timestamp = timeNow
		mtx = lockItem.(*lockEntry).mutex
		l.locks.Store(lockName, lockItem)
	}

	// Call tryLock to avoid runtime error of unlocking an unlocked mutex
	mtx.TryLock()
	mtx.Unlock()
}

// Clean removes locks that are not used for a long time.
// It is based on the timestamp.
// It is called in the background by the locks storage (executed when creating new locks storage).
// It is called every 5 seconds.
// invalidateAfter is set when creating new locks storage. It is set by developer.
func clean(l *LockingServiceImpl) {
	mutex.Lock()
	defer mutex.Unlock()

	timeNow := time.Now().Unix()

	l.locks.Range(func(key, value interface{}) bool {
		if !isTimedOut(timeNow, value.(*lockEntry), l.invalidateAfter) {
			// lock still valid
			return true
		}

		if value.(*lockEntry).mutex.TryLock() {
			// mutex was unlocked and is timed out, we can safely delete
			log.Debug(context.Background(), "Lock: %s cleaning up.", key)
			value.(*lockEntry).mutex.Unlock()
			l.locks.Delete(key)
		} else {
			// mutex is locked and invalid, what now ?
			// keep it in locks map to retry next time (the request should time out eventually)
			log.Warn(context.Background(), "Lock: %s is locked and timed out. Clean delayed.", key)

			// in case this is taking a long time (mutex locked and timed out) forcefully unlock and refresh timestamp
			if isTimedOut(timeNow, value.(*lockEntry), l.invalidateAfter*2) {
				log.Warn(context.Background(), "Lock: %s is locked and timed out. Unlocking forcefully.", key)
				value.(*lockEntry).mutex.Unlock()
				// Refresh the timestamp
				swap := &lockEntry{name: value.(*lockEntry).name, mutex: value.(*lockEntry).mutex, timestamp: timeNow}
				l.locks.Store(key, swap)
			}
		}

		// mutex is unlocked but used / referenced from somewhere else
		// e.g. a request has acquired a mutex but has not locked it yet
		// -> this is a known rare race condition
		// -> should not happen with reasonable invalidate and cleanAfter settings

		return true
	})
}

func isTimedOut(timeNow int64, value *lockEntry, invalidateAfter time.Duration) bool {
	return timeNow-value.timestamp > int64(invalidateAfter.Seconds())
}

// runCleanerJob is used to start a cleaner job that will remove locks that
// haven't been used (acquired) for over invalidateAfter duration.
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
