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

		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	ex := &ContentExample{
		Name: "Test",
		ID:   1,
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
}
