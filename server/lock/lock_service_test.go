package lock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const testID = "testID"

// basic test unlock
func TestLocksStorageImpl_AcquireLock(t *testing.T) {
	runJobEvery := time.Second * 5
	lockStorage := NewLockingService(time.Second*10, &runJobEvery)
	lockStorage.Acquire(testID).Lock()

	// assert
	value, _ := lockStorage.(*LockingServiceImpl).locks.Load(testID)
	lock := value.(*LockItem).mutex.TryLock()
	if lock == true {
		t.Fatalf("Double locking")
	}

	lockStorage.Unlock(testID)
}

// concurrent locking + unlocking + clean
func TestLocksStorageImpl_Concurrent(t *testing.T) {
	runJobEvery := time.Second
	lck := NewLockingService(time.Second*2, &runJobEvery)
	var wg sync.WaitGroup
	workers := 1000
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(lck, &wg, 100)
	}

	wg.Wait()

	fmt.Println("counter: ", counter)

	if counter != workers*100 {
		t.Fatalf("Lock service expected to deliver %d value of counter", workers*100) // FIXME
	}

	return
}

func TestLockStorageImpl_UnlockCleanedMutex(t *testing.T) {
	runJobEvery := time.Millisecond
	lck := NewLockingService(time.Millisecond*2, &runJobEvery)
	lck.Acquire(testID).Lock()
	lck.(*LockingServiceImpl).locks.Delete(testID)
	lck.Unlock(testID)
}

var counter = 0

func worker(l LockingService, wg *sync.WaitGroup, iterations int) {
	for i := 0; i < iterations; i++ {
		l.Acquire(testID).Lock()
		// do work
		counter++
		time.Sleep(time.Millisecond)
		l.Unlock(testID)
	}

	wg.Done()
}
