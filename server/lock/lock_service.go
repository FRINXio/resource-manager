package lock

import (
	"time"
)

// LockingService provides a locking mechanism for the server.
type LockingService struct {
	locks LocksStorage
}

// NewLockingService creates a new locking service.
func NewLockingService() *LockingService {
	return &LockingService{locks: NewLocksStorageImpl(2 * time.Minute)}
}

// Lock acquires mutex by lockName and locks it.
func (l *LockingService) Lock(lockName string) {
	l.locks.AcquireLock(lockName).Lock()
}

// Unlock unlocks a lockName.
func (l *LockingService) Unlock(lockName string) {
	l.locks.AcquireLock(lockName).Unlock()
}
