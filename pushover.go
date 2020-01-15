// Package pushover implements access to the Pushover API
//
// This documentation can be considered a supplement to the
// official Pushover API documentation at
// https://pushover.net/api. Refer to that official
// documentation for the details on how to us these
// library functions.
package pushover

import (
	"fmt"
)

const (
	keyDevice    = "device"
	keyDevices   = "devices"
	keyErrors    = "errors"
	keyExpire    = "expire"
	keyGroup     = "group"
	keyHTML      = "html"
	keyLicenses  = "licenses"
	keyMessage   = "message"
	keyMonospace = "monospace"
	keyPriority  = "priority"
	keyReceipt   = "receipt"
	keyRequest   = "request"
	keyRetry     = "retry"
	keySound     = "sound"
	keyStatus    = "status"
	keyTimestamp = "timestamp"
	keyTitle     = "title"
	keyToken     = "token"
	keyURL       = "url"
	keyURLTitle  = "url_title"
	keyUser      = "user"
)

// ErrInvalidRequest indicates invalid request data
// was sent to a library function
type ErrInvalidRequest struct{}

func (ir *ErrInvalidRequest) Error() string {
	return "Invalid request"
}

// ErrInvalidResponse indicates an invalid response body
// was received from the Pushover API
type ErrInvalidResponse struct{}

func (ir *ErrInvalidResponse) Error() string {
	return "Invalid response"
}

var messagesURL = "https://api.pushover.net/1/messages.json"
var validateURL = "https://api.pushover.net/1/users/validate.json"

func mapKeyToInt(key string, m map[string]interface{}) (int, bool) {
	var value float64
	var result int
	var ok bool

	if value, ok = m[key].(float64); ok {
		result = int(value)
	}

	return result, ok
}

func interfaceArrayToStringArray(key string, m map[string]interface{}) []string {
	var interfaceArray []interface{}
	var stringArray []string
	var ok bool

	if interfaceArray, ok = m[key].([]interface{}); ok {
		stringArray = make([]string, len(interfaceArray))
		for i, v := range interfaceArray {
			stringArray[i] = fmt.Sprintf("%v", v)
		}
	} else {
		stringArray = []string{}
	}

	return stringArray
}

func interfaceMapToStringMap(inMap map[string]interface{}) map[string]string {
	outMap := make(map[string]string)

	for k, v := range inMap {
		outMap[k] = fmt.Sprintf("%v", v)
	}

	return outMap
}
