package client

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Creates a new multipart/form-data http request with form data
func multipartForm(endpoint string, params url.Values) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for name := range params {
		err := writer.WriteField(name, params.Get(name))
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
