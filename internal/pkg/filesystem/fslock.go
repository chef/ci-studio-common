package filesystem

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/zbiljic/go-filelock"
)

// ErrLocked - the filelock error.
var ErrLocked = filelock.ErrLocked

// Locker is an interface of filelock.
type Locker interface {
	GetLock(filename string) (filelock.TryLockerSafe, error)
	GetRetryAttempts() uint
	GetRetryDelay() time.Duration
	GetRetryDelayType() retry.DelayTypeFunc
}

// OsLock defines what is available for retries.
type OsLock struct {
	RetryAttempts  uint
	RetryDelay     time.Duration
	RetryDelayType retry.DelayTypeFunc
}

// GetLock creates a new filelock.
func (o *OsLock) GetLock(filename string) (filelock.TryLockerSafe, error) {
	return filelock.New(filename)
}

// GetRetryAttempts returns the max number of retries allowed.
func (o *OsLock) GetRetryAttempts() uint {
	return o.RetryAttempts
}

// GetRetryDelay returns the timed allowed between retries.
func (o *OsLock) GetRetryDelay() time.Duration {
	return o.RetryDelay
}

// GetRetryDelayType returns the type of retry delay.
func (o *OsLock) GetRetryDelayType() retry.DelayTypeFunc {
	return o.RetryDelayType
}
