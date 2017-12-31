package client_test

import (
	"fmt"
	"testing"

	client "github.com/ponzu-cms/go-client"
	"github.com/ponzu-cms/ponzu/system/item"
)

func TestToValues(t *testing.T) {
	type ContentExample struct {
		item.Item

		Name string   `json:"name"`
		ID   int      `json:"id"`
		Tags []string `json:"tags"`
	}

	ex := &ContentExample{
		Name: "Test case name",
		ID:   1,
		Tags: []string{"first", "second", "third"},
	}

	data, err := client.ToValues(ex)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if data.Get("name") != ex.Name {
		t.Fail()
	}

	if data.Get("id") != fmt.Sprintf("%d", ex.ID) {
		t.Fail()
	}

	if data.Get("tags.0") != "first" {
		t.Fail()
	}
	if data.Get("tags.1") != "second" {
		t.Fail()
	}
	if data.Get("tags.2") != "third" {
		t.Fail()
	}
}

func TestParseReferenceURI(t *testing.T) {
	cases := map[string]client.Target{
		"/api/content?type=Test&id=1": client.Target{Type: "Test", ID: 1},
	}

	for in, expected := range cases {
		got, err := client.ParseReferenceURI(in)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		if got.ID != expected.ID {
			fmt.Printf("expected: %v got: %v\n", expected.ID, got.ID)
			t.Fail()
		}
	}
}

func TestParseReferenceURIErrors(t *testing.T) {
	cases := map[string]string{
		"/api/content":                  "improperly formatted reference URI: /api/content",
		"/api/content?type=Test&noID=1": "reference URI missing 'id' value: /api/content?type=Test&noID=1",
		"/api/content?noType=Test&id=1": "reference URI missing 'type' value: /api/content?noType=Test&id=1",
	}

	for in, expected := range cases {
		_, err := client.ParseReferenceURI(in)
		if err.Error() != expected {
			fmt.Printf("got: %v, expected: %s\n", err, expected)
			t.Fail()
		}
	}
}
