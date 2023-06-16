package lock

import (
	"sync"
	"time"
)

type LockingService struct {
	locks *LocksStorageImpl
}

func NewLockingService() *LockingService {
	return &LockingService{locks: NewLocksStorage(2 * time.Minute)}
}

func (l *LockingService) Lock(lockName string) (bool, error) {
	mtx, ok := l.locks.Load(lockName)

	if !ok {
		mtx = &sync.Mutex{}
		l.locks.Store(lockName, mtx)
	} else {
		l.locks.UpdateTimestamp(lockName)
	}

	mtx.Lock()

	return true, nil
}

func (l *LockingService) Unlock(lockName string) (bool, error) {
	if mtx, ok := l.locks.Load(lockName); !ok {
		return false, nil
	} else {
		mtx.Unlock()
	}
	return true, nil
}
