package http

import (
	"net/http"

	"github.com/pkg/errors"
)

/*
GetURLHeaders returns a map array of all available headers.
@param string - URL given
@return *http.Header
*/
func GetURLHeaders(url string) (http.Header, error) {
	response, err := http.Head(url)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to download URL (%s)", url)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "http status = %s", response.Status)
	}

	return response.Header, nil
}
