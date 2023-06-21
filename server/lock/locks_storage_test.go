package lock

import (
	"sync"
	"testing"
	"time"
)

const testID = "testID"

func TestLocksStorageImpl_AcquireLock(t *testing.T) {
	lockStorage := NewLocksStorageImpl(time.Second * 10)
	lockStorage.AcquireLock(testID).Lock()

	// assert
	value, _ := lockStorage.(*LocksStorageImpl).locks.Load(testID)
	lock := value.(*LockItem).mutex.TryLock()
	if lock == true {
		t.Fatalf("Double locking")
	}

	// todo add unlock
}

// basic test unlock

// concurrent locking + unlocking + clean
func TestLocksStorageImpl_Concurrent(t *testing.T) {

	lck := NewLockingService(time.Minute)
	var wg sync.WaitGroup
	workers := 10000
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(lck, &wg, 100)
	}

	wg.Wait()

	if counter != workers*100 {
		t.Fatalf("") // FIXME
	}
}

var counter = 0

func worker(l *LockingService, wg *sync.WaitGroup, iterations int) {
	for i := 0; i < iterations; i++ {
		l.Lock(testID)
		// do work
		counter++
		time.Sleep(time.Millisecond)
		l.Unlock(testID)
	}

	wg.Done()
}
