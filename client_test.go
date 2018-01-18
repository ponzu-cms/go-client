package client

import (
	"net/http"
	"testing"
)

func TestMergeHeader(t *testing.T) {
	cases := map[string]string{
		"Content-Type": "text/plain",
		"X-Test-1":     "Value 1",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	// add some default to the original request header
	req.Header.Set("Content-Type", "text/plain")

	header := http.Header{}
	header.Set("X-Test-1", "Value 1")

	req = mergeHeader(req, header)

	for test, exp := range cases {
		res := req.Header.Get(test)
		if res != exp {
			t.Errorf("expected %s, got %s", exp, res)
		}
	}

	// Test multipart requests

	req, err = multipartForm("http://localhost:8080", NewValues(), nil)
	if err != nil {
		t.Errorf("Error creating multipart request: %v", err)
	}

	cases = map[string]string{
		"Content-Type": req.Header.Get("Content-Type"), // will have header from multipartForm
		"X-Test-1":     "Value 1",
	}

	req = mergeHeader(req, header)

	for test, exp := range cases {
		res := req.Header.Get(test)
		if res != exp {
			t.Errorf("expected %s, got %s", exp, res)
		}
	}
}
