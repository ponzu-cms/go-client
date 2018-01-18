package client

import (
	"fmt"
	"net/url"
	"strings"
)

// Values is a modified, Ponzu-specific version of the Go standard library's
// url.Values, which implements most all of the same behavior. The exceptions are
// that its `Add(k, v string)` method converts keys into the expected format for
// Ponzu data containing slice type fields, and the `Get(k string)` method returns
// an `interface{}` which will either assert to a `string` or `[]string`.
type Values struct {
	values   url.Values
	keyIndex map[string]int
}

// Add updates the Values by including a properly formatted form value to data Values.
func (v *Values) Add(key, value string) {
	if v.keyIndex[key] == 0 {
		v.values.Set(key, value)
		v.keyIndex[key]++
		return
	}

	if v.keyIndex[key] == 1 {
		val := v.values.Get(key)
		v.values.Del(key)
		k := key + ".0"
		v.values.Add(k, val)
	}

	keyIdx := fmt.Sprintf("%s.%d", key, v.keyIndex[key])
	v.keyIndex[key]++
	v.values.Set(keyIdx, value)
}

// NewValues creates and returns an empty set of Ponzu values.
func NewValues() *Values {
	return &Values{
		values:   make(url.Values),
		keyIndex: make(map[string]int),
	}
}

// Del deletes a key and its value(s) from the data set.
func (v *Values) Del(key string) {
	if v.keyIndex[key] != 0 {
		// delete all key.0, key.1, etc
		n := v.keyIndex[key]
		for i := 0; i < n; i++ {
			v.values.Del(fmt.Sprintf("%s.%d", key, i))
			v.keyIndex[key]--
		}

		v.keyIndex[key] = 0
		return
	}

	v.values.Del(key)
	v.keyIndex[key]--
}

// Encode prepares the data set into a URL query encoded string.
func (v *Values) Encode() string { return v.values.Encode() }

// Get returns an `interface{}` value for the key provided, which will assert to
// either a `string` or `[]string`.
func (v *Values) Get(key string) interface{} {
	if strings.Contains(key, ".") {
		return v.values.Get(key)
	}

	if v.keyIndex[key] == 0 {
		return ""
	}

	if v.keyIndex[key] == 1 {
		return v.values.Get(key)
	}

	var results []string
	for i := 0; i < v.keyIndex[key]; i++ {
		keyIdx := fmt.Sprintf("%s.%d", key, i)
		results = append(results, v.values.Get(keyIdx))
	}

	return results
}

// Set sets a value for a key provided. If Set/Add has already been called, this
// will override all values at the key.
func (v *Values) Set(key, value string) {
	if v.keyIndex[key] == 0 {
		v.values.Set(key, value)
		v.keyIndex[key]++
		return
	}

	v.Del(key)
	v.keyIndex[key] = 0
	v.Set(key, value)
}
