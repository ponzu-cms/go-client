package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ToValues converts a Content type to a *Values to use with a Ponzu Go client.
func ToValues(p interface{}) (*Values, error) {
	// encode p to JSON
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	// decode json to Values
	var kv map[string]interface{}
	err = json.Unmarshal(j, &kv)
	if err != nil {
		return nil, err
	}

	vals := NewValues()
	for k, v := range kv {
		switch v.(type) {
		case []interface{}:
			vv := v.([]interface{})
			for i := range vv {
				vals.Add(k, fmt.Sprintf("%v", vv[i]))
			}
		default:
			vals.Add(k, fmt.Sprintf("%v", v))
		}
	}

	return vals, nil
}

// Target represents required criteria to lookup single content items from the
// Ponzu Content API.
type Target struct {
	Type string
	ID   int
}

// ParseReferenceURI is a helper method which accepts a reference path / URI from
// a parent Content type, and retrns a Target containing a content item's Type
// and ID.
func ParseReferenceURI(uri string) (Target, error) {
	return parseReferenceURI(uri)
}

func parseReferenceURI(uri string) (Target, error) {
	const prefix = "/api/content?"
	if !strings.HasPrefix(uri, prefix) {
		return Target{}, fmt.Errorf("improperly formatted reference URI: %s", uri)
	}

	uri = strings.TrimPrefix(uri, prefix)

	q, err := url.ParseQuery(uri)
	if err != nil {
		return Target{}, fmt.Errorf("failed to parse reference URI: %s, %v", prefix+uri, err)
	}

	if q.Get("type") == "" {
		return Target{}, fmt.Errorf("reference URI missing 'type' value: %s", prefix+uri)
	}

	if q.Get("id") == "" {
		return Target{}, fmt.Errorf("reference URI missing 'id' value: %s", prefix+uri)
	}

	// convert query id string to int
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		return Target{}, err
	}

	return Target{Type: q.Get("type"), ID: id}, nil
}
