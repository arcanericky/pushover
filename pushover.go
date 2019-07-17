// Package pushover implements access to the Pushover API
//
// This documentation can be considered a supplement to the
// official Pushover API documentation at
// https://pushover.net/api. Refer to that official
// documentation for the details on how to us these
// library functions.
package pushover

import "errors"

const (
	keyDevice    = "device"
	keyDevices   = "devices"
	keyErrors    = "errors"
	keyGroup     = "group"
	keyHTML      = "html"
	keyLicenses  = "licenses"
	keyMessage   = "message"
	keyMonospace = "monospace"
	keyPriority  = "priority"
	keyRequest   = "request"
	keySound     = "sound"
	keyStatus    = "status"
	keyTimestamp = "timestamp"
	keyTitle     = "title"
	keyToken     = "token"
	keyURL       = "url"
	keyURLTitle  = "url_title"
	keyUser      = "user"
)

// ErrInvalidToken indicates an invalid token
var ErrInvalidToken = errors.New("Invalid token")

// ErrInvalidUser indicates an invalid user
var ErrInvalidUser = errors.New("Invalid user")

// ErrInvalidMessage indicates invalid message text
// was sent to a library function
var ErrInvalidMessage = errors.New("Invalid message")

// ErrInvalidRequest indicates invalid request data
// was sent to a library function
var ErrInvalidRequest = errors.New("Invalid request")

// ErrInvalidResponse indicates an invalid response body
// was received from the Pushover API
var ErrInvalidResponse = errors.New("Invalid response")

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

func interfaceArrayToStringArray(key string, m map[string]interface{}) ([]string, error) {
	var interfaceArray []interface{}
	var stringArray []string
	var ok bool

	if interfaceArray, ok = m[key].([]interface{}); ok {
		stringArray = make([]string, len(interfaceArray))
		for i, v := range interfaceArray {
			if stringArray[i], ok = v.(string); !ok {
				return nil, ErrInvalidResponse
			}
		}
	} else {
		stringArray = []string{}
	}

	return stringArray, nil
}

func interfaceMapToStringMap(inMap map[string]interface{}) (map[string]string, error) {
	var ok bool
	outMap := make(map[string]string)

	for k, v := range inMap {
		if outMap[k], ok = v.(string); !ok {
			return nil, ErrInvalidResponse
		}
	}

	return outMap, nil
}
