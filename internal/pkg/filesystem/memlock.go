package filesystem

import (
	"path/filepath"
	"time"

	"github.com/avast/retry-go"
	"github.com/zbiljic/go-filelock"
)

type (
	// MemLock - definition of the memory lock.
	MemLock struct {
		path   string
		locked bool
		err    error
	}

	// MemMapLock - This is a super simple implementation and could probably be done
	// with something better than a map, but heh its mostly for testing.
	MemMapLock struct {
		RetryAttempts  uint
		RetryDelay     time.Duration
		RetryDelayType retry.DelayTypeFunc

		Locks map[string]*MemLock

		// The error you want GetLock() to throw
		Err error

		// The error you want TryLock() to throw
		LockErr error
	}
)

// GetRetryAttempts - get retry count.
func (o *MemMapLock) GetRetryAttempts() uint {
	return o.RetryAttempts
}

// GetRetryDelay - time between each retry.
func (o *MemMapLock) GetRetryDelay() time.Duration {
	return o.RetryDelay
}

// GetRetryDelayType - type of retry delay.
func (o *MemMapLock) GetRetryDelayType() retry.DelayTypeFunc {
	return o.RetryDelayType
}

// GetLock - create a new lock.
func (o *MemMapLock) GetLock(filename string) (filelock.TryLockerSafe, error) {
	lock, exists := o.Locks[filename]

	if exists {
		return lock, o.Err
	}

	lock = &MemLock{
		path:   filename,
		locked: false,
		err:    o.LockErr,
	}

	o.Locks[filename] = lock

	return lock, o.Err
}

func (f *MemLock) String() string {
	return filepath.Base(f.path)
}

// TryLock - attempt to get the lock.
func (f *MemLock) TryLock() (bool, error) {
	if f.locked {
		return false, ErrLocked
	}

	f.locked = true

	return true, f.err
}

// Lock - lock the fs.
func (f *MemLock) Lock() error {
	f.locked = true

	return f.err
}

// Unlock - unlock the fs.
func (f *MemLock) Unlock() error {
	f.locked = false

	return f.err
}

// Must - nil.
func (f *MemLock) Must() filelock.TryLocker {
	return nil
}

// Destroy - destroy the lock.
func (f *MemLock) Destroy() error {
	return f.err
}
