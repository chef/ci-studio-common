package filesystem

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/zbiljic/go-filelock"
)

var (
	ErrLocked = filelock.ErrLocked
)

type Locker interface {
	GetLock(filename string) (filelock.TryLockerSafe, error)
	GetRetryAttempts() uint
	GetRetryDelay() time.Duration
	GetRetryDelayType() retry.DelayTypeFunc
}

type OsLock struct {
	RetryAttempts  uint
	RetryDelay     time.Duration
	RetryDelayType retry.DelayTypeFunc
}

func (o *OsLock) GetLock(filename string) (filelock.TryLockerSafe, error) {
	return filelock.New(filename)
}

func (o *OsLock) GetRetryAttempts() uint {
	return o.RetryAttempts
}

func (o *OsLock) GetRetryDelay() time.Duration {
	return o.RetryDelay
}

func (o *OsLock) GetRetryDelayType() retry.DelayTypeFunc {
	return o.RetryDelayType
}
