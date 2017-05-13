package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Creates a new multipart/form-data http request with form data
func multipartForm(endpoint string, params url.Values, fileParams []string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for name := range params {
		if keyIsFile(name, fileParams) {
			if len(params[name]) > 1 {
				// iterate through file paths (has multiple, like FileRepeater),
				// make name => name.0, name.1, ... name.N
				files := params[name]
				for i := range files {
					fieldName := fmt.Sprintf("%s.%d", name, i)
					err := addFileToWriter(fieldName, files[i], writer)
					if err != nil {
						return nil, err
					}
				}
			} else {
				err := addFileToWriter(name, params.Get(name), writer)
				if err != nil {
					return nil, err
				}
			}
		} else {
			err := writer.WriteField(name, params.Get(name))
			if err != nil {
				return nil, err
			}
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

func keyIsFile(key string, fileKeys []string) bool {
	for i := range fileKeys {
		if key == fileKeys[i] {
			return true
		}
	}

	return false
}

func addFileToWriter(fieldName, filePath string, w *multipart.Writer) error {
	paths := strings.Split(filePath, string(filepath.Separator))
	filename := paths[len(paths)-1]
	part, err := w.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(part, src)
	if err != nil {
		return err
	}
	return nil
}
