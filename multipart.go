package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func multipartForm(endpoint string, params *Values, fileParams []string) (*http.Request, error) {
	if params == nil {
		return nil, errors.New("form data params must not be nil")
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for name := range params.values {
		if keyIsFile(name, fileParams) {
			if len(params.values[name]) > 1 {
				// iterate through file paths (has multiple, like FileRepeater),
				// make name => name.0, name.1, ... name.N
				files := params.values[name]
				for i := range files {
					fieldName := fmt.Sprintf("%s.%d", name, i)
					err := addFileToWriter(fieldName, files[i], writer)
					if err != nil {
						return nil, err
					}
				}
			} else {
				err := addFileToWriter(name, params.values.Get(name), writer)
				if err != nil {
					return nil, err
				}
			}
		} else {
			err := writer.WriteField(name, params.values.Get(name))
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
