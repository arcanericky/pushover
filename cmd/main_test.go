package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const id = "deadbeef-dead-beef-dead-deadbeefdead"

func serverMessageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.ParseMultipartForm(0)

	// Check html and monospace
	html, _ := r.Form["html"]
	monospace, _ := r.Form["monospace"]
	if len(html) > 0 && len(monospace) > 0 && html[0] == "1" && monospace[0] == "1" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"html":"cannot be set with monospace","monospace":"cannot be set with html","errors":["html and monospace are mutually exclusive"],"status":0,"request":"%s"}`, id)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":1,"request":"%s","receipt":"1337"}`, id)
}

func TestPushoverMessageCLI(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(serverMessageHandler))
	defer apiServer.Close()

	// Test valid input and output
	savedArgs := os.Args
	baseArgs := []string{
		"pushover",
		"message",
		"--pushoverurl", apiServer.URL,
		"--token", "token",
		"--user", "user",
		"--message", "message",
	}

	os.Args = baseArgs

	// Nothing to check - exercising code
	main()

	// Test image attachment with valid file
	os.Args = append(os.Args, "--image", savedArgs[0])

	// Nothing to check - exercising code
	main()

	// Test image attachment with invalid file
	os.Args = append(os.Args, "--image", "invalidfile")

	// Nothing to check - exercising code
	main()

	// Test invalid input and output
	os.Args = baseArgs
	os.Args = append(os.Args, "--html", "1", "--monospace", "1")

	// Nothing to check - exercising code
	main()

	// Test no server
	apiServer.Close()
	main()

	os.Args = savedArgs
}

func serverValidateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token, _ := r.Form["token"]
	if len(token) > 0 && token[0] == "fail" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"token":"invalid","errors":["application token is invalid"],"status":0,"request":"e8488e4e-dbe1-4795-a253-3ef644aa14a6"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":1,"group":0,"devices":["pixel2xl"],"licenses":["Android"],"request":"%s"}`, id)
}

func TestPushoverValidateCLI(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(serverValidateHandler))
	defer apiServer.Close()

	// Test valid input and output
	savedArgs := os.Args
	os.Args = []string{
		"pushover",
		"validate",
		"--pushoverurl", apiServer.URL,
		"--token", "token",
		"--user", "user",
	}

	// Nothing to check - exercising code
	main()

	// Nothing to check - exercising code
	os.Args[5] = "fail"
	main()

	// Test no server
	apiServer.Close()
	main()

	os.Args = savedArgs
}
