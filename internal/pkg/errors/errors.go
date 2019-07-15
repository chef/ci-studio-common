package errors

import (
	"github.com/pkg/errors"
)

// LayeredWrap will conditionally layer two error responses together
func LayeredWrap(upstreamErr error, upstreamReason string, recentErr error, recentReason string) error {
	if upstreamErr == nil {
		return errors.Wrap(recentErr, recentReason)
	}

	parentErr := errors.Wrap(upstreamErr, upstreamReason)
	return errors.Wrapf(parentErr, "[on cleanup after error] %s", recentReason)
}
