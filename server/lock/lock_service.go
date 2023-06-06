package lock

import (
	"sync"
	"time"

	cache "github.com/net-auto/resourceManager/cache"
)

type LockingService interface {
	Lock(lockName string, timeout int) (bool, error)
	Unlock(lockName string) (bool, error)
}

type LockingServiceImpl struct {
	locks *cache.BigCacheServiceImpl[sync.Mutex]
}

func NewLockingService() *LockingServiceImpl {
	return &LockingServiceImpl{locks: cache.NewBigCacheService[sync.Mutex](10 * time.Minute)}
}

func (l *LockingServiceImpl) Lock(lockName string) (bool, error) {
	if _, err := l.locks.Get(lockName); err != nil {
		l.locks.Set(lockName, sync.Mutex{})
	}

	mtx, err := l.locks.Get(lockName)

	if err != nil {
		return false, err
	}

	mtx.Lock()
	return true, nil
}

func (l *LockingServiceImpl) Unlock(lockName string) (bool, error) {
	mtx, err := l.locks.Get(lockName)

	if err != nil {
		return false, err
	}

	mtx.Unlock()
	return true, nil
}
