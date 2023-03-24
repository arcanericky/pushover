package pushover

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// MessageRequest is the data for the POST to the Pushover
// REST API. Some fields in this request should contain
// numbers but the Pushover API parameters are strings.
// There is no validation performed for these fields. If
// invalid data is submitted to the Pushover API, it will
// be rejected with an appropriate error.
//
// See the Pushover API documentation for
// more information on these parameters.
type MessageRequest struct {
	// The URL for the Pushover REST API POST.
	//
	// Leave this empty unless you wish to override the URL.
	PushoverURL string

	// Required fields

	// Pushover API token
	Token string

	// The user's token for message delivery
	User string

	// The message sent to the user
	Message string

	// Optional Fields

	// Message title
	Title string

	// Embedded URL
	URL string

	// The displayed text for the URL
	//
	// If the field URL is missing, the title will be
	// displayed as normal text
	URLTitle string

	// If enabled together will be rejected by Pushover

	// Enable HTML formatting of the message
	//
	// See the Pushover REST API documentation for allowed tags
	HTML string

	// Enable monospace formatting of the message
	Monospace string

	// Sound name for the sound on the user's device
	//
	// See the Pushover REST API documentation for valid
	// values. Invalid sound names will not be rejected by
	//Pushover
	Sound string

	// The device to send the message to rather than all the
	// user's devices.
	//
	// Devices not registered will not be rejected by Pushover
	// and will therefore fail silently.
	Device string

	// Priority number for the message
	//
	// See the Pushover REST API documentation for values and
	// what they mean
	//
	// Invalid priority numbers will be rejected by Pushover
	Priority string

	// How often in seconds the Pushover servers will send
	// the same notification to the user
	//
	// Must be set when Priority is set to "2" and must
	// have a value of at least 30 seconds between retries
	Retry string

	// How many seconds your notification will continue to
	// be retried for (every retry seconds)
	//
	// Must be set when Priority is set to "2" and must
	// have a maximum value of at most 10800 seconds
	// (3 hours)
	Expire string

	// Callback url for the message
	//
	// Optional be set when Priority is set to "2"
	Callback string

	// Unix timestamp for the message rather than the time
	// the message was received by the Pushover REST API
	//
	// Invalid timestamps will not be rejected by Pushover
	Timestamp string

	// Reader for image (attachment) data
	ImageReader io.Reader

	// Optional image name
	//
	// Leave blank to default to image.jpg
	ImageName string
}

// MessageResponse is the response from this API. It is read from
// the body of the Pushover REST API response and translated
// to this response structure.
//
// For access to the original, untranslated response, access
// the ResponseBody field.
type MessageResponse struct {
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

	// ID assigned to the request by Pushover
	Request string

	// Receipt
	//
	// When a priority of 2 is given, a receipt field
	// is returned
	Receipt string

	// List of errors returned
	//
	// Empty if no errors
	Errors []string

	// Map of parameters and corresponding errors
	//
	// Empty if no errors
	ErrorParameters map[string]string
}

// MessageContext will submit a request to the Pushover
// Message API. This function will send a
// message, triggering a notification on a user's
// device or a group's devices.
//
//	  resp, err := pushover.MessageContext(context.Background(),
//	    pushover.MessageRequest{
//		     Token:   token,
//		     User:    user,
//		     Message: message,
//	  })
func MessageContext(ctx context.Context, request MessageRequest) (*MessageResponse, error) {
	var requestData io.Reader
	var contentType string

	if len(request.PushoverURL) == 0 {
		request.PushoverURL = messagesURL
	}

	if len(request.ImageName) == 0 {
		request.ImageName = "image.jpg"
	}

	fields := []struct {
		field string
		value string
	}{
		{field: keyToken, value: request.Token},
		{field: keyUser, value: request.User},
		{field: keyMessage, value: request.Message},
		{field: keyTitle, value: request.Title},
		{field: keyURL, value: request.URL},
		{field: keyURLTitle, value: request.URLTitle},
		{field: keyHTML, value: request.HTML},
		{field: keyMonospace, value: request.Monospace},
		{field: keySound, value: request.Sound},
		{field: keyDevice, value: request.Device},
		{field: keyPriority, value: request.Priority},
		{field: keyRetry, value: request.Retry},
		{field: keyExpire, value: request.Expire},
		{field: keyCallback, value: request.Callback},
		{field: keyTimestamp, value: request.Timestamp},
	}

	if request.ImageReader == nil {
		formData := url.Values{}
		for _, v := range fields {
			if len(v.value) > 0 {
				formData.Set(v.field, v.value)
			}
		}

		requestData = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else {
		requestBody := &bytes.Buffer{}
		writer := multipart.NewWriter(requestBody)
		part, _ := writer.CreateFormFile("attachment", request.ImageName)
		_, _ = io.Copy(part, request.ImageReader)

		for _, v := range fields {
			if len(v.value) > 0 {
				_ = writer.WriteField(v.field, v.value)
			}
		}
		writer.Close()

		requestData = requestBody
		contentType = writer.FormDataContentType()
	}

	req, err := http.NewRequest(http.MethodPost, request.PushoverURL, requestData)
	if err != nil {
		return nil, &ErrInvalidRequest{}
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	r := new(MessageResponse)

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, &ErrInvalidResponse{}
	}

	r.ResponseBody = body.String()
	r.HTTPStatus = resp.Status
	r.HTTPStatusCode = resp.StatusCode

	// Decode json response
	var result map[string]interface{}
	if e := json.NewDecoder(strings.NewReader(string(r.ResponseBody))).Decode(&result); e != nil {
		return nil, &ErrInvalidResponse{}
	}

	var ok bool

	// Populate request status
	if r.APIStatus, ok = mapKeyToInt(keyStatus, result); !ok {
		return nil, &ErrInvalidResponse{}
	}
	delete(result, keyStatus)

	// Populate request ID
	if r.Request, ok = result[keyRequest].(string); !ok {
		return nil, &ErrInvalidResponse{}
	}
	delete(result, keyRequest)

	// Populate receipt
	if r.Receipt, ok = result[keyReceipt].(string); ok {
		delete(result, keyReceipt)
	}

	// Populate errors
	r.Errors = interfaceArrayToStringArray(keyErrors, result)
	delete(result, keyErrors)

	// Populate parameters with corresponding errors
	r.ErrorParameters = interfaceMapToStringMap(result)

	return r, nil
}

// Message will submit a request to the Pushover
// Message API. This function will send a
// message, triggering a notification on a user's
// device or a group's devices.
//
//	  resp, err := pushover.Message(pushover.MessageRequest{
//		     Token:   token,
//		     User:    user,
//		     Message: message,
//	  })
func Message(request MessageRequest) (*MessageResponse, error) {
	return MessageContext(context.Background(), request)
}
