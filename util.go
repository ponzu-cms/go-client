package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ToValues converts a Content type to a url.Values to use with a Ponzu Go client
func ToValues(p interface{}) (url.Values, error) {
	// encode p to JSON
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	// decode json to url.Values
	var kv map[string]interface{}
	err = json.Unmarshal(j, &kv)
	if err != nil {
		return nil, err
	}

	vals := make(url.Values)
	for k, v := range kv {
		switch v.(type) {
		case []string:
			s := v.([]string)
			for i := range s {
				if i == 0 {
					vals.Set(k, s[i])
				}

				vals.Add(k, s[i])
			}
		default:
			vals.Set(k, fmt.Sprintf("%v", v))
		}
	}

	return vals, nil
}
