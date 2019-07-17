package pushover

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// ValidateRequest is the data to POST to the Pushover
// REST API. See the Pushover Validate API documentation
// for more information on these parameters.
type ValidateRequest struct {
	// The URL for the Pushover REST API POST.
	//
	// Leave this empty unless you wish to override the URL.
	PushoverURL string

	// Required fields

	// Pushover API token
	Token string

	// The user's token to validate
	User string

	// User device to validate (optional)
	Device string
}

// ValidateResponse is the response from this API. It is read from
// the body of the Pushover REST API response and translated
// to this response structure.
//
// For access to the original, untranslated response, access
// the ResponseBody field.
type ValidateResponse struct {
	// Original response body from POST
	ResponseBody string

	// HTTP Status string
	HTTPStatus string

	// HTTP Status Code
	HTTPStatusCode int

	// The status as returned by the Pushover API.
	//
	// Value of 1 indicates 200 response received.
	// Any other value indicates an error with the
	// input.
	APIStatus int

	// This field is returned but not documented
	// by the Pushover Validate API
	Group int

	// ID assigned to the request by Pushover
	Request string

	// List of registered devices
	Devices []string

	// List of licensed platforms
	Licenses []string

	// List of errors returned
	//
	// Empty if no errors
	Errors []string

	// Map of parameters and corresponding errors
	//
	// Empty if no errors
	ErrorParameters map[string]string
}

func validateWithoutValidation(request ValidateRequest) (*ValidateResponse, error) {
	if len(request.PushoverURL) == 0 {
		request.PushoverURL = validateURL
	}

	formData := url.Values{
		keyToken: {request.Token},
		keyUser:  {request.User},
	}

	if len(request.Device) > 0 {
		formData.Set(keyDevice, request.Device)
	}

	resp, err := http.PostForm(request.PushoverURL, formData)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	r := new(ValidateResponse)

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, ErrInvalidResponse
	}

	r.ResponseBody = body.String()
	r.HTTPStatus = resp.Status
	r.HTTPStatusCode = resp.StatusCode

	// Decode json response
	var result map[string]interface{}
	if e := json.NewDecoder(strings.NewReader(string(r.ResponseBody))).Decode(&result); e != nil {
		return nil, ErrInvalidResponse
	}

	var ok bool

	// Populate request status
	if r.APIStatus, ok = mapKeyToInt(keyStatus, result); !ok {
		return nil, ErrInvalidResponse
	}
	delete(result, keyStatus)

	// Populate request ID
	if r.Request, ok = result[keyRequest].(string); !ok {
		return nil, ErrInvalidResponse
	}
	delete(result, keyRequest)

	// Populate group
	if r.Group, ok = mapKeyToInt(keyGroup, result); ok {
		delete(result, keyGroup)
	}

	// Populate licenses
	if r.Licenses, err = interfaceArrayToStringArray(keyLicenses, result); err != nil {
		return nil, err
	}
	delete(result, keyLicenses)

	// Populate devices
	if r.Devices, err = interfaceArrayToStringArray(keyDevices, result); err != nil {
		return nil, err
	}
	delete(result, keyDevices)

	// Populate errors
	if r.Errors, err = interfaceArrayToStringArray(keyErrors, result); err != nil {
		return nil, err
	}
	delete(result, keyErrors)

	// Populate parameters with corresponding errors
	if r.ErrorParameters, err = interfaceMapToStringMap(result); err != nil {
		return nil, err
	}

	return r, nil
}

// Validate will submit a POST request to the Pushover
// Validate API after validating the required fields are
// present This function will check a user or group token
// to determine if it is valid.
//
// The required fields are: Token, User
//   resp, err := pushover.Validate(pushover.ValidateRequest{
//	   Token:   token,
//	   User:    user,
//   })
func Validate(request ValidateRequest) (*ValidateResponse, error) {
	// Validate Token
	if len(request.Token) == 0 {
		return nil, ErrInvalidToken
	}

	// Validate User
	if len(request.User) == 0 {
		return nil, ErrInvalidUser
	}

	return validateWithoutValidation(request)
}
