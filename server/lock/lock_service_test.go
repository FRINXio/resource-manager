package lock

import (
	"sync"
	"testing"
	"time"
)

const testID = "testID"

// basic test unlock
func TestLocksStorageImpl_AcquireLock(t *testing.T) {
	runJobEvery := time.Second * 2
	lockStorage := NewLockingService(time.Second*5, &runJobEvery)
	lockStorage.Acquire(testID).Lock()

	// assert
	value, _ := lockStorage.(*LockingServiceImpl).locks.Load(testID)
	lock := value.(*lockEntry).mutex.TryLock()
	if lock == true {
		t.Fatalf("Double locking")
	}

	lockStorage.Unlock(testID)
	lockAgain := value.(*lockEntry).mutex.TryLock()
	if lockAgain == false {
		t.Fatalf("Failed to unlock")
	}
	value.(*lockEntry).mutex.Unlock()

	for i := 0; i < 20; i++ {
		// count numbers of locks in the map
		var length int
		lockStorage.(*LockingServiceImpl).locks.Range(func(k, v interface{}) bool {
			length++
			return true
		})

		if length == 0 {
			// Mutex removed by cleanup
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("Mutex has not been removed by cleanup")
}

// concurrent locking + unlocking + clean
func TestLocksStorageImpl_Concurrent(t *testing.T) {
	runJobEvery := time.Second
	lck := NewLockingService(time.Second*10, &runJobEvery)
	var wg sync.WaitGroup
	workers := 10000
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(lck, &wg, 100)
	}

	wg.Wait()

	if counter != workers*100 {
		t.Fatalf("Lock service expected to deliver %d value of counter", workers*100)
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
		l.Unlock(testID)
	}

	wg.Done()
}
