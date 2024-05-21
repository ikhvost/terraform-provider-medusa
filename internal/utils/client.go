package utils

import (
	"github.com/ikhvost/medusajs-go-sdk/medusa"
	"net/http"
)

func GetClient(data any) medusa.ClientWithResponsesInterface {
	c, ok := data.(medusa.ClientWithResponsesInterface)
	if !ok {
		panic("invalid client type")
	}
	return c
}

func CleanHeaders(headers http.Header, keep ...string) http.Header {
	for key := range headers {
		if !contains(keep, key) {
			headers.Del(key)
		}
	}
	return headers
}

func contains(keep []string, key string) bool {
	for _, k := range keep {
		if k == key {
			return true
		}
	}
	return false
}
