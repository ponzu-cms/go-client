package client

import (
	"fmt"
	"io/ioutil"
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

	req, err = multipartForm("http://localhost:8080", nil, nil)
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

func TestContentBySlug(t *testing.T) {
	// just show the response
	slug := "item-id-50407fd4-8d59-4448-a032-992812664ea6"
	resp, err := http.Get("http://localhost:8080/api/content?slug=" + slug)
	if err != nil {
		t.Errorf("failed to make request: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
	}
	resp.Body.Close()

	fmt.Println(string(b))
}
