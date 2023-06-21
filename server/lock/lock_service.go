package lock

import "time"

// LockingService provides a locking mechanism for the server.
type LockingService struct {
	locks LocksStorage
}

// NewLockingService creates a new locking service.
func NewLockingService(lockTimeout time.Duration) *LockingService {
	return &LockingService{locks: NewLocksStorageImpl(lockTimeout)}
}

// Lock acquires mutex by lockName and locks it.
func (l *LockingService) Lock(lockName string) {
	l.locks.AcquireLock(lockName).Lock()
}

// Unlock unlocks a lockName.
func (l *LockingService) Unlock(lockName string) {
	l.locks.AcquireLock(lockName).Unlock()
}

// TODOs
// fix unlock - do not call accquire lock to create a new mutex
// merge lock service and lock storage
// tests

// worker acquire lock1
// storage create lock1

// sleep

// cleaner delete lock 1

// worker lock1.lock()

// worker unlock optional
