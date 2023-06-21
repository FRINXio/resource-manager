package lock

import (
	"sync"
	"time"
)

// This mutex handle synchronized access to the locks' storage.
// It is used as prevention before race conditions.
var (
	mutex = &sync.Mutex{}
)

// LocksStorage provides a storage mechanism for the locks.
// It provides methods to acquire and release locks.
type LocksStorage interface {
	AcquireLock(lockName string) *sync.Mutex
	ReleaseLock(lockName string)
}

// LockItem is a structure that represents a lock.
// It contains a testID, a mutex and a timestamp.
// The timestamp is used to indicate when the lock was last used.
// We use timestamp to help us to know which locks are not used for a long time.
type LockItem struct {
	name      string
	mutex     *sync.Mutex
	timestamp int64
}

type LocksStorageImpl struct {
	locks           *sync.Map
	invalidateAfter time.Duration
}

func NewLocksStorageImpl(invalidateAfter time.Duration) LocksStorage {
	locks := &sync.Map{}
	locksStorage := &LocksStorageImpl{locks: locks, invalidateAfter: invalidateAfter}
	go runCleanerJob(locksStorage)

	return locksStorage
}

// AcquireLock acquires mutex by lockName.
// If the lockName is not in the storage, it creates a new lock with new mutex and set pointer to it.
// If the lockName is in the storage, it returns the mutex.
// It also updates the timestamp of the lock to the current time.
func (l *LocksStorageImpl) AcquireLock(lockName string) *sync.Mutex {
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

// ReleaseLock releases a lockName.
// It sets the timestamp of the lock to 0 and sets the mutex to nil.
// It also removes the lockName from the storage.
func (l *LocksStorageImpl) ReleaseLock(lockName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if lockItem, ok := l.locks.Load(lockName); ok {
		lockItem.(*LockItem).timestamp = 0
		lockItem.(*LockItem).mutex = nil
		l.locks.Delete(lockName)
	}
}

// Clean removes locks that are not used for a long time.
// It is based on the timestamp.
// It is called in the background by the locks storage (executed when creating new locks storage).
// It is called every 5 seconds.
// invalidateAfter is set when creating new locks storage. It is set by developer.
func clean(l *LocksStorageImpl) {
	mutex.Lock()
	defer mutex.Unlock()

	timeNow := time.Now().Unix()

	l.locks.Range(func(key, value interface{}) bool {
		if timeNow-value.(*LockItem).timestamp > int64(l.invalidateAfter.Seconds()) {
			load, _ := l.locks.LoadAndDelete(key)
			if load.(LockItem).mutex.TryLock() {
				// mutex was unlocked and invalid, we can delete
				load.(LockItem).mutex.Unlock()
			} else {
				// mutex is locked and invalid, what now ?
				// keep it in locks map]
				//log.Warn("Lock taking too long")
				l.locks.Store(key, load)
			}

			// mutex is unlocked but used / referenced from somewhere else -> this is a known rare race condition we will live with
		}

		return true
	})
}

// runCleanerJob is used to start a cleaner job that will remove locks that are not used longer then their invalidateAfter.
func runCleanerJob(l *LocksStorageImpl) {
	for {
		time.Sleep(time.Millisecond) // TODO make configurable ... default 10 seconds
		clean(l)
	}
}
