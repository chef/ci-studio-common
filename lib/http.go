package lib

// Based on https://gist.github.com/josue/598ac048f78362f593f49436680b60b7

import (
	"log"
	"net/http"
	"strings"
)

/*
GetURLHeaders returns a map array of all available headers.
@param string - URL given
@return map[string]interface{}
*/
func GetURLHeaders(url string) map[string]interface{} {
	response, err := http.Head(url)
	if err != nil {
		log.Fatal("Error: Unable to download URL (", url, ") with error: ", err)
	}

	if response.StatusCode != http.StatusOK {
		log.Fatal("Error: HTTP Status = ", response.Status)
	}

	headers := make(map[string]interface{})

	for k, v := range response.Header {
		headers[strings.ToLower(k)] = string(v[0])
	}

	return headers
}

/*
GetURLHeaderByKey returns the header value from a given header key, if available, else returns empty string.
@param string - URL given
@param string - Header key
@return string
*/
func GetURLHeaderByKey(url string, key string) string {
	headers := GetURLHeaders(url)
	key = strings.ToLower(key)

	if value, ok := headers[key]; ok {
		return value.(string)
	}

	return ""
}
