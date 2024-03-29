package pushover

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
User Valid
{"status":1,"group":0,"devices":["pixel2xl"],"licenses":["Android"],"request":"6bbfa652-b15b-486e-ac8d-238639ead6a2"}

User Invalid
{"user":"invalid","errors":["user key is invalid"],"status":0,"request":"4fecdd5d-6e46-486d-8f3e-23a07f94695e"}

Token Invaild
{"token":"invalid","errors":["application token is invalid"],"status":0,"request":"e8488e4e-dbe1-4795-a253-3ef644aa14a6"}
*/

//const id = "deadbeef-dead-beef-dead-deadbeefdead"

func validateServerHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	// Check token
	value := r.Form["token"]
	if len(value) == 0 || len(value[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"token":"invalid","errors":["application token is invalid"],"status":0,"request":"%s"}`, id)
		return
	}

	// Check user
	user := r.Form["user"]
	if len(user) == 0 || len(user[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"user":"invalid","errors":["user key is invalid"],"status":0,"request":"%s"}`, id)
		return
	}

	if user[0] == "failstatus" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"abc","group":0,"devices":["pixel2xl"],"licenses":["Android"],"request":"%s"}`, id)
		return
	}

	if user[0] == "failrequest" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":1,"group":0,"devices":["pixel2xl"],"licenses":["Android"],"request":1337}`)
		return
	}

	if user[0] == "failjson" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"priority":1337,"errors":"priority is invalid"],"status":0,"request":"%s"}`, id)
		return
	}

	if user[0] == "failbody" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "1")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":1,"group":0,"devices":["pixel2xl"],"licenses":["Android"],"request":"%s"}`, id)
}

func TestPushoverValidate(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(validateServerHandler))
	defer apiServer.Close()

	var request ValidateRequest

	// Default Pushover URL
	validateURL = apiServer.URL
	r, e := ValidateContext(context.TODO(), request)
	if e != nil || r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "application token is invalid" || r.ErrorParameters["token"] != "invalid" {
		t.Error("Default Pushover URL")
	}

	// Handling of no token
	r, _ = ValidateContext(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "application token is invalid" || r.ErrorParameters["token"] != "invalid" {
		t.Error("Handling of no token without validation")
	}

	r, _ = Validate(request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "application token is invalid" || r.ErrorParameters["token"] != "invalid" {
		t.Error("Handling of no token without validation")
	}

	// Handling of no user
	request.Token = "testtoken"
	r, _ = ValidateContext(context.TODO(), request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "user key is invalid" ||
		r.ErrorParameters["user"] != "invalid" {
		t.Error("Handling of no user without validation")
	}

	r, _ = Validate(request)
	if r.HTTPStatusCode != http.StatusBadRequest || r.APIStatus != 0 || r.Request != id ||
		r.Errors[0] != "user key is invalid" ||
		r.ErrorParameters["user"] != "invalid" {
		t.Error("Handling of no user without validation")
	}

	// Valid submission
	request.User = "testuser"
	r, e = Validate(request)
	if e != nil || r.HTTPStatusCode != http.StatusOK || r.APIStatus != 1 || r.Request != id ||
		len(r.Errors) > 0 || len(r.ErrorParameters) > 0 {
		fmt.Println(e)
		t.Error("Valid submit data")
	}

	// Context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Millisecond)
	_, e = ValidateContext(ctx, request)
	if e != context.DeadlineExceeded {
		t.Error("Context deadline exceeded")
	}
	cancel()

	// Invalid API Status in response
	request.User = "failstatus"
	_, e = Validate(request)
	if _, ok := e.(*ErrInvalidResponse); !ok {
		t.Error("Invalid API status in response")
	}

	// Invalid request ID in response
	request.User = "failrequest"
	_, e = Validate(request)
	if _, ok := e.(*ErrInvalidResponse); !ok {
		t.Error("Invalid request ID in response")
	}

	// Invalid json response
	request.User = "failjson"
	_, e = Validate(request)
	if _, ok := e.(*ErrInvalidResponse); !ok {
		t.Error("Invalid response JSON")
	}

	// Load all the fields
	request.User = "testuser"
	request.Device = "device"
	r, e = Validate(request)
	if e != nil || r.HTTPStatusCode != http.StatusOK || r.APIStatus != 1 || r.Request != id ||
		len(r.Errors) > 0 || len(r.ErrorParameters) > 0 {
		t.Error("All fields submitted")
	}

	// Invalid body
	request.User = "failbody"
	_, e = Validate(request)
	if _, ok := e.(*ErrInvalidResponse); !ok {
		t.Error("Invalid response body")
	}

	// Test http.PostForm() returning error
	apiServer.Close()
	_, e = Validate(request)
	if e == nil {
		t.Error("No API server")
	}
}

func TestErrorTypes(t *testing.T) {
	errRequest := &ErrInvalidRequest{}
	if len(errRequest.Error()) == 0 {
		t.Error("ErrInvalidRequest does not return string on Error()")
	}

	errResponse := &ErrInvalidResponse{}
	if len(errResponse.Error()) == 0 {
		t.Error("ErrInvalidResponse does not return string on Error()")
	}
}
