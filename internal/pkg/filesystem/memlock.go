package filesystem

import (
	"path/filepath"
	"time"

	"github.com/avast/retry-go"
	"github.com/zbiljic/go-filelock"
)

type (
	MemLock struct {
		path   string
		locked bool
		err    error
	}

	// This is a super simple implementation and could probably be done
	// with something better than a map, but heh its mostly for testing
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

func (o *MemMapLock) GetRetryAttempts() uint {
	return o.RetryAttempts
}

func (o *MemMapLock) GetRetryDelay() time.Duration {
	return o.RetryDelay
}

func (o *MemMapLock) GetRetryDelayType() retry.DelayTypeFunc {
	return o.RetryDelayType
}

func (f *MemMapLock) GetLock(filename string) (filelock.TryLockerSafe, error) {
	lock, exists := f.Locks[filename]

	if exists {
		return lock, f.Err
	}

	lock = &MemLock{
		path:   filename,
		locked: false,
		err:    f.LockErr,
	}

	f.Locks[filename] = lock
	return lock, f.Err
}

func (f *MemLock) String() string {
	return filepath.Base(f.path)
}

func (f *MemLock) TryLock() (bool, error) {
	if f.locked {
		return false, ErrLocked
	}

	f.locked = true
	return true, f.err
}

func (f *MemLock) Lock() error {
	f.locked = true
	return f.err
}

func (f *MemLock) Unlock() error {
	f.locked = false
	return f.err
}

func (f *MemLock) Must() filelock.TryLocker {
	return nil
}

func (f *MemLock) Destroy() error {
	return f.err
}
