package lock

import (
	"sync"
	"time"
)

type LocksStorage interface {
	Store(lockName string, mtx *sync.Mutex)
	Load(lockName string) (mtx *sync.Mutex, ok bool)
	Delete(lockName string)
	UpdateTimestamp(lockName string)
}

type LocksMap struct {
	mutex     *sync.Mutex
	timestamp int64
}

type LocksStorageImpl struct {
	locks           sync.Map
	invalidateAfter time.Duration
}

func NewLocksStorage(invalidateAfter time.Duration) *LocksStorageImpl {
	lockingServiceCleaner := &LocksStorageImpl{locks: sync.Map{}, invalidateAfter: invalidateAfter}

	go startCleanerJob(lockingServiceCleaner)

	return lockingServiceCleaner
}

func (l *LocksStorageImpl) Store(lockName string, mtx *sync.Mutex) {
	l.locks.Store(lockName, LocksMap{mutex: mtx, timestamp: time.Now().Unix()})
}

func (l *LocksStorageImpl) Load(lockName string) (mtx *sync.Mutex, ok bool) {
	if v, ok := l.locks.Load(lockName); ok {
		mtx = v.(LocksMap).mutex
		ok = true

		return mtx, ok
	}

	return nil, false
}

func (l *LocksStorageImpl) Delete(lockName string) {
	l.locks.Delete(lockName)
}

func (l *LocksStorageImpl) Clean() {
	timeNow := time.Now().Unix()

	l.locks.Range(func(key, value interface{}) bool {
		if timeNow-value.(LocksMap).timestamp > int64(l.invalidateAfter.Seconds()) {
			l.locks.Delete(key)
		}

		return true
	})
}

func (l *LocksStorageImpl) UpdateTimestamp(lockName string) {
	if v, ok := l.locks.Load(lockName); ok {
		l.locks.Store(lockName, LocksMap{mutex: v.(LocksMap).mutex, timestamp: time.Now().Unix()})
	}
}

func startCleanerJob(l *LocksStorageImpl) {
	for {
		time.Sleep(time.Minute * 5)
		l.Clean()
	}
}
