package client

import (
	"fmt"
	"net/url"
)

// Search makes a GET request to return the Ponzu Search API response for the
// enpoint: `/api/search?type=<Type>&q=<Query String>`
func (c *Client) Search(contentType string, q string, opts *QueryOptions) (*APIResponse, error) {
	opts = setDefaultOpts(opts)

	endpoint := fmt.Sprintf(
		"%s/api/search?type=%s&q=%s&count=%d&offest=%d",
		c.Conf.Host, contentType, url.QueryEscape(q), opts.Count, opts.Offset,
	)

	if c.CacheEnabled() {
		ok, resp := c.Cache.Check(endpoint)
		if ok {
			return resp, nil
		}
	}

	w, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, err
}
