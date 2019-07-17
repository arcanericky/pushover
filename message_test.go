package pushover

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
Messages seen during development

Empty Message
400 Bad Request
{"message":"cannot be blank","errors":["message cannot be blank"],"status":0,"request":"4e8c4667-5904-46b9-a367-ff4869a3f401"}

Invalid Token
400 Bad Request
{"token":"invalid","errors":["application token is invalid"],"status":0,"request":"2e8c0c49-8d8d-46d0-b4a6-020c68755037"}

Invalid User
400 Bad Request
{"user":"invalid","errors":["user identifier is not a valid user, group, or subscribed user key"],"status":0,"request":"3cf7ebfd-6f23-4fee-8189-790a08e96892"}

Fields: Token, User, Message
200 OK
{"status":1,"request":"3f7224d5-c408-422b-be1f-dc0da5172497"}

HTML and Monospace both set to 1
400 Bad Request
{"html":"cannot be set with monospace","monospace":"cannot be set with html","errors":["html and monospace are mutually exclusive"],"status":0,"request":"0acffab5-5023-4d97-8492-f4f5d860c83a"}

Invalid Priority
400 Bad Request
{"priority":"is invalid, can only be -2, -1, 0, 1, or 2","errors":["priority is invalid"],"status":0,"request":"8cb3b0d5-d5b7-4283-9bcf-908dcfacff6b"}
*/

const id = "deadbeef-dead-beef-dead-deadbeefdead"

func serverHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.ParseMultipartForm(0)

	// Check message
	value, _ := r.Form["message"]
	if len(value) == 0 || len(value[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"cannot be blank","errors":["message cannot be blank"],"status":0,"request":"%s"}`, id)
		return
	}

	// Check token
	value, _ = r.Form["token"]
	if len(value) == 0 || len(value[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"token":"invalid","errors":["application token is invalid"],"status":0,"request":"%s"}`, id)
		return
	}

	// Check user
	user, _ := r.Form["user"]
	if len(user) == 0 || len(user[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"user":"invalid","errors":["user identifier is not a valid user, group, or subscribed user key"],"status":0,"request":"%s"}`, id)
		return
	}

	// Check html and monospace
	html, _ := r.Form["html"]
	monospace, _ := r.Form["monospace"]
	if len(html) > 0 && len(monospace) > 0 && html[0] == "1" && monospace[0] == "1" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"html":"cannot be set with monospace","monospace":"cannot be set with html","errors":["html and monospace are mutually exclusive"],"status":0,"request":"%s"}`, id)
		return
	}

	value, _ = r.Form["priority"]
	if len(value) > 0 && len(value[0]) > 0 {
		priority, err := strconv.Atoi(value[0])
		if err != nil || priority < -2 || priority > 2 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"priority":"is invalid, can only be -2, -1, 0, 1, or 2","errors":["priority is invalid"],"status":0,"request":"%s"}`, id)
			return
		}
	}

	if user[0] == "failstatus" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status":"abc","request":"%s"}`, id)
		return
	}

	if user[0] == "failrequest" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status":1,"request":1337}`)
		return
	}

	if user[0] == "failjson" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"priority":"is invalid, can only be -2, -1, 0, 1, or 2","errors":"priority is invalid"],"status":0,"request":"%s"}`, id)
		return
	}

	if user[0] == "failbody" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "1")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":1,"request":"%s"}`, id)
}

func TestPushoverMessage(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer apiServer.Close()

	var request MessageRequest

	// Default Pushover URL
	messagesURL = apiServer.URL
	r, e := messageWithoutValidation(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "message cannot be blank" || r.ErrorParameters["message"] != "cannot be blank" {
		t.Error("Default Pushover URL")
	}

	// Invalid Pushover URL
	request.PushoverURL = "\x7f"
	r, e = messageWithoutValidation(context.TODO(), request)
	if e != ErrInvalidRequest {
		t.Error("Invalid Pushover URL")
	}

	// Handling of no message
	request.PushoverURL = apiServer.URL
	r, e = messageWithoutValidation(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "message cannot be blank" || r.ErrorParameters["message"] != "cannot be blank" {
		t.Error("Handling of no message without validation")
	}

	r, e = Message(request)
	if e != ErrInvalidMessage {
		t.Error("Handling of no message with validation")
	}

	// Handling of no token
	request.Message = "test message"
	r, e = messageWithoutValidation(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "application token is invalid" || r.ErrorParameters["token"] != "invalid" {
		t.Error("Handling of no token without validation")
	}

	r, e = Message(request)
	if e != ErrInvalidToken {
		t.Error("Handling of no token with validation")
	}

	// Handling of no user
	request.Token = "testtoken"
	r, e = messageWithoutValidation(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "user identifier is not a valid user, group, or subscribed user key" ||
		r.ErrorParameters["user"] != "invalid" {
		t.Error("Handling of no user without validation")
	}

	r, e = Message(request)
	if e != ErrInvalidUser {
		t.Error("Handling of no user with validation")
	}

	// Valid submission
	request.User = "testuser"
	r, e = Message(request)
	if r.HTTPStatusCode != http.StatusOK || r.APIStatus != 1 || r.Request != id ||
		len(r.Errors) > 0 || len(r.ErrorParameters) > 0 {
		t.Error("Valid submit data")
	}

	// Invalid API Status in response
	request.User = "failstatus"
	r, e = Message(request)
	if e != ErrInvalidResponse {
		t.Error("Invalid API status in response")
	}

	// Invalid request ID in response
	request.User = "failrequest"
	r, e = Message(request)
	if e != ErrInvalidResponse {
		t.Error("Invalid request ID in response")
	}

	// Invalid json response
	request.User = "failjson"
	r, e = Message(request)
	if e != ErrInvalidResponse {
		t.Error("Invalid response JSON")
	}

	// Invalid body
	request.User = "failbody"
	r, e = Message(request)
	if e != ErrInvalidResponse {
		t.Error("Invalid response body")
	}

	// Load all the fields
	request.User = "testuser"
	request.Title = "title"
	request.URL = "url"
	request.URLTitle = "urlTitle"
	request.HTML = "0"
	request.Monospace = "0"
	request.Sound = "sound"
	request.Device = "device"
	request.Priority = "0"
	request.Timestamp = "timestamp"
	r, e = Message(request)
	if r.HTTPStatusCode != http.StatusOK || r.APIStatus != 1 || r.Request != id ||
		len(r.Errors) > 0 || len(r.ErrorParameters) > 0 {
		t.Error("All fields submitted")
	}

	// Context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Millisecond)
	r, e = MessageContext(ctx, request)
	if e != context.DeadlineExceeded {
		t.Error("Context deadline exceeded")
	}
	cancel()

	// Image attachment
	request.ImageReader = strings.NewReader("image data")
	r, e = Message(request)
	if r.HTTPStatusCode != http.StatusOK || r.APIStatus != 1 || r.Request != id ||
		len(r.Errors) > 0 || len(r.ErrorParameters) > 0 {
		t.Error("Image attachment")
	}

	// Test http.PostForm() returning error
	apiServer.Close()
	r, e = Message(request)
	if e == nil {
		t.Error("No API server")
	}
}
